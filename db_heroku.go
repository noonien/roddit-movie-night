// this file is only included in heroku builds
// +build heroku

package main

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func initDB() *gorm.DB {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL env missing")
	}

	var err error
	db, err = gorm.Open("postgresql", dbURL)
	if err != nil {
		log.Fatalf("couldn't connect to db: %v", err)
	}

	return db
}
