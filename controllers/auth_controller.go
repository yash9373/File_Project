package controllers

import (
	"strings"

	"file_project/database"
	"file_project/models"
	"file_project/utils"

	"github.com/gofiber/fiber/v2"
)

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

func Register(c *fiber.Ctx) error {
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
	var count int64
	database.DB.Model(&models.User{}).Where("email = ?", body.Email).Count(&count)
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already registered"})
	}

	hash, err := utils.HashPassword(body.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not hash password"})
	}
	user := models.User{Name: body.Name, Email: body.Email, Password: hash}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": token, "user": fiber.Map{"id": user.ID, "name": user.Name, "email": user.Email}})
}

func Login(c *fiber.Ctx) error {
	var body LoginRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	}
	email := strings.TrimSpace(strings.ToLower(body.Email))
	if !validateEmail(email) || !validatePassword(body.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed: email/password"})
	}

	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	if !utils.CheckPasswordHash(user.Password, body.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}
	return c.JSON(fiber.Map{"token": token})
}

func Me(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"user_id": c.Locals("user_id"),
		"email":   c.Locals("email"),
	})
}
