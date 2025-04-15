package main

import "github.com/brunocordeiro180/go-rate-limiter/internal/database"

func main() {
	redisConn := database.NewRedisConnection("localhost:6379", "", 0)
	db := database.NewDatabaseConnector(redisConn)

	if err := db.Connect(); err != nil {
		panic(err)
	}
	defer db.Disconnect()
}
