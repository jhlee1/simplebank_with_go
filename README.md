# This is an example source code from the [go backend master course](https://www.youtube.com/watch?v=prh0hTyI1sU)


## Things to install besides the go packages in the source code
- [migrate](https://github.com/golang-migrate/migrate)
- [sqlc](https://sqlc.dev/)
- [Docker](https://www.docker.com/products/docker-desktop/)
- postgreSQL docker image

## Setup PostgreSQL

```bash
  docker pull postgres:15.2-alpine
  docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecretpassword -d postgres
  docker exec -it postgres15 createdb --username=root --owner=root simple_bank
```

## create a migration version
```bash
$ migrate create -ext sql -dir db/migrations -seq initial_schema
```

## Build docker image
```bash
$ docker build -t simplebank:latest .
```

## Run built image with IP address to connect DB from the server
1. Get the IP address from
```shell
$ docker inspect postgres15
```
2. Set the environment variable for the DB
```shell
$ docker run --name simplebank -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:mysecretpassword@172.17.0.3:5432/simple_bank?sslmode=disable" simplebank:latest
```

## Set the nextwork using bridge
The IP address is changed everytime. You cannot connect the DB using IP address every time
1. Get the network info
```shell
$ docker network ls
```
```text
NETWORK ID     NAME      DRIVER    SCOPE
58ae87dbc62d   bridge    bridge    local
a7fcbd91b5d2   host      host      local
05fb840db4c1   none      null      local
```

```shell
$ docker inspect bridge
```
```json
{
  "Containers": {
    "638c00a5e2817c24e269e07861c19718d85bd2ab1631166f446a8efd5906f7f6": {
      "Name": "simplebank",
      "EndpointID": "021dd91a960b57f58ff1b430dfa1832ac1341b263f41e86333cc64ed700a2342",
      "MacAddress": "02:42:ac:11:00:02",
      "IPv4Address": "172.17.0.2/16",
      "IPv6Address": ""
    },
    "cac873975734704ab29e08712d56b769f0bb7b6d551170964383f3a5871c19b8": {
      "Name": "postgres15",
      "EndpointID": "81e59dbb75948f4787967f59f49f991eb721d048782698e04f2bfa90f5d88d29",
      "MacAddress": "02:42:ac:11:00:03",
      "IPv4Address": "172.17.0.3/16",
      "IPv6Address": ""
    }
  }
}
```
Usually docker containers can connect to each other using the name if they are in the same network. However, default bridge network does not support that

2. Create a new network
```shell
$ docker network create bank-network
$ docker network connect bank-network postgres15
```
Check if it is connected using 
```shell
$ docker network inspect bank-network
$ docker container inspect postgres15
```

```json
{
  "Containers": {
    "cac873975734704ab29e08712d56b769f0bb7b6d551170964383f3a5871c19b8": {
        "Name": "postgres15",
        "EndpointID": "935cb1dea1b8d1c9458b9d906b13586b9ed9fedbfdd753bce1da75bd3250141c",
        "MacAddress": "02:42:ac:12:00:02",
        "IPv4Address": "172.18.0.2/16",
        "IPv6Address": ""
    }
  }
}
```

## Create a container that connects to the created network
```shell
$ docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:mysecretpassword@postgres15:5432/simple_bank?sslmode=disable" simplebank:latest
```

## How to setup github actions to build and push docker image to aws ECR
I have followed [this guide](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-amazon-web-services)
Most tutorials and resources use `user` and `user group`. However, the recent [aws-credential](https://github.com/aws-actions/configure-aws-credentials) has changed the way of authentication. It requires to use OIDC (OpenID Connect).

### Create a new Identity provider
1. Go to IAM > Identity providers > Add provider
2. Select OpenID Connect
3. Set the Provider URL to `https://token.actions.githubusercontent.com`
4. Set the Audience to `sts.amazonaws.com`

### Assign a role to the Identity provider
1. Go to IAM > Roles > Create role
2. Select Web identity
3. Select the Identity provider created above
4. Set the Audience to `sts.amazonaws.com`
5. Add the permission `AmazonEC2ContainerRegistryFullAccess`
6. Set the role name to `anything you want`

### Add ARN to the deployment yaml file

1. Copy the ARN from the role created above
2. Add the ARN to aws-credential step in the yaml file
3. Add permissions to the job

This is a part of the deployment yaml file

```yaml
name: Deploy to production

on:
  push:
    branches: [ "main" ]

jobs:
  build:
    name: Build image
    permissions:
      id-token: write
      contents: read
    runs-on: ubuntu-latest

    steps:
      - name: Checkout repo
        uses: actions/checkout@v3
        
      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2 # More information on this action can be found below in the 'AWS Credentials' section
        with:
          role-to-assume: arn:aws:iam::847768247459:role/github-action
          aws-region: us-east-2
  
```

## Create the DB using RDS
## Set the env variables using Secrets Manager from aws

Create a random symmetric key using
```shell
openssl rand -hex 64 | head -c 32
```

Get the secret value using awscli
```shell
aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text
```
Install jq to parse the json
```shell
brew install jq
aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > app.env
```





