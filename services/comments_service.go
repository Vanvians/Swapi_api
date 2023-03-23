package services

import (
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/jcezetah/Swapi_api/models"
	"github.com/jcezetah/Swapi_api/db"
)

type CommentServiceInt interface {
	ListMovies() ([]models.Movie, error)
	AddComment(movieID int, comment *models.Comment, db *sqlx.DB) error
	GetComments(movieID int, db *sqlx.DB) ([]models.Comment, error)
}


func (*CommentServiceInt) AddComment(movieID int, comment *models.Comment, db *sqlx.DB) error {
	// Insert new comment into database
	_, err := db.Exec(`
		INSERT INTO comments (movie_id, text, ip_address, created_at)
		VALUES ($1, $2, $3, $4)`,
		movieID, comment.Comment, comment.IPAddress, time.Now().UTC())
	if err != nil {
		return err
	}
	return nil
}

func (*CommentServiceInt) GetComments(movieID int, db *sqlx.DB) ([]models.Comment, error) {
	// Retrieve comments for movie from database
	rows, err := db.Queryx(`
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

func AddComment(movieID string, comment *models.Comment, db *sqlx.DB) error {
    // Build the SQL statement
    stmt, err := db.Prepare("INSERT INTO comments (comment_ID, movie_id, text, IPaddress, Time) VALUES (?, ?, ?)")
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

func GetCommentCount(movieID int) (int, error) {
	// Prepare a SQL query to count the comments for the specified movie ID
    query := fmt.Sprintf("SELECT COUNT(*) FROM comments WHERE movie_id = %d", movieID)

    // Execute the query and retrieve the comment count
    var commentCount int
    err = db.QueryRow(query).Scan(&commentCount)
    if err != nil {
        return 0, fmt.Errorf("error retrieving comment count: %v", err)
    }

    return commentCount, nil
}