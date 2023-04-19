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
