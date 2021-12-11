package main

import (
	"log"

	"github.com/nikochiko/chat-go/chat"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres dbname=chatgo port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	listener := chat.NewHTTPListener(db, log.Default())
	listener.AutoMigrate()

	listener.Listen()
}
