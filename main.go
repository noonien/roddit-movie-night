package main // import "github.com/noonien/roddit-movie-night"

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/packr"
)

func main() {
	db = initDB()
	migrate(db)
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
	r.Route("/api", api)

	http.ListenAndServe(":"+port, r)
}
