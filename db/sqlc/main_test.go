package db

import (
	"database/sql"
	"github.com/jhlee1/simplebank/db/util"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig(".") // "." because the config file is in the current folder
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
