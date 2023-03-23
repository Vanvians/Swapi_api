package models


type Movie struct {
	MovieId   int    `json:"episode_id"`
	Title       string `json:"title"`
	OpeningCrawl string `json:"opening_crawl"`
	CommentCount int    `json:"comment_count"`
	ReleaseDate string `json:"release_date"`
}

type Character struct {
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	HeightCm  int    `json:"height_cm"`
	HeightFt  string `json:"height_ft"`
}