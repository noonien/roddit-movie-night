package main // import "github.com/noonien/roddit-movie-night"

import (
	"net/http"
	"os"
	"time"

	"github.com/noonien/roddit-movie-night/frontend"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	r.Mount("/", http.FileServer(frontend.Box))
	r.Route("/api", api)

	http.ListenAndServe(":"+port, r)
}
