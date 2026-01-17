package server

import (
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/repository"
	"log"

	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App
	db          database.Service
	authHandler *handlers.AuthHandler
	postHandler *handlers.PostHandler
}

func New(dbAddr string) *FiberServer {
	db := database.New(dbAddr)
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	if err := userRepo.InitTable(); err != nil {
		log.Printf("Warning: Failed to initialize users table: %v", err)
	} else {
		log.Println("Users table ready")
	}
	if err := postRepo.InitTable(); err != nil {
		log.Printf("Warning: Failed to initialize posts table: %v", err)
	} else {
		log.Println("Posts table ready")
	}

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "backend",
			AppName:      "backend",
		}),
		db:          db,
		authHandler: handlers.NewAuthHandler(userRepo),
		postHandler: handlers.NewPostHandler(postRepo),
	}

	return server
}
