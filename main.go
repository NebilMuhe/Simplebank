package main

import (
	"database/sql"
	"log"
	"nebil/golang/api"
	db "nebil/golang/db/sqlc"
	"nebil/golang/utils"

	_ "github.com/lib/pq"
)

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://root:mysecretpassword@localhost:5432/simplebank?sslmode=disable"
// 	address  = ":8080"
// )

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load config file", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not connect to db", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("can not create server %v", err)
	}

	err = server.Start(config.Address)
	if err != nil {
		log.Fatal("Can not start the server", err)
	}
}

// func main() {
// 	slice := []int{1, 2, 3}

// 	updatedArray := Update(slice)

// 	fmt.Println("Updated Slice ", updatedArray)
// 	fmt.Println("Old Slice ", slice)
// 	fmt.Println("Cap Updated Slice ", cap(updatedArray))
// 	fmt.Println("Cap Slice  ", cap(slice))
// }

// func Update(nums []int) []int {
// 	nums[1] = 4
// 	nums = append(nums, 5)
// 	return nums
// }
