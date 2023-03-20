package handlers

import (
	"net/http"
	"strconv"

	"Swapi_api/models"
	"Swapi_api/services"
	"Swapi_api/utils"

	"github.com/gorilla/mux"
)

type CharacterHandler struct{
	service services.CharacterService
}

func NewCharacterHandler(service services.CharacterService) *MovieHandler{
	return &MovieHandler{
		service:service,
	}
}


// returs a handle function that lists characters for a movie
func (mh *MovieHandler) GetCharactersHandler() http.HandlerFunc{
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