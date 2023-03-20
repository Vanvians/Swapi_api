package models

import "time"

type Movie struct{
	ID string `json:"id"`
	EpisodeID int `json:"episode_id"`
	Title string `json:"title"`
	Director string `json:"director"`
	Producer string `json:"producer"`
	ReleaseDate time.Time `json:"release_date"`
	OpeningCrawl string `json:"opening_crawl"`
	Characters []string `json:"characters"`
	Comments []Comment `json:"comments"`
}