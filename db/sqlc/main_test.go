package db

import (
	"database/sql"
	"fmt"
	"log"
	"nebil/golang/utils"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://root:mysecretpassword@localhost:5432/simplebank?sslmode=disable"
// )

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../../")
	if err != nil {
		log.Fatal("can not load config file")
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		fmt.Println(err)
		return
	}

	testQueries = New(testDB)
	fmt.Println(testQueries)

	os.Exit(m.Run())
}
