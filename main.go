package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/go-redis/redis/v8"

	"github.com/jcezetah/Swapi_api/cache"
	"github.com/jcezetah/Swapi_api/db"
	"github.com/jcezetah/Swapi_api/handlers"
	"github.com/jcezetah/Swapi_api/services"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatalf("Failed to load env file: %s", err)
	}

	db.Connect()

    // Initialize Redis client
    redisClient := redis.NewClient(&redis.Options{
        Addr: os.Getenv("REDIS_HOST"),
        DB: 1,
    })

    // Initialize cache
    cache, err := cache.NewRedisCache(redisClient)

	if err != nil{
        log.Fatalf("failed to create RedisCache: %s", err)
	}

    // Initialize movie service
    movieService := services.NewMovieService(cache)

	//Initialize comment service
	commentService := services.NewCommentService(cache, db.DB)

    // Initialize movie handler
    movieHandler := handlers.NewMovieHandler(movieService, commentService)

	//Initialize comment handler
	commentHandler := handlers.NewCommentsHandler(commentService, db.DB)

    // Initialize router
    router := mux.NewRouter()

    // Define routes
    router.HandleFunc("/movies", movieHandler.ListMovies).Methods(http.MethodGet)
    router.HandleFunc("/movies/{id}/characters", movieHandler.ListCharacters).Methods(http.MethodGet)
	router.HandleFunc("/movies/{id}/comments", commentHandler.AddComment).Methods("POST")
	router.HandleFunc("/movies/{id}/comments", commentHandler.ListComments).Methods("GET")


    // Start server, remeber to change 8080 when you create the actual cache
	port := os.Getenv("PORT")
	log.Printf("Starting server on port %s\n", port)
    log.Fatal(http.ListenAndServe(":8080", router))

	
}
