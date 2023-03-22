package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcezetah/swapi_api/models"
	"github.com/jcezetah/swapi_api/services"
	"github.com/jcezetah/swapi_api/utils/errors"
)

type CommentsHandler struct {
	service services.CommentService
	logger  *log.Logger
}

func NewCommentsHandler(service services.CommentService, logger *log.Logger) *CommentsHandler {
	return &CommentsHandler{
		service: service,
		logger:  logger,
	}
}

func (h *CommentsHandler) handleAddComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		h.logger.Printf("error while decoding comment from request body: %v\n", err)
		errors.NewBadRequestError("invalid request body").WriteTo(w)
		return
	}

	addedComment, err := h.service.AddComment(&comment)
	if err != nil {
		errors.HandleError(err).WriteTo(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(addedComment)
}

func (h *CommentsHandler) handleListComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID, ok := vars["movie_id"]
	if !ok {
		errors.NewBadRequestError("missing movie_id parameter").WriteTo(w)
		return
	}

	comments, err := h.service.ListComments(movieID)
	if err != nil {
		errors.HandleError(err).WriteTo(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (h *CommentsHandler) handleGetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, ok := vars["comment_id"]
	if !ok {
		errors.NewBadRequestError("missing comment_id parameter").WriteTo(w)
		return
	}

	comment, err := h.service.GetComment(commentID)
	if err != nil {
		errors.HandleError(err).WriteTo(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}
