package main

import (
	"hourBot/bot"
	"sync"
)

var lock = &sync.Mutex{}

func main() {
	connectionString := "host=localhost port=5432 user=user dbname=gorm password=pass sslmode=disable"
	bot.Run(connectionString)
}
