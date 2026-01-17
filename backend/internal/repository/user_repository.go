package repository

import (
	"backend/internal/database"
	"backend/internal/models"
	"fmt"
	"strings"
	"time"
)

type UserRepository struct {
	db database.Service
}

func NewUserRepository(db database.Service) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) InitTable() error {
	query := "CREATE TABLE users (id INT NOT NULL, email VARCHAR(255), password VARCHAR(255), name VARCHAR(255), image VARCHAR(500), role VARCHAR(50), created_at INT, PRIMARY KEY (id))"

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

func (r *UserRepository) CreateUser(params models.CreateUserParams) (*models.User, error) {

	id := time.Now().UnixNano() / 1000000
	createdAt := time.Now().Unix()

	email := escapeString(params.Email)
	password := escapeString(params.Password)
	name := escapeString(params.Name)
	image := escapeString(params.Image)

	query := fmt.Sprintf(
		"INSERT INTO users VALUES (%d, '%s', '%s', '%s', '%s', 'user', %d)",
		id, email, password, name, image, createdAt,
	)

	err := r.db.Execute(query)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        id,
		Email:     params.Email,
		Password:  params.Password,
		Name:      params.Name,
		Image:     params.Image,
		Role:      "user",
		CreatedAt: time.Unix(createdAt, 0),
	}, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	email = escapeString(email)
	query := fmt.Sprintf("SELECT id, email, password, name, image, role, created_at FROM users WHERE email = '%s'", email)

	row, err := r.db.QueryRow(query)
	if err != nil {
		return nil, err
	}

	if len(row) < 7 {
		return nil, fmt.Errorf("invalid row data")
	}

	createdAtUnix := row[6].(int64)

	return &models.User{
		ID:        row[0].(int64),
		Email:     row[1].(string),
		Password:  row[2].(string),
		Name:      row[3].(string),
		Image:     row[4].(string),
		Role:      row[5].(string),
		CreatedAt: time.Unix(createdAtUnix, 0),
	}, nil
}

func (r *UserRepository) GetUserById(id int64) (*models.User, error) {
	query := fmt.Sprintf("SELECT id, email, password, name, image, role, created_at FROM users WHERE id = %d", id)

	row, err := r.db.QueryRow(query)
	if err != nil {
		return nil, err
	}

	if len(row) < 7 {
		return nil, fmt.Errorf("invalid row data")
	}

	createdAtUnix := row[6].(int64)

	return &models.User{
		ID:        row[0].(int64),
		Email:     row[1].(string),
		Password:  row[2].(string),
		Name:      row[3].(string),
		Image:     row[4].(string),
		Role:      row[5].(string),
		CreatedAt: time.Unix(createdAtUnix, 0),
	}, nil
}

func escapeString(s string) string {

	result := ""
	for _, c := range s {
		if c == '\'' {
			result += "''"
		} else {
			result += string(c)
		}
	}
	return result
}
