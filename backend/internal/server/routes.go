package server

import (
	"backend/internal/middleware"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// DEBUG: Check environment variable
	frontendURL := os.Getenv("FRONTEND_URL")
	log.Printf("=== CORS DEBUG ===")
	log.Printf("FRONTEND_URL env var: '%s'", frontendURL)
	log.Printf("Is empty: %v", frontendURL == "")
	log.Printf("==================")

	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     frontendURL,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)
	api := s.App.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", s.authHandler.Register)
	auth.Post("/login", s.authHandler.Login)
	auth.Get("/me", middleware.AuthMiddleware, s.authHandler.GetMe)
	posts := api.Group("/posts")
	posts.Get("/", s.postHandler.GetAllPosts)
	posts.Get("/my/posts", middleware.AuthMiddleware, s.postHandler.GetMyPosts)
	posts.Get("/:id", s.postHandler.GetPost)
	posts.Post("/", middleware.AuthMiddleware, s.postHandler.CreatePost)
	posts.Put("/:id", middleware.AuthMiddleware, s.postHandler.UpdatePost)
	posts.Delete("/:id", middleware.AuthMiddleware, s.postHandler.DeletePost)

}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
