package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	flags "github.com/rglonek/jeddevdk-goflags"
)

// TODO: handle 'help' calls using custom help

var apiQueue = make(chan apiQueueItem, 1024)

type apiQueueItem struct {
	w      http.ResponseWriter
	r      *http.Request
	finish chan int
}

func (c *restCmd) handleApi(w http.ResponseWriter, r *http.Request) {
	f := make(chan int, 1)
	apiQueue <- apiQueueItem{
		w:      w,
		r:      r,
		finish: f,
	}
	<-f
}

func (c *restCmd) handleApiDo() {
	for {
		c.handleApiDoLoop()
	}
}

func (c *restCmd) handleApiDoLoop() {
	queueItem := <-apiQueue
	w := queueItem.w
	r := queueItem.r
	defer close(queueItem.finish)
	if strings.Trim(r.URL.Path, "/") == "quit" {
		w.Write([]byte("OK"))
		os.Exit(0)
	}
	// command = []string{"xdr","connect"}
	// handle and parse payload to command struct
	// execute command struct .Execute
	pr, pw, err := os.Pipe()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	os.Stdout = pw
	os.Stderr = pw
	log.SetOutput(pw)
	buf := new(bytes.Buffer)
	defer func() {
		os.Stdout = c.stdout
		os.Stderr = c.stderr
		log.SetOutput(c.logout)
	}()
	defer pr.Close()
	defer pw.Close()
	go c.copy(buf, pr)

	var na = &aerolab{
		opts: new(commands),
	}
	na.parser = flags.NewParser(na.opts, flags.HelpFlag|flags.PassDoubleDash)
	na.iniParser = flags.NewIniParser(na.parser)
	na.parseFile()
	na.parser.ParseArgs([]string{})

	command := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	keys := []string{}
	keyField := reflect.ValueOf(na.opts).Elem()
	v, err := c.findCommand(keyField, strings.Join(keys, "."), "", []string{}, command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(v.Addr().Interface())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if command[0] == "attach" {
		if len(c.getTail(v)) == 0 {
			http.Error(w, "Tail is not optional for attach commands via the rest api", http.StatusBadRequest)
			return
		}
	}

	c.apiRunCommand(v, w, buf)
}

func (c *restCmd) copy(buf *bytes.Buffer, pr *os.File) {
	for {
		_, err := io.Copy(buf, pr)
		if err != nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (c *restCmd) getTail(v *reflect.Value) []string {
	tail := v.FieldByName("Tail")
	if !tail.IsValid() {
		return []string{}
	}
	return tail.Interface().([]string)
}

func (c *restCmd) apiRunCommand(v *reflect.Value, w http.ResponseWriter, buf *bytes.Buffer) {
	// run the execute command and stream results back to response.data, if error from command, also set response.err
	tailv := c.getTail(v)
	tail := []reflect.Value{reflect.ValueOf(tailv)}
	outv := v.Addr().MethodByName("Execute").Call(tail)
	out := outv[0].Interface()
	switch out.(type) {
	case error:
		w.WriteHeader(http.StatusInternalServerError)
	}
	time.Sleep(20 * time.Millisecond)
	io.Copy(w, buf)
	switch err := out.(type) {
	case error:
		w.Write([]byte(err.Error()))
	}
}

func (c *restCmd) findCommand(keyField reflect.Value, start string, tags reflect.StructTag, tagStack []string, command []string) (v *reflect.Value, err error) {
	tagCommand := tags.Get("command")
	if tagCommand != "" {
		tagStack = append(tagStack, tagCommand)
	}
	for i, t := range tagStack {
		if (len(command) > i && command[i] != t) || len(command) <= i {
			return
		}
	}
	if testEq(command, tagStack) {
		return &keyField, nil
	}
	switch keyField.Type().Kind() {
	case reflect.Struct:
		for i := 0; i < keyField.NumField(); i++ {
			fieldName := keyField.Type().Field(i).Name
			fieldTag := keyField.Type().Field(i).Tag
			if len(fieldName) > 0 && fieldName[0] >= 97 && fieldName[0] <= 122 {
				if keyField.Field(i).Type().Kind() != reflect.Struct {
					continue
				}
				v, err = c.findCommand(keyField.Field(i), start, fieldTag, tagStack, command)
				if err != nil || v != nil {
					return
				}
			}
			if len(fieldName) == 0 || fieldName[0] < 65 || fieldName[0] > 90 {
				continue
			}
			if start != "" {
				fieldName = start + "." + fieldName
			}
			if strings.HasPrefix(fieldName, "Config.Defaults.") || fieldName == "DryRun" {
				continue
			}
			v, err = c.findCommand(keyField.Field(i), fieldName, fieldTag, tagStack, command)
			if err != nil || v != nil {
				return
			}
		}
	}
	return
}

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}