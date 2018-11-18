package main

import (
	"net/http"
	"time"

	"github.com/rs/xid"
)

type VotedMovie struct {
	Movie
	Votes int `json:"numVotes"`
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	poll, _ := r.Context().Value("poll").(*Poll)

	var res []VotedMovie
	err := db.Table("movies").
		Select("movies.*, count(votes.user_id) as votes").
		Joins("LEFT JOIN votes ON votes.movie_id = movies.id").
		Where("movies.poll_id = ?", poll.ID).
		Group("movies.id").
		Scan(&res).Error
	if err != nil {
		serverError(w, "db error", err)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func registerVote(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user").(*User)
	poll, _ := r.Context().Value("poll").(*Poll)

	if poll.ClosesAt.Before(time.Now()) {
		userError(w, "poll is closed")
		return
	}

	var movies struct{ Votes []string }
	err := decodeJSON(r, &movies)
	if err != nil {
		userError(w, "invalid json: "+err.Error())
		return
	}

	// TODO: validate movies

	err = db.Where("user_id = ? AND poll_id = ?", user.ID, poll.ID).
		Delete(Vote{}).Error
	if err != nil {
		serverError(w, "db error", err)
		return
	}

	for _, mID := range movies.Votes {
		vote := &Vote{
			UserID:  user.ID,
			PollID:  poll.ID,
			MovieID: mID,
		}

		err = db.Create(&vote).Error
		if err != nil {
			serverError(w, "db error", err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func suggestMovies(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user").(*User)
	poll, _ := r.Context().Value("poll").(*Poll)

	if poll.ClosesAt.Before(time.Now()) {
		userError(w, "poll is closed")
		return
	}

	var req struct {
		IMDBURL string `json:"imdbURL"`
	}
	err := decodeJSON(r, &req)
	if err != nil {
		userError(w, "invalid json: "+err.Error())
		return
	}

	movie := &Movie{
		ID:          xid.New().String(),
		PollID:      poll.ID,
		Name:        "name",
		Image:       "image",
		IMDBURL:     req.IMDBURL,
		SuggestedBy: user.ID,
	}

	err = db.Create(&movie).Error
	if err != nil {
		serverError(w, "db error", err)
		return
	}

	writeJSON(w, http.StatusOK, movie)
}

func setPoll(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var poll Poll
	err := decodeJSON(r, &poll)
	if err != nil {
		userError(w, "invalid json: "+err.Error())
		return
	}

	if poll.ID == "" {
		poll.ID = xid.New().String()
		err = db.Create(&poll).Error
	} else {
		err = db.Save(&poll).Error
	}
	if err != nil {
		serverError(w, "db error", err)
		return
	}

	writeJSON(w, http.StatusOK, &poll)
}
