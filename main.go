package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

    "github.com/go-redis/redis/v8"

    "github.com/jcezetah/Swapi_api/cache"
    "github.com/jcezetah/Swapi_api/handlers"
    "github.com/jcezetah/Swapi_api/services"
    "github.com/jcezetah/Swapi_api/utils"
	"github.com/jcezetah/Swapi_api/db"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load env file: %s", err)
	}

	db.Connect()

    // Initialize Redis client
    redisClient := redis.NewClient(&redis.Options{
        Addr: "env",
		Password: "env",
		DB: "env",
    })

    // Initialize cache
    cache, err := cache.NewRedisCache(redisClient)

	if err != nil{

	}

    // Initialize movie service
    movieService := services.NewMovieService(cache)

	//Initialize comment service
	commentService := services.NewCommentService()

    // Initialize movie handler
    movieHandler := handlers.NewMovieHandler(movieService)

	//Initialize comment handler
	commentHandler := handlers.CommentsHandler(commentService, db)

    // Initialize router
    router := mux.NewRouter()

    // Define routes
    router.HandleFunc("/movies", movieHandler.ListMovies).Methods(http.MethodGet)
    router.HandleFunc("/movies/{id}/characters", movieHandler.ListCharacters).Methods(http.MethodGet)
	router.HandleFunc("/movies/{id}/comments", commentHandler.AddComment).Methods("POST")
	router.HandleFunc("/movies/{id}/comments", commentHandler.ListComments).Methods("GET")


    // Start server, remeber to change 8080 when you create the actual cache
	port := utils.Getenv("PORT")
	log.Printf("Starting server on port %s\n", port)
    log.Fatal(http.ListenAndServe(":8080", router))

	
}
