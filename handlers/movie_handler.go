package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/jcezetah/Swapi_api/services"
	"github.com/jcezetah/Swapi_api/utils"
)

type MovieHandler struct {
	movieService *services.MovieService
	commentService *services.CommentService
}

func NewMovieHandler(movieService *services.MovieService, commentService *services.CommentService) *MovieHandler {
	return &MovieHandler{movieService, commentService}
}

func (h *MovieHandler) ListMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.movieService.ListMovies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sort movies by release date from earliest to newest
	sortedMovies := h.movieService.SortMoviesByReleaseDate(movies)

	// Add comment count to each movie
	for i := range sortedMovies {
		commentCount, err := h.commentService.GetCommentCount(sortedMovies[i].MovieId)
		if err == nil {
			sortedMovies[i].CommentCount = commentCount
		}
	}

	// Marshal movies to JSON
	jsonMovies, err := json.Marshal(sortedMovies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonMovies)
}

func (h *MovieHandler) ListCharacters(w http.ResponseWriter, r *http.Request) {
	// Get movie ID from path parameter
	vars := mux.Vars(r)
	movieID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse query parameters
	params := r.URL.Query()
	sortBy := params.Get("sortBy")
	sortOrder := params.Get("sortOrder")
	filterByGender := params.Get("filterByGender")

	// Get list of characters for the movie
	characters, err := h.movieService.ListCharacters(movieID, sortBy, sortOrder, filterByGender)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate metadata
	totalCharacters := len(characters)
	totalHeightCm, err := h.movieService.GetTotalHeight(characters)
	if err != nil{
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("error getting total height"))
	}

	totalHeightFtIn, err := utils.ConvertCmToFtIn(totalHeightCm)
	if err != nil{
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("error converting height to ft/in"))
	}

	// Create metadata JSON object
	metadata := map[string]interface{}{
		"totalCharacters": totalCharacters,
		"totalHeightCm":   totalHeightCm,
		"totalHeightFtIn": totalHeightFtIn,
	}

	// Create response object
	response := map[string]interface{}{
		"metadata":   metadata,
		"characters": characters,
	}

	// Marshal response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
