package handlers

import (
	"backend/internal/auth"
	"backend/internal/models"
	"backend/internal/repository"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
}

func NewAuthHandler(userRepo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	user, err := h.userRepo.CreateUser(models.CreateUserParams{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Image: "https://api.dicebear.com/7.x/bottts-neutral/svg?seed=" +
			url.QueryEscape(req.Email) +
			"&backgroundColor=3b82f6,8b5cf6,ec4899,f59e0b,10b981",
	})
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "duplicate") || strings.Contains(errMsg, "unique") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already registered",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating user: " + err.Error(),
		})
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Name, user.Image, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "nimbledb-test_token",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   os.Getenv("IS_PROD") == "true",
		SameSite: "None",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return c.JSON(AuthResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Name, user.Image, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "nimbledb-test_token",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   os.Getenv("IS_PROD") == "true",
		SameSite: "None",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return c.JSON(AuthResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userClaims, ok := c.Locals("user").(*auth.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user claims",
		})
	}

	user, err := h.userRepo.GetUserById(userClaims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"image": user.Image,
			"role":  user.Role,
		},
	})
}
