package controllers

import (
	"net/url"

	"file_project/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FileController struct {
	Files *services.FileService
}

func (fc *FileController) Upload(c *fiber.Ctx) error {
	password := c.FormValue("password")
	if len(password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password must be >= 6 chars"})
	}
	file, err := c.FormFile("file")
	if err != nil || file == nil || file.Size == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file is required"})
	}
	ownerIDAny := c.Locals("user_id")
	ownerID, _ := ownerIDAny.(uint)
	meta, err := fc.Files.SaveAndEncrypt(ownerID, file, password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":       meta.ID,
		"filename": meta.Filename,
		"size":     meta.Size,
	})
}

func (fc *FileController) Download(c *fiber.Ctx) error {
	idStr := c.Params("id")
	pwd := c.Query("password")
	if pwd == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password required"})
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	ownerIDAny := c.Locals("user_id")
	ownerID, _ := ownerIDAny.(uint)
	plain, filename, err := fc.Files.DecryptAndRead(ownerID, id, pwd)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid password or file not found"})
	}
	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", "attachment; filename=\""+url.QueryEscape(filename)+"\"")
	return c.Send(plain)
}

func (fc *FileController) ChangePassword(c *fiber.Ctx) error {
	type req struct{ OldPassword, NewPassword string }
	var body req
	if err := c.BodyParser(&body); err != nil || len(body.NewPassword) < 6 || body.OldPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	ownerIDAny := c.Locals("user_id")
	ownerID, _ := ownerIDAny.(uint)
	if err := fc.Files.ChangePassword(ownerID, id, body.OldPassword, body.NewPassword); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid password or file not found"})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}

func (fc *FileController) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	ownerIDAny := c.Locals("user_id")
	ownerID, _ := ownerIDAny.(uint)
	if err := fc.Files.Delete(ownerID, id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "file not found"})
	}
	return c.JSON(fiber.Map{"status": "deleted"})
}

func (fc *FileController) List(c *fiber.Ctx) error {
	ownerIDAny := c.Locals("user_id")
	ownerID, _ := ownerIDAny.(uint)
	list, err := fc.Files.List(ownerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}
