package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/jcezetah/Swapi_api/cache"
	"github.com/jcezetah/Swapi_api/models"
	"github.com/jcezetah/Swapi_api/utils"
)



type MovieService struct {
	cache *cache.RedisCache
}

func NewMovieService(cache *cache.RedisCache) *MovieService {
	return &MovieService{cache}
}

func (s *MovieService) SortMoviesByReleaseDate(movies []models.Movie) []models.Movie {
	sort.Slice(movies, func(i, j int) bool {
		return movies[i].ReleaseDate < movies[j].ReleaseDate
	})
	return movies
}

func (s *MovieService) ListCharacters(movieID int, sortBy string, sortOrder string, filterByGender string) ([]models.Character, error) {
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

func (s *MovieService) GetTotalHeight(characters []models.Character) (int, error) {
    // Calculate the total height of the specified list of characters
    totalHeightCm := 0
    for _, c := range characters {
        totalHeightCm += c.HeightCm
    }
    
    return totalHeightCm, nil
}

func GetCharacters(movieID int) ([]models.Character, error) {
    // Construct URL for characters endpoint
    url := fmt.Sprintf("%s/%d/", utils.BaseUrlPeaople, movieID)

    // Fetch character data from API
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Decode response JSON into struct
    var movie struct {
        Characters []string `json:"characters"`
    }
    err = json.NewDecoder(resp.Body).Decode(&movie)
    if err != nil {
        return nil, err
    }

    // Initialize slice of Character objects
    characters := make([]models.Character, len(movie.Characters))

    // Fetch character data for each character URL
    for i, characterUrl := range movie.Characters {
        resp, err := http.Get(characterUrl)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()

        // Decode character JSON into struct
        var character models.Character
        err = json.NewDecoder(resp.Body).Decode(&character)
        if err != nil {
            return nil, err
        }

        // Add character to list of characters
        characters[i] = character
    }

    return characters, nil
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
	moviesObjects := make([]models.Movie, len(response.Results))
	for i, result := range response.Results {
		moviesObjects[i] = models.Movie{
			MovieId:   result.EpisodeID,
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
	err = s.cache.Set(context.Background(),"movies", string(jsonMovies), utils.CacheExpiration)
	if err != nil {
		return nil, err
	}

	return moviesObjects, nil
}







