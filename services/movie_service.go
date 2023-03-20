package services

import(
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"Swapi_api/models"
	"Swapi_api/utils"
)

type MovieService struct {
	db *gorm.DB
	cache *redis.Client
}

func NewMovieService(db *gorm.DB, cache *redis.Client) *MovieService{
	return &MovieService{
		db: db,
		cache: cache,
	}
}

//fetches all movies and sorts them by release date
func (s *MovieService) GetMovies() ([]*models.movie, error){
	movies, err := s.getMoviesFromCache()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get movies from cache")
	}

	if len(movies) == 0{
		movies, err := s.getMoviesFromApi()
		if err != nil{
			return nil, errors.Wrap("failed to get movies from API")
		}

		err = s.cacheMovies(movies)
		if err != nil{
			return nil, errors.Wrap(err, "failed to cache movies")
		}
	}

	sort.Slice(movies, func(i, j int)bool{
		return movies[i].ReleaseDate.Before(movies[j].ReleaseDate)
	})

	return movies, nil
}

//Fetches a single movie by ID
func (s *MovieService) GetMovie(id string)(*models.Movie, error){
	movie, err := s.getMovieFromCache(id)
	if err != nil{
		return nil, errors.Wrap(err, "failed to get movie from cache")
	}

	if movie == nil{
		movie, err := s.getMoviesFromApi(id)
		if err != nil{
			return nil, errors.Wrap(err, "failed to get movie from API")
		}

		err = s.cacheMovie(movie)
		if err != nil{
			return nil, errors.Wrap(err, "failed to cache movie")
		}
	}
	return movie, nil
}

