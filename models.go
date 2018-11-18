package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type AccessMode string

const (
	AccessAllow = "allow"
	AccessDeny  = "deny"
)

type IP struct {
	gorm.Model

	IP     string `gorm:"unique; not null"`
	Access AccessMode
	Note   string
}

type User struct {
	gorm.Model

	Name        string `gorm:"unique; not null"`
	IP          string `gorm:"not null"`
	RandID      string `gorm:"not null"`
	Fingerprint string `gorm:"not null"`
}

type Poll struct {
	gorm.Model

	Name      string `gorm:"not null"`
	Info      string
	ExpiresAt time.Time
}

type Movie struct {
	gorm.Model

	Name       string `json:"name"`
	IMDBId     string `json:"imdb_id"`
	TrailerURL string `json:"trailer_url"`

	SuggestedBy int `gorm:"foreignkey:User"`
}

type Vote struct {
	gorm.Model

	UserID  int
	MovieID int
}
