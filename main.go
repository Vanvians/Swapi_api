package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/jcezetah/Swapi_api/handlers"
	"github.com/jcezetah/Swapi_api/services"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Initialize Postgres client
	db, err := services.NewPostgresClient(os.Getenv("DB_URI"))
	if err != nil {
		log.Fatalf("Error initializing Postgres client: %v", err)
	}

	// Initialize movie service
	movieService := services.NewMovieService(redisClient)

	// Initialize comment service
	commentService := services.NewCommentService(db)

	// Initialize character service
	characterService := services.NewCharacterService()

	// Initialize router
	r := mux.NewRouter()

	// Register movie handlers
	movieHandler := handlers.NewMovieHandler(movieService)
	r.HandleFunc("/movies", movieHandler.ListMovies).Methods("GET")
	r.HandleFunc("/movies/comments", movieHandler.AddComment).Methods("POST")
	r.HandleFunc("/movies/{id}/comments", movieHandler.ListComments).Methods("GET")
	r.HandleFunc("/movies/{id}/characters", movieHandler.ListCharacters).Methods("GET")

	// Register character handlers
	characterHandler := handlers.NewCharacterHandler(characterService)
	r.HandleFunc("/characters", characterHandler.ListCharacters).Methods("GET")

	// Serve API
	fmt.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
