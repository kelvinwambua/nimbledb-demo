package handlers

import (
	"backend/internal/auth"
	"backend/internal/models"
	"backend/internal/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	postRepo *repository.PostRepository
}

func NewPostHandler(postRepo *repository.PostRepository) *PostHandler {
	return &PostHandler{
		postRepo: postRepo,
	}
}

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	userClaims, ok := c.Locals("user").(*auth.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	post, err := h.postRepo.CreatePost(models.CreatePostParams{
		UserID:  userClaims.UserID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create post: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"post": post,
	})
}

func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	posts, err := h.postRepo.GetAllPosts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch posts: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"posts": posts,
	})
}

func (h *PostHandler) GetPost(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	post, err := h.postRepo.GetPostByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	return c.JSON(fiber.Map{
		"post": post,
	})
}

func (h *PostHandler) GetMyPosts(c *fiber.Ctx) error {
	userClaims, ok := c.Locals("user").(*auth.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	posts, err := h.postRepo.GetPostsByUserID(userClaims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch posts: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"posts": posts,
	})
}

func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	userClaims, ok := c.Locals("user").(*auth.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	isOwner, err := h.postRepo.CheckPostOwnership(id, userClaims.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}
	if !isOwner {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to update this post",
		})
	}

	var req UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	err = h.postRepo.UpdatePost(id, models.UpdatePostParams{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update post: " + err.Error(),
		})
	}

	post, _ := h.postRepo.GetPostByID(id)

	return c.JSON(fiber.Map{
		"post": post,
	})
}

func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	userClaims, ok := c.Locals("user").(*auth.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	isOwner, err := h.postRepo.CheckPostOwnership(id, userClaims.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}
	if !isOwner {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to delete this post",
		})
	}

	err = h.postRepo.DeletePost(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete post: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}
