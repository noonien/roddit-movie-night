package main // import "github.com/noonien/roddit-movie-night"

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/packr"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func main() {
	initDB()
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.Recoverer)

	frontBox := packr.NewBox("frontend/dist")
	r.Mount("/", http.FileServer(frontBox))
	r.Route("api", api)

	http.ListenAndServe(":"+port, r)
}

func initDB() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL env missing")
	}

	var err error
	db, err = gorm.Open(dbURL)
	if err != nil {
		log.Fatalf("couldn't connect to db: %v", err)
	}

	db.AutoMigrate(&IP{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Poll{})
	db.AutoMigrate(&Movie{})
	db.AutoMigrate(&Vote{})
}
