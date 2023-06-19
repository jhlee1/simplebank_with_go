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









