package models

import (
	"net/http"
	"time"
)

type Comment struct{
	ID int `json:"id"`
	Comment string `json:"comment"`
	MovieID int `json:"movie_id"`
	IPAddress string `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentRepository interface{
	ListComments(EpisodeID) ([]Comment, error)
	HandleAddComment(http.Response, http.Request)
	HandleListComments(http.Response, http.Request)

}

type CommentService interface{

}