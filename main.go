package main

import (
	"database/sql"
	"github.com/jhlee1/simplebank/api"
	db "github.com/jhlee1/simplebank/db/sqlc"
	"github.com/jhlee1/simplebank/db/util"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	config, err := util.LoadConfig(".") // "." because the config file is in the current folder
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
