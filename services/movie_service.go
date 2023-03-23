package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/jcezetah/Swapi_api/cache"
	"github.com/jcezetah/Swapi_api/utils"
	"github.com/jcezetah/Swapi_api/models"
)



type MovieService struct {
	cache cache.RedisCache
}

func NewMovieService(cache cache.RedisCache) *MovieService {
	return &MovieService{cache}
}

func (s *MovieService) ListMovies() ([]models.Movie, error) {
	// Check if movies are already cached
	movies, err := s.cache.Get("movies")
	if err == nil {
		// Movies found in cache, unmarshal and return
		var cachedMovies []models.Movie
		err = json.Unmarshal([]byte(movies), &cachedMovies)
		if err == nil {
			return cachedMovies, nil
		}
	}

	// Movies not found in cache, fetch from swapi.dev
	resp, err := http.Get(utils.BaseUrlFilms)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Results []struct {
			EpisodeID   int    `json:"episode_id"`
			Title       string `json:"title"`
			OpeningCrawl string `json:"opening_crawl"`
			ReleaseDate string `json:"release_date"`
		} `json:"results"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	// Convert response to Movie objects
	movies = make([]models.Movie, len(response.Results))
	for i, result := range response.Results {
		movies[i] = .Movie{
			EpisodeId:   result.EpisodeID,
			Title:       result.Title,
			OpeningCrawl: result.OpeningCrawl,
			ReleaseDate: result.ReleaseDate,
		}
	}

	// Cache movies
	jsonMovies, err := json.Marshal(movies)
	if err != nil {
		return nil, err
	}
	err = s.cache.Set("movies", string(jsonMovies), utils.CacheExpiration)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (s *MovieService) SortMoviesByReleaseDate(movies []models.Movie) []Movie {
	sort.Slice(movies, func(i, j int) bool {
		return movies[i].ReleaseDate < movies[j].ReleaseDate
	})
	return movies
}

func (s *MovieService) ListCharacters(movieID int, sortBy string, sortOrder string, filterByGender string) ([]Character, error) {
    // Get the list of characters for the specified movie ID
    characters, err := GetCharacters(movieID)
    if err != nil {
        return nil, err
    }
    
    // Apply the sorting and filtering options
    switch sortBy {
    case "name":
        sort.Slice(characters, func(i, j int) bool {
            if sortOrder == "desc" {
                return characters[i].Name > characters[j].Name
            }
            return characters[i].Name < characters[j].Name
        })
    case "height":
        sort.Slice(characters, func(i, j int) bool {
            if sortOrder == "desc" {
                return characters[i].HeightCm > characters[j].HeightCm
            }
            return characters[i].HeightCm < characters[j].HeightCm
        })
    }
    if filterByGender != "" {
        var filteredCharacters []models.Character
        for _, c := range characters {
            if c.Gender == filterByGender {
                filteredCharacters = append(filteredCharacters, c)
            }
        }
        characters = filteredCharacters
    }
    
    return characters, nil
}

func GetTotalHeight(characters []models.Character) (int, error) {
    // Calculate the total height of the specified list of characters
    totalHeightCm := 0
    for _, c := range characters {
        totalHeightCm += c.HeightCm
    }
    
    return totalHeightCm, nil
}

