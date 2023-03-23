package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/jcezetah/Swapi_api/services"
	"github.com/jcezetah/Swapi_api/utils"
)

type MovieHandler struct {
	movieService *services.MovieService
}

func NewMovieHandler(movieService *services.MovieService) *MovieHandler {
	return &MovieHandler{movieService}
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
		commentCount, err := h.movieService.GetCommentCount(sortedMovies[i].EpisodeId)
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
	totalHeightCm := h.movieService.GetTotalHeight(characters)
	totalHeightFtIn, err := utils.ConvertCmToFtIn(totalHeightCm)

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

func getMovieHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    movieID := vars["movieID"]

    // check if the movie exists
    movie, ok := movies[movieID]
    if !ok {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    // return the movie information
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(movie)
}