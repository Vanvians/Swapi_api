package handlers

import (
	"net/http"
	"strconv"

	"Swapi_api/models"
	"Swapi_api/services"
	"Swapi_api/utils"
t66
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

