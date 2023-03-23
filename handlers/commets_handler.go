package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/jcezetah/Swapi_api/models"
	"github.com/jcezetah/Swapi_api/services"
	"github.com/jcezetah/Swapi_api/utils"
)

type CommentsHandler struct {
	commentService services.CommentServiceInt
	db           *sqlx.DB
}

func NewCommentsHandler(movieService services.MovieService, db *sqlx.DB) *CommentsHandler {
	return &CommentsHandler{
		commentService: services.CommentServiceInt,
		db:           db,
	}
}

func (h *CommentsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	// Parse movie ID from URL parameter
	vars := mux.Vars(r)
	movieID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("Invalid movie ID:%s", movieID))
		return
	}

	// Parse comment from request body
	var comment models.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("Invalid comment payload"))
		return
	}

	// Validate comment length
	if len(comment.Comment) > 500 {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("Comment length exceeds limit of 500 characters"))
		return
	}

	// Get commenter's public IP address
	comment.IPAddress = r.RemoteAddr

	// Store comment in database
	err = h.commentService.AddComment(movieID, &comment, h.db)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("Failed to store comment"))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, comment)
}

func (h *CommentsHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	// Parse movie ID from URL parameter
	vars := mux.Vars(r)
	movieID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("Invalid movie ID"))
		return
	}

	// Retrieve comments for movie from database
	comments, err := h.commentService.GetComments(movieID, h.db)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("Failed to retrieve comments"))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, comments)
}
