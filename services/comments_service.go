package services

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/jcezetah/Swapi_api/models"
	"github.com/jcezetah/Swapi_api/cache"
)

type CommentServiceInt interface {
	ListMovies() ([]models.Movie, error)
	AddComment(movieID int, comment *models.Comment, db *sqlx.DB) error
	GetComments(movieID int, db *sqlx.DB) ([]models.Comment, error)
}

type CommentService struct {
    cache *cache.RedisCache
	db *sqlx.DB
}

func NewCommentService(cache *cache.RedisCache, db *sqlx.DB) *CommentService{
	return &CommentService{cache, db}
}

func (s *CommentService) GetComments(movieID int) ([]models.Comment, error) {
	// Retrieve comments for movie from database
	rows, err := s.db.Queryx(`
		SELECT text, ip_address, created_at
		FROM comments
		WHERE movie_id = $1
		ORDER BY created_at DESC`,
		movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Transform database rows into models.Comment format
	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err = rows.Scan(&comment.Comment, &comment.IPAddress, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (s *CommentService) AddComment(movieID int, comment *models.Comment) error {
    // Build the SQL statement
    stmt, err := s.db.Prepare("INSERT INTO comments (comment_ID, movie_id, text, IPaddress, Time) VALUES (?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    // Execute the SQL statement
    _, err = stmt.Exec(comment.ID, movieID, comment.Comment, comment.IPAddress, comment.CreatedAt)
    if err != nil {
        return err
    }

    return nil
}

func (s *CommentService)GetCommentCount(movieID int) (int, error) {
	// Prepare a SQL query to count the comments for the specified movie ID
    query := fmt.Sprintf("SELECT COUNT(*) FROM comments WHERE movie_id = %d", movieID)

    // Execute the query and retrieve the comment count
    var commentCount int
    err := s.db.QueryRow(query).Scan(&commentCount)
    if err != nil {
        return 0, fmt.Errorf("error retrieving comment count: %v", err)
    }

    return commentCount, nil
}