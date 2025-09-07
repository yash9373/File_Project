package controllers

import (
	"net/url"

	"file_project/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ShareController struct {
	Shares *services.ShareLinkService
	Files  *services.FileService
}

type CreateShareRequest struct {
	FileID           string `json:"file_id"`
	ExpiresInMinutes *int   `json:"expires_in_minutes"`
	MaxDownloads     *int   `json:"max_downloads"`
}

// Create a share link for a file owned by the requester
func (sc *ShareController) Create(c *fiber.Ctx) error {
	var body CreateShareRequest
	if err := c.BodyParser(&body); err != nil || body.FileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	fileID, err := uuid.Parse(body.FileID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid file_id"})
	}
	ownerIDAny := c.Locals("user_id")
	ownerID, _ := ownerIDAny.(uint)
	// We rely on FileService to validate ownership when downloading via private endpoint; here we trust file id exists. Optionally could check.
	link, err := sc.Shares.CreateShareLink(fileID, ownerID, body.ExpiresInMinutes, body.MaxDownloads)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// return public URL path
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"token":         link.Token,
		"url":           "/share/" + url.PathEscape(link.Token) + "/download",
		"expires_at":    link.ExpiresAt,
		"max_downloads": link.MaxDownloads,
	})
}

// Public download using share token; still requires password query param
func (sc *ShareController) PublicDownload(c *fiber.Ctx) error {
	token := c.Params("token")
	pwd := c.Query("password")
	if token == "" || pwd == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token and password required"})
	}
	l, err := sc.Shares.ValidateAndRecord(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	plain, filename, err := sc.Files.DecryptAndRead(l.File.OwnerID, l.FileID, pwd)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid password or file not found"})
	}
	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", "attachment; filename=\""+url.QueryEscape(filename)+"\"")
	return c.Send(plain)
}

// Delete a share link by token (owner only)
func (sc *ShareController) Delete(c *fiber.Ctx) error {
	token := c.Params("token")
	ownerIDAny := c.Locals("user_id")
	ownerID, _ := ownerIDAny.(uint)
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token required"})
	}
	if err := sc.Shares.Delete(token, ownerID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(fiber.Map{"status": "deleted"})
}
