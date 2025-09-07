package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"file_project/models"
	"file_project/repositories"

	"github.com/google/uuid"
)

type ShareLinkService struct {
	Links repositories.ShareLinkRepository
}

func NewShareLinkService(links repositories.ShareLinkRepository) *ShareLinkService {
	return &ShareLinkService{Links: links}
}

// CreateShareLink creates a share token for a file owned by the user with optional expiry/max-download limit
func (s *ShareLinkService) CreateShareLink(fileID uuid.UUID, createdBy uint, expiresInMinutes *int, maxDownloads *int) (*models.ShareLink, error) {
	token, err := generateToken(32)
	if err != nil {
		return nil, err
	}
	var expiresAt *time.Time
	if expiresInMinutes != nil && *expiresInMinutes > 0 {
		t := time.Now().Add(time.Duration(*expiresInMinutes) * time.Minute)
		expiresAt = &t
	}
	link := &models.ShareLink{
		Token:         token,
		FileID:        fileID,
		ExpiresAt:     expiresAt,
		MaxDownloads:  maxDownloads,
		Downloads:     0,
		CreatedByUser: createdBy,
	}
	if err := s.Links.Create(link); err != nil {
		return nil, err
	}
	return link, nil
}

// ValidateAndRecord checks limits and increments download counter
func (s *ShareLinkService) ValidateAndRecord(token string) (*models.ShareLink, error) {
	l, err := s.Links.FindByToken(token)
	if err != nil {
		return nil, err
	}
	// expiry
	if l.ExpiresAt != nil && time.Now().After(*l.ExpiresAt) {
		return nil, errors.New("link expired")
	}
	// max downloads
	if l.MaxDownloads != nil && l.Downloads >= *l.MaxDownloads {
		return nil, errors.New("download limit reached")
	}
	if err := s.Links.IncrementDownload(token); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *ShareLinkService) Delete(token string, createdBy uint) error {
	return s.Links.Delete(token, createdBy)
}

func generateToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// URL-safe base64 without padding
	return base64.RawURLEncoding.EncodeToString(b), nil
}
