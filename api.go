package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/jmcvetta/randutil"
)

var adminHashes = [][]byte{
	[]byte{
		0xf7, 0xd9, 0xae, 0x94, 0xe9, 0x22, 0xa7, 0x35,
		0x19, 0x63, 0xb9, 0xea, 0x45, 0xd7, 0xed, 0x9f,
		0xbc, 0x3a, 0x8, 0x4c, 0x1a, 0x5e, 0xc5, 0x91,
		0xa6, 0x20, 0x54, 0x2a, 0xc4, 0x81, 0xdb, 0x28,
	},
}

func api(r chi.Router) {
	r.Route("/polls/{pollID}", func(r chi.Router) {
		r.Use(identify)
		r.Use(getPoll)
		r.Get("/movies", getMovies)
		r.Post("/vote", registerVote)
		r.Post("/suggest", suggestMovies)
	})
	r.With(checkAdmin).Post("/polls", setPoll)
}

func checkAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adminCookie, _ := r.Cookie("admin")
		if adminCookie == nil || adminCookie.Value == "" {
			userError(w, "nope")
			return
		}

		digest := sha256.Sum256([]byte(adminCookie.Value))

		var found bool
		for _, valid := range adminHashes {
			if bytes.Equal(valid, digest[:]) {
				found = true
				break
			}
		}

		if !found {
			userError(w, "nope")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)

		var ip IP
		err := db.First(&ip, "ip = ?", remoteIP).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			serverError(w, "db error", err)
			return
		}

		if ip.Access == AccessDeny {
			userError(w, "gtfo")
			return
		}

		var user User

		uidCookie, _ := r.Cookie("id")
		if uidCookie != nil && uidCookie.Value != "" {
			userID := uidCookie.Value

			err = db.First(&user, "id = ?", userID).Error
			if err != nil && !gorm.IsRecordNotFoundError(err) {
				serverError(w, "db error", err)
				return
			}

			if user.ID != "" {
				ctx := context.WithValue(r.Context(), "user", &user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		fpCookie, _ := r.Cookie("fp")
		if fpCookie == nil || len(fpCookie.Value) != 32 {
			userError(w, "invalid request")
			return
		}
		fingerprint := fpCookie.Value

		if ip.Access == AccessAllowFingerprinted {
			err = db.First(&user, "ip = ? AND fingerprint = ?", remoteIP, fpCookie.Value).Error
		} else {
			err = db.First(&user, "ip = ?", remoteIP).Error
		}
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			serverError(w, "db error", err)
			return
		}

		if user.ID != "" {
			http.SetCookie(w, &http.Cookie{
				Name:  "id",
				Value: user.ID,
			})

			ctx := context.WithValue(r.Context(), "user", &user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userID, err := randutil.AlphaString(32)
		if err != nil {
			serverError(w, "randutil", err)
			return
		}

		user = User{
			ID:          userID,
			IP:          remoteIP,
			Fingerprint: fingerprint,
		}

		err = db.Create(&user).Error
		if err != nil {
			serverError(w, "db error", err)
			return
		}

		ctx := context.WithValue(r.Context(), "user", &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPoll(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pollID := chi.URLParam(r, "pollID")

		var poll Poll
		var err error
		if pollID == "latest" {
			err = db.Order("created_at desc").First(&poll).Error
		} else {
			err = db.First(&poll, "id = ?", pollID).Error
		}
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				userError(w, "invalid poll id")
			} else {
				serverError(w, "db error", err)
			}

			return
		}

		ctx := context.WithValue(r.Context(), "poll", &poll)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type errorResponse struct {
	Error string `json:"error"`
}

func userError(w http.ResponseWriter, err string) {
	writeJSON(w, http.StatusBadRequest, errorResponse{Error: err})
}

func serverError(w http.ResponseWriter, context string, err error) {
	log.Printf("%s: %v", context, err)

	writeJSON(w, http.StatusInternalServerError, errorResponse{
		Error: "something bad happened",
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func decodeJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
