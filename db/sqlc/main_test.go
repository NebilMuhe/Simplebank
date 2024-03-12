package db

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:mysecretpassword@localhost:5432/simplebank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		fmt.Println(err)
		return
	}

	testQueries = New(testDB)
	fmt.Println(testQueries)

	os.Exit(m.Run())
}
