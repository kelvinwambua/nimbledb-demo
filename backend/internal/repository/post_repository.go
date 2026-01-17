package repository

import (
	"backend/internal/database"
	"backend/internal/models"
	"fmt"
	"strings"
	"time"
)

type PostRepository struct {
	db database.Service
}

func NewPostRepository(db database.Service) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) InitTable() error {
	query := "CREATE TABLE posts (id INT NOT NULL, user_id INT NOT NULL, title VARCHAR(255), content VARCHAR(5000), created_at INT, updated_at INT, PRIMARY KEY (id))"

	err := r.db.Execute(query)
	if err != nil {
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "already exists") ||
			strings.Contains(errMsg, "duplicate") ||
			strings.Contains(errMsg, "exists") {
			return nil
		}
		return err
	}
	return nil
}

func (r *PostRepository) CreatePost(params models.CreatePostParams) (*models.Post, error) {
	id := time.Now().UnixNano() / 1000000
	now := time.Now().Unix()

	title := escapeString(params.Title)
	content := escapeString(params.Content)

	query := fmt.Sprintf(
		"INSERT INTO posts VALUES (%d, %d, '%s', '%s', %d, %d)",
		id, params.UserID, title, content, now, now,
	)

	err := r.db.Execute(query)
	if err != nil {
		return nil, err
	}

	return &models.Post{
		ID:        id,
		UserID:    params.UserID,
		Title:     params.Title,
		Content:   params.Content,
		CreatedAt: time.Unix(now, 0),
		UpdatedAt: time.Unix(now, 0),
	}, nil
}

func (r *PostRepository) GetAllPosts() ([]models.Post, error) {

	query := "SELECT id, user_id, title, content, created_at, updated_at FROM posts ORDER BY created_at DESC"

	cols, rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	_ = cols

	posts := make([]models.Post, 0, len(rows))
	for _, row := range rows {
		if len(row) < 6 {
			continue
		}

		post := models.Post{
			ID:        row[0].(int64),
			UserID:    row[1].(int64),
			Title:     row[2].(string),
			Content:   row[3].(string),
			CreatedAt: time.Unix(row[4].(int64), 0),
			UpdatedAt: time.Unix(row[5].(int64), 0),
		}

		userQuery := fmt.Sprintf("SELECT name, email, image FROM users WHERE id = %d", post.UserID)
		userRow, err := r.db.QueryRow(userQuery)
		if err == nil && len(userRow) >= 3 {
			post.AuthorName = userRow[0].(string)
			post.AuthorEmail = userRow[1].(string)
			post.AuthorImage = userRow[2].(string)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) GetPostByID(id int64) (*models.Post, error) {
	query := fmt.Sprintf("SELECT id, user_id, title, content, created_at, updated_at FROM posts WHERE id = %d", id)

	row, err := r.db.QueryRow(query)
	if err != nil {
		return nil, err
	}

	if len(row) < 6 {
		return nil, fmt.Errorf("invalid row data")
	}

	post := &models.Post{
		ID:        row[0].(int64),
		UserID:    row[1].(int64),
		Title:     row[2].(string),
		Content:   row[3].(string),
		CreatedAt: time.Unix(row[4].(int64), 0),
		UpdatedAt: time.Unix(row[5].(int64), 0),
	}

	userQuery := fmt.Sprintf("SELECT name, email, image FROM users WHERE id = %d", post.UserID)
	userRow, err := r.db.QueryRow(userQuery)
	if err == nil && len(userRow) >= 3 {
		post.AuthorName = userRow[0].(string)
		post.AuthorEmail = userRow[1].(string)
		post.AuthorImage = userRow[2].(string)
	}

	return post, nil
}

func (r *PostRepository) GetPostsByUserID(userID int64) ([]models.Post, error) {
	query := fmt.Sprintf("SELECT id, user_id, title, content, created_at, updated_at FROM posts WHERE user_id = %d ORDER BY created_at DESC", userID)

	cols, rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	_ = cols

	posts := make([]models.Post, 0, len(rows))
	for _, row := range rows {
		if len(row) < 6 {
			continue
		}

		posts = append(posts, models.Post{
			ID:        row[0].(int64),
			UserID:    row[1].(int64),
			Title:     row[2].(string),
			Content:   row[3].(string),
			CreatedAt: time.Unix(row[4].(int64), 0),
			UpdatedAt: time.Unix(row[5].(int64), 0),
		})
	}

	return posts, nil
}

func (r *PostRepository) UpdatePost(id int64, params models.UpdatePostParams) error {
	now := time.Now().Unix()

	title := escapeString(params.Title)
	content := escapeString(params.Content)

	query := fmt.Sprintf(
		"UPDATE posts SET title = '%s', content = '%s', updated_at = %d WHERE id = %d",
		title, content, now, id,
	)

	return r.db.Execute(query)
}

func (r *PostRepository) DeletePost(id int64) error {
	query := fmt.Sprintf("DELETE FROM posts WHERE id = %d", id)
	return r.db.Execute(query)
}

func (r *PostRepository) CheckPostOwnership(postID, userID int64) (bool, error) {
	query := fmt.Sprintf("SELECT user_id FROM posts WHERE id = %d", postID)

	row, err := r.db.QueryRow(query)
	if err != nil {
		return false, err
	}

	if len(row) < 1 {
		return false, fmt.Errorf("post not found")
	}

	ownerID := row[0].(int64)
	return ownerID == userID, nil
}
