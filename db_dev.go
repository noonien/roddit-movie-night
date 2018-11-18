// this file is only included in non-heroku builds
// +build !heroku

package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func initDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("couldn't connect to db: %v", err)
	}

	return db
}
