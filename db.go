package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type AccessMode string

const (
	AccessAllowFingerprinted = "allow_fingerprinted"
	AccessDeny               = "deny"
)

type IP struct {
	IP     string `gorm:"primarykey; not null"`
	Access AccessMode
	Note   string
}

type User struct {
	ID string `gorm:primarykey; not null"`

	IP          string `gorm:"not null"`
	Fingerprint string `gorm:"not null"`

	CreatedAt time.Time
}

type Poll struct {
	ID string `json:"id" gorm:"primarykey; not null"`

	Name string `json:"name" gorm:"not null"`
	Info string `json:"info"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ClosesAt  time.Time `json:"closes_at"`
}

type MovieRating struct {
	Source string `json:"source"`
	Rating string `json:"rating"`
}

type MovieRatings []MovieRating

func (mr MovieRatings) Value() (driver.Value, error) {
	j, err := json.Marshal(mr)
	return j, err
}

func (mr *MovieRatings) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		sourceStr, ok := src.(string)
		if !ok {
			return errors.New("could not type assert MovieRatings")
		}
		source = []byte(sourceStr)
	}

	err := json.Unmarshal(source, mr)
	if err != nil {
		return err
	}

	return nil
}

type Movie struct {
	ID     string `json:"id" gorm:"primarykey; not null"`
	PollID string `json:"-" gorm:"not null"`

	Title    string       `json:"title"`
	Year     string       `json:"year"`
	Rated    string       `json:"rated"`
	Runtime  string       `json:"runtime"`
	Genre    string       `json:"genre"`
	Director string       `json:"director"`
	Writer   string       `json:"writer"`
	Actors   string       `json:"actors"`
	Plot     string       `json:"plot"`
	Language string       `json:"language"`
	Poster   string       `json:"poster"`
	Ratings  MovieRatings `json:"ratings" gorm:"type:text"`
	URL      string       `json:"url"`

	SuggestedBy string    `json:"-" gorm:"foreignkey:User"`
	CreatedAt   time.Time `json:"created_at"`
}

type Vote struct {
	PollID  string
	UserID  string
	MovieID string

	CreatedAt time.Time
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&IP{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Poll{})
	db.AutoMigrate(&Movie{})
	db.AutoMigrate(&Vote{})
}
