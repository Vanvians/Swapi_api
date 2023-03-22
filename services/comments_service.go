package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jcezetah/Swapi_api/models"
	"github.com/jcezetah/Swapi_api/utils"
)

type commentService struct {
	cache *utils.Cache
	db    *models.DB
}

func NewCommentService(repo models.CommentRepository) commentService {
	return &commentService{
		repo: repo,
	}
}

func (s *commentsService) getCommentCount(movieID int) (int, error) {
	key := "comment_count_" + strconv.Itoa(movieID)

	// Try to retrieve the comment count from the cache
	if count, ok := s.cache.Get(key); ok {
		return count.(int), nil
	}

	// If the comment count is not in the cache, retrieve it from the database
	count, err := s.db.GetCommentCount(movieID)
	if err != nil {
		return 0, err
	}

	// Cache the comment count for future requests
	s.cache.Set(key, count)

	return count, nil
}

func (s *commentsService) addComment(movieID int, text string, ipAddress string) error {
	comment := models.Comment{
		MovieID:    movieID,
		Text:       text,
		IPAddress: ipAddress,
	}

	err := s.db.InsertComment(comment)
	if err != nil {
		return utils.NewError(http.StatusInternalServerError, "Error inserting comment into database")
	}

	// Invalidate the comment count cache for the movie
	key := "comment_count_" + strconv.Itoa(movieID)
	s.cache.Delete(key)

	return nil
}

func (s *commentsService) listComments(movieID int) ([]models.Comment, error) {
	comments, err := s.db.ListComments(movieID)
	if err != nil {
		return nil, utils.NewError(http.StatusInternalServerError, "Error retrieving comments from database")
	}

	// Reverse the order of comments to display them in reverse chronological order
	for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
		comments[i], comments[j] = comments[j], comments[i]
	}

	return comments, nil
}

func (s *commentsService) handleAddComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieID, err := strconv.Atoi(params["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid movie ID")
		return
	}

	var requestBody struct {
		Text string `json:"text"`
	}
	err = utils.ParseJSONRequest(r, &requestBody)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ipAddress := r.Header.Get("X-Real-Ip")
	err = s.addComment(movieID, requestBody.Text, ipAddress)
	if err != nil {
		utils.RespondWithError(w, err.StatusCode, err.Message)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Comment added successfully"})
}

// handleListComments handles the GET request to list all comments for a movie
func handleListComments(w http.ResponseWriter, r *http.Request) {
	// Get the movie ID from the URL parameters
	vars := mux.Vars(r)
	movieID := vars["id"]

	// Get the comments for the movie from the database
	comments, err := GetCommentsByMovieID(movieID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to get comments", err)
		return
	}

	// Reverse the order of comments
	reverseComments(comments)

	// Return the comments as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to encode comments", err)
		return
	}
}

// handleAddComment handles the POST request to add a comment for a movie
func handleAddComment(w http.ResponseWriter, r *http.Request) {
	// Get the movie ID from the URL parameters
	vars := mux.Vars(r)
	movieID := vars["id"]

	// Get the commenter's IP address
	ip := getIPAddress(r)

	// Get the current date and time
	now := time.Now().UTC()

	// Decode the request body into a Comment struct
	var comment Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		handleError(w, http.StatusBadRequest, "Failed to decode request body", err)
		return
	}

	// Check if the comment is too long
	if len(comment.Text) > 500 {
		handleError(w, http.StatusBadRequest, "Comment is too long", ErrCommentTooLong)
		return
	}

	// Add the comment to the database
	if err := AddComment(movieID, comment.Text, ip, now); err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to add comment", err)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated)
}

// handleError sends an error response with the given status code, message, and error
func handleError(w http.ResponseWriter, statusCode int, message string, err error) {
	log.Printf("%s: %s", message, err.Error())
	response := ErrorResponse{
		Message: message,
		Error:   err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode error response: %s", err.Error())
	}
}

// getIPAddress returns the IP address of the requester
func getIPAddress(r *http.Request) string {
	// Check for X-Forwarded-For header first (used by Heroku)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ", ")[0]
	}
	// Fall back to remote address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("Failed to get IP address: %s", err.Error())
	}
	return ip
}

// reverseComments reverses the order of comments in place
func reverseComments(comments []Comment) {
	for i := 0; i < len(comments)/2; i++ {
		j := len(comments) - i - 1
		comments[i], comments[j] = comments[j], comments[i]
	}
}
