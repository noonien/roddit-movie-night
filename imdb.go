package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

var imdbURLRe = regexp.MustCompile(`imdb.com/title/(tt[0-9]{7})`)
var imdbBaseURL = "https://www.imdb.com/title/%s/"
var omdbBaseURL = "https://www.omdbapi.com/?apikey=%s&i=%s"

var omdbAPIKey string

func init() {
	omdbAPIKey = os.Getenv("OMDB_KEY")
	if omdbAPIKey == "" {
		omdbAPIKey = "BanMePlz"
	}
}

type OMDBResponse struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Runtime  string `json:"Runtime"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Language string `json:"Language"`
	Country  string `json:"Country"`
	Awards   string `json:"Awards"`
	Poster   string `json:"Poster"`
	Ratings  []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	ImdbID     string `json:"imdbID"`
	Type       string `json:"Type"`
	DVD        string `json:"DVD"`
	BoxOffice  string `json:"BoxOffice"`
	Production string `json:"Production"`
	Website    string `json:"Website"`
	Response   string `json:"Response"`
}

func getMovieInfo(imdbURL string) (*Movie, error) {
	match := imdbURLRe.FindStringSubmatch(imdbURL)
	if len(match) == 0 {
		return nil, errors.New("invalid imdb url")
	}

	imdbID := match[1]
	imdbURL = fmt.Sprintf(imdbBaseURL, imdbID)
	omdbURL := fmt.Sprintf(omdbBaseURL, omdbAPIKey, imdbID)

	resp, err := http.Get(omdbURL)
	if err != nil {
		return nil, errors.New("could not fetch movie info")
	}
	defer resp.Body.Close()

	var oresp OMDBResponse
	json.NewDecoder(resp.Body).Decode(&oresp)
	if err != nil && oresp.Response != "True" {
		return nil, errors.New("could not parse movie info")
	}

	if oresp.Type != "movie" {
		return nil, errors.New("not a movie")
	}

	var ratings []MovieRating
	for _, r := range oresp.Ratings {
		ratings = append(ratings, MovieRating{
			Source: r.Source,
			Rating: r.Value,
		})
	}

	return &Movie{
		Title:    oresp.Title,
		Year:     oresp.Year,
		Rated:    oresp.Rated,
		Runtime:  oresp.Runtime,
		Genre:    oresp.Genre,
		Director: oresp.Director,
		Writer:   oresp.Writer,
		Actors:   oresp.Actors,
		Plot:     oresp.Plot,
		Language: oresp.Language,
		Poster:   oresp.Poster,
		Ratings:  ratings,
		URL:      imdbURL,
	}, nil
}
