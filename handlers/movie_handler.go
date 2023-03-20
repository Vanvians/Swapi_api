package handlers

import (
	"net/http"
	"strconv"

	"swapi_api/models"
	"swapi_api/services"
	"swapi_api/utils"
	"swapi_api/db"
	"github.com/gorilla/mux"
)

type MovieHandler struct{
	service services.MovieService
}

func NewMovieHandler(service services.MovieService) *MovieHandler{
	return &MovieHandler{
		service:service,
	}
}

//returns a handler function that lists movies containing the name, opening and comment count
func (mh *MovieHandler) GetMoviesHandler() http.HandlerFunc{
	return func(w *http.ResponseWriter, r *http.request){
		movies, err := mh.service.GetMovies()
		if err != nil{
			utils.WriteError(w,err)
			return
		}

		utils.WriteJson(w, http.StatusOK, movies)
	}
}

// returs a handle function that lists characters for a movie
func (mh *CharacterHandler) GetCharactersHandler() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		movieID := mux.Vars(r)["id"]
		sort := r.Url.Query().Get("sort")
		filter := r.URL.Query().Get("filter")

		//parse limit and ofset query parameters
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)

		characters, err := mh.service.GetCharacters(movieID, sort, filter, limit,offset)
		if err != nil{
			utils.WriteError(w,err)
			return
		}
		utils.WriteJson(w, http.StatusOK, characters)
	}
}