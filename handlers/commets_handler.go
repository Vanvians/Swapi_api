package handlers

import (
	"net/http"
	"strconv"

	"github.com/jcezetah/swapi_api/models"
	"github.com/swapi_api/services"
	"github.com/swapi_api/utils"

	"github.com/gorilla/mux"
)

type CommentsHandler struct{
	service services.CommentService
}

func NewCommentsHandler(service services.CommentService) *MovieHandler{
	return &MovieHandler{
		service:service,
	}
}

//returns a handler function that lists comments for a movie
func (mh *MovieHandler) GetCommentsHandler() http.HandlerFunc{
	return func(w http.response, r *http.Request) {
		movieID := mux.Vars(r)["id"]

		comments, err := mh.service.GetComments(movieID)
		if err != nil{
			utils.WriteError(w,err)
			return
		}
		utils.WriteJson(w, http.StatusOK, comments)
	}
}

//returns a handler function that adds a new comment for a movie
func (mh *MovieHandler) AddCommentHandler() http.HandlerFunc{
	return func (w http.ResponeWriter, r *http.Request) {
		movieID := mux.Vars(r)["id"]
		comment := &models.Comment{}
		if err != nil{
			utils.WriteError(w, utils.NewApiError(http.StatusBadRequest, "Invalid request payload"))
			return
		}
		
		err = mh.service.AddComment(movieID,comment)
		if err != nil{
			utils.WriteError(w,err)
			return
		}

		utils.WriteJson(w, http.statusCreated, comment)
	}
}