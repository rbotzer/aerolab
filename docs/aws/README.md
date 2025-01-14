# How to setup and use AWS with Aerolab

## Prerequisites

### Create Credentials file

There are two ways to create a credentials file:

#### Using aws-cli

1. Download and install [aws-cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
2. Run the command `aws configure`

#### Manually

Create `~/.aws` directory, inside the directory create 2 files:

`~/.aws/credentials`

```
[default]
aws_access_key_id = KEYID
aws_secret_access_key = SECRETKEY
```

`~/.aws/config`

```
[default]
region = DEFAULT_REGION_TO_USE
```

### Configure AWS account

#### Security Groups

AeroLab clusters require a security group. The following rules will allow for full connectivity:
1. Create a security group (sg-xxxx) with the rule to allow all outbound (default) and inbound: port 22(tcp) from ANY IP
2. Edit the security group (sg-xxxx) and add a rule to allow all inbound on all ports coming from itself (source: sg-xxxx)

If planning to deploy AMS or other clients:
1. Create a security group (sg-yyyy) with a rule to allow all outbound (default) and allow inbound from ANY IP to the following TCP ports: 22, 3000, 8888, 8080
2. Edit the security group (sg-yyyy) adding 2 rules:
   a) allow all ports from self source (sg-yyyy)
   b) allow all ports from server source (sg-xxxx)
3. Edit the security group (sg-xxxx), and add a rule allowing all ports from the client source security group (sg-yyyy)

Use (sg-xxxx) for clusters and (sg-yyyy) for client machines.

#### Subnets

If creating a new subnet and/or VPC, configure the VPC and Subnet such that:
* the instances will have automatically assigned public DNS
* the instances will have automatically assigned public IP addresses

### Configure the backend in aerolab

The most basic configuration is: `aerolab config backend -t aws`

To specify a custom location where SSH keys will be stored and override the default aws region config, extra parameters may be supplied:

```
aerolab config backend -t aws -p /path/where/keys/will/be/stored -r AWS_REGION
```

## Deploy cluster in aws

Extra parameters are required when working with the `aws` backend as opposed to the `docker` backend.

Executing `aerolab cluster create help` once the backend has been selected will display the relevant options.

### Example:

```bash
./aerolab cluster create -n testcluster -c 3 -m mesh -I t3a.medium -E 20 -S sg-03430d698bffb44a3 -U subnet-06cc8a834647c4cc3
```

## Destroy cluster
```bash
./aerolab cluster-destroy -f -n testcluster
```

## Other commands

All commands are supported on both `aws` and `docker` backends and should behave exactly the same.

## Working with multiple regions

Working with multiple regions can be achieved by switching the backend, as so:

```
aerolab config backend -t aws -r eu-west-1
...commands...
aerolab config backend -t aws -r us-east-1
...commands...
```

Alternatively, if using multiple regions on many occassions, multiple configuration files may be utilised:

```
# create a config called us.conf
AEROLAB_CONFIG_FILE=us.conf
aerolab config backend -t aws -r us-east-1

# create a config called eu.conf
AEROLAB_CONFIG_FILE=eu.conf
aerolab config backend -t aws -r eu-west-1

# since eu is the exported region variable, default commands execute against it
aerolab cluster create
aerolab attach shell -- asadm -e info

# execute an ad-hoc command on another region
AEROLAB_CONFIG_FILE=us.conf aerolab cluster create

# keep running in eu region
aerolab cluster destroy
```

## Note on shared AWS accounts and KeyPairs

AeroLab aims to create and destroy SSH key pairs as needed. having said that, if a particular cluster is created by user X, user Y can only access the cluster if user X shares their key pair for that cluster.

By default keys are stored in `${HOME}/aerolab-keys`
