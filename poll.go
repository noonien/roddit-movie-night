package main

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
)

type VotedMovie struct {
	Movie
	Votes int `json:"votes"`
}

func getPoll(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)
	poll := r.Context().Value("poll").(*Poll)

	var res struct {
		Poll   *Poll        `json:"poll"`
		Movies []VotedMovie `json:"movies"`
		Votes  []string     `json:"votes"`
	}
	res.Poll = poll

	err := db.Table("movies").
		Select("movies.*, count(votes.user_id) as votes").
		Joins("LEFT JOIN votes ON votes.movie_id = movies.id").
		Where("movies.poll_id = ?", poll.ID).
		Group("movies.id").
		Scan(&res.Movies).Error
	if err != nil {
		serverError(w, "db error", err)
		return
	}

	movies := res.Movies
	sort.Slice(movies, func(i, j int) bool {
		mi, mj := movies[i], movies[j]
		return mi.Votes > mj.Votes || strings.Compare(mi.Title, mj.Title) < 0
	})

	var votes []Vote
	err = db.Where("user_id = ? AND poll_id = ?", user.ID, poll.ID).
		Find(&votes).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		serverError(w, "db error", err)
		return
	}

	res.Votes = make([]string, 0, len(votes))
	for _, v := range votes {
		res.Votes = append(res.Votes, v.MovieID)
	}

	writeJSON(w, http.StatusOK, res)
}

func registerVote(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)
	poll := r.Context().Value("poll").(*Poll)

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

	vs := map[string]bool{}
	votes := make([]string, 0, len(movies.Votes))
	for _, mID := range movies.Votes {
		if _, exists := vs[mID]; exists {
			continue
		}

		votes = append(votes, mID)
	}

	var count int
	err = db.Model(Movie{}).Where("id IN (?)", votes).Count(&count).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		serverError(w, "db error", err)
		return
	}
	if count != len(votes) {
		userError(w, "invalid movie id")
		return
	}

	err = db.Where("user_id = ? AND poll_id = ?", user.ID, poll.ID).
		Delete(Vote{}).Error
	if err != nil {
		serverError(w, "db error", err)
		return
	}

	for _, mID := range votes {
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

	var count int
	err := db.Model(Movie{}).
		Where("poll_id = ? AND suggested_by = ?", poll.ID, user.ID).
		Count(&count).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		serverError(w, "db error", err)
		return
	}
	if count >= 2 {
		userError(w, "cannot suggest more movies")
		return
	}

	var req struct {
		URL string `json:"url"`
	}
	err = decodeJSON(r, &req)
	if err != nil {
		userError(w, "invalid json: "+err.Error())
		return
	}

	movie, err := getMovieInfo(req.URL)
	if err != nil {
		userError(w, err.Error())
		return
	}
	movie.ID = xid.New().String()
	movie.PollID = poll.ID
	movie.SuggestedBy = user.ID

	var emovie Movie
	err = db.First(&emovie, "poll_id = ? AND url = ?", poll.ID, movie.URL).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		serverError(w, "db error", err)
		return
	} else if err == nil {
		userError(w, "movie was already suggested")
		return
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
