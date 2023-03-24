package models

import (
	"time"
)

type Comment struct{
	ID int `json:"id"`
	Comment string `json:"comment"`
	MovieID int `json:"movie_id"`
	IPAddress string `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

