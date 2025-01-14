# Jupyter client

Jupyter is a web-based interface allowing for interactive testing and development of code. It is intended for testing purposes only.

Supported jupyter kernels are: `go,python,dotnet,java,node`

## Help pages

```
aerolab client create jupyter help
```

## Create a jupyter client machine with all kernels

```
aerolab client create jupyter -n jupyter
```

## Create a jupyter client machine with go and python kernels only

```
aerolab client create jupyter -n jupyter -k go,python
```

## Add dotnet kernel to existing jupyter client

```
aerolab client create jupyter -n jupyter -k dotnet
```

## Connect

First get IP from:

```
aerolab client list
```

Then, in the browser, navigate to:

```
http://IP:8888
```

## Seed IPs

The `-s` switch can be used when creating a jupyter client. This will make the creation process auto-fill the IPs in the example code with the cluster IPs for seeding.

If you haven't used the `-s` switch when creating a jupyter client, remember to adjust the seed IP address inside the jupyter GUI when editing your code, or you won't be able to connect.

## Connectivity notes on docker desktop

Please note that if using `docker desktop`, accessing the containers directly from your host machine is not allowed. In this case, when running the create commands above, add the `-e 8888:8888` switch to expose port 8888 to the host. Then access the jupyter server in your browser using `http://127.0.0.1:8888`

Example:

```
aerolab client create jupyter -n jupyter -e 8888:8888
```
