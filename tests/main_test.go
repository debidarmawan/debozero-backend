package database_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	database "github.com/debidarmawan/debozero/database/sqlc"
	"github.com/debidarmawan/debozero/utils"
	_ "github.com/lib/pq"
)

var testQuery *database.Queries

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Could not load env config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Could not connect to database", err)
	}

	testQuery = database.New(conn)

	os.Exit(m.Run())
}
