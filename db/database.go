package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// Comment represents a single comment.
type Comment struct {
    ID        int
    PostID    int
    Content   string
    CreatedAt string
}

// db represents the database connection.
var db *sql.DB

func Connect() (*sqlx.DB, error) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %s", err)
	}

	// Create database URL string
	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	// Connect to Postgres database
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// InitializeDB initializes the database connection.
func InitializeDB(dataSourceName string) (*sql.DB, error) {
    var err error
    db, err = sql.Open("mysql", dataSourceName)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db, nil
}

// SaveComment saves a comment to the database.
func SaveComment(ctx context.Context, comment Comment) error {
    query := "INSERT INTO comments (post_id, content) VALUES (?, ?)"
    result, err := db.ExecContext(ctx, query, comment.PostID, comment.Content)
    if err != nil {
        return err
    }

    commentID, err := result.LastInsertId()
    if err != nil {
        return err
    }

    comment.ID = int(commentID)

    return nil
}

// GetComments retrieves comments from the database for a given post ID.
func GetComments(ctx context.Context, postID int) ([]Comment, error) {
    var comments []Comment

    query := "SELECT id, post_id, content, created_at FROM comments WHERE post_id = ?"
    rows, err := db.QueryContext(ctx, query, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var comment Comment
        err := rows.Scan(&comment.ID, &comment.PostID, &comment.Content, &comment.CreatedAt)
        if err != nil {
            return nil, err
        }
        comments = append(comments, comment)
    }

    return comments, nil
}

// CloseDB closes the database connection.
func CloseDB() error {
    return db.Close()
}
