package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func api(r chi.Router) {
	r.Use(identify)
	r.Use(getPoll)
	r.Get("/polls/{pollID}/movies", getMovies)
	r.Post("polls/{pollID}/vote", registerVote)
	r.Post("/polls/{pollID}/suggest", suggestMovies)
}

func identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//ctx := context.WithValue(r.Context(), "article", article)
		//next.ServeHTTP(w, r.WithContext(ctx))
		next.ServeHTTP(w, r)
	})
}

func getPoll(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pollID := chi.URLParam(r, "pollID")

		var poll Poll
		if pollID == "latest" {
			err := db.Order("created_at desc").First(&poll)
			if err != nil {
				apiError(w, "something bad happened")
			}
		} else {
			pID, err := strconv.Atoi(pollID)
			if err != nil {
				apiError(w, "invalid poll id")
				return
			}

			err = db.First(&poll, "id = ?", pID).Error
		}

		ctx := context.WithValue(r.Context(), "poll", poll)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getMovies(w http.ResponseWriter, r *http.Request) {

}

func registerVote(w http.ResponseWriter, r *http.Request) {
}

func suggestMovies(w http.ResponseWriter, r *http.Request) {
}

type errorResponse struct {
	Error string `json:"error"`
}

func apiError(w http.ResponseWriter, err string) {
	errResp := errorResponse{Error: err}
	json.NewEncoder(w).Encode(errResp)
}
