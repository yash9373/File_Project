package controllers

import (
	"strings"

	"file_project/models"
	"file_project/repositories"
	"file_project/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	Users repositories.UserRepository
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Basic validation helpers
func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".") && len(email) >= 5
}

func validatePassword(p string) bool {
	return len(p) >= 6
}

func (a *AuthController) Register(c *fiber.Ctx) error {
	var body RegisterRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	}
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	body.Name = strings.TrimSpace(body.Name)

	if body.Name == "" || !validateEmail(body.Email) || !validatePassword(body.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed: name/email/password"})
	}

	// Check existing user
	count, _ := a.Users.CountByEmail(body.Email)
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already registered"})
	}

	hash, err := utils.HashPassword(body.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not hash password"})
	}
	user := models.User{Name: body.Name, Email: body.Email, Password: hash}
	if err := a.Users.Create(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": token, "user": fiber.Map{"id": user.ID, "name": user.Name, "email": user.Email}})
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	var body LoginRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	}
	email := strings.TrimSpace(strings.ToLower(body.Email))
	if !validateEmail(email) || !validatePassword(body.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed: email/password"})
	}

	user, err := a.Users.FindByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	if !utils.CheckPasswordHash(user.Password, body.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}
	return c.JSON(fiber.Map{"token": token})
}

func (a *AuthController) Me(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"name":    c.Locals("name"),
		"user_id": c.Locals("user_id"),
		"email":   c.Locals("email"),
	})
}
