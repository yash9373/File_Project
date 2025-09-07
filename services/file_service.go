package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"file_project/models"
	"file_project/repositories"

	"github.com/google/uuid"
)

type FileService struct {
	Files repositories.FileRepository
}

func NewFileService(files repositories.FileRepository) *FileService {
	return &FileService{Files: files}
}

// SaveAndEncrypt saves the uploaded file to storage encrypted under the provided password
func (s *FileService) SaveAndEncrypt(ownerID uint, header *multipart.FileHeader, password string) (*models.EncryptedFile, error) {
	if header == nil || header.Size == 0 {
		return nil, errors.New("empty file")
	}
	// Read all bytes (for simplicity). For very large files, stream-chunk approach would be needed.
	src, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	plain, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}
	enc, err := EncryptBytes(plain, password)
	if err != nil {
		return nil, err
	}
	id := uuid.New()
	// Ensure storage directory exists
	_ = os.MkdirAll("storage", 0755)
	filenameSafe := strings.ReplaceAll(header.Filename, "..", "_")
	path := filepath.Join("storage", fmt.Sprintf("%s_%s.enc", id.String(), filenameSafe))
	if err := os.WriteFile(path, enc, 0600); err != nil {
		return nil, err
	}
	meta := &models.EncryptedFile{
		ID:       id,
		OwnerID:  ownerID,
		Filename: header.Filename,
		Path:     path,
		Size:     int64(len(enc)),
	}
	if err := s.Files.Create(meta); err != nil {
		_ = os.Remove(path)
		return nil, err
	}
	return meta, nil
}

// DecryptAndRead loads the encrypted file and decrypts with password
func (s *FileService) DecryptAndRead(ownerID uint, id uuid.UUID, password string) ([]byte, string, error) {
	meta, err := s.Files.FindByID(id, ownerID)
	if err != nil {
		return nil, "", err
	}
	enc, err := os.ReadFile(meta.Path)
	if err != nil {
		return nil, "", err
	}
	plain, err := DecryptBytes(enc, password)
	if err != nil {
		return nil, "", err
	}
	return plain, meta.Filename, nil
}

// ChangePassword re-encrypts file content with a new password
func (s *FileService) ChangePassword(ownerID uint, id uuid.UUID, oldPassword, newPassword string) error {
	meta, err := s.Files.FindByID(id, ownerID)
	if err != nil {
		return err
	}
	enc, err := os.ReadFile(meta.Path)
	if err != nil {
		return err
	}
	plain, err := DecryptBytes(enc, oldPassword)
	if err != nil {
		return err
	}
	reenc, err := EncryptBytes(plain, newPassword)
	if err != nil {
		return err
	}
	return os.WriteFile(meta.Path, reenc, 0600)
}

// Delete removes the file from disk and database
func (s *FileService) Delete(ownerID uint, id uuid.UUID) error {
	meta, err := s.Files.FindByID(id, ownerID)
	if err != nil {
		return err
	}
	_ = os.Remove(meta.Path)
	return s.Files.Delete(meta.ID, ownerID)
}

// List returns encrypted files for owner
func (s *FileService) List(ownerID uint) ([]models.EncryptedFile, error) {
	return s.Files.ListByOwner(ownerID)
}
