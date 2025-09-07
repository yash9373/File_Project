package models

import (
	"time"

	"github.com/google/uuid"
)

// ShareLink represents a public share token for an encrypted file
// Token is a URL-safe random string (primary key)
type ShareLink struct {
	Token         string        `gorm:"primaryKey;size:64" json:"token"`
	FileID        uuid.UUID     `gorm:"type:uuid;not null" json:"file_id"`
	File          EncryptedFile `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	ExpiresAt     *time.Time    `json:"expires_at,omitempty"`
	MaxDownloads  *int          `json:"max_downloads,omitempty"`
	Downloads     int           `json:"downloads"`
	CreatedByUser uint          `json:"created_by_user"`
	CreatedAt     time.Time     `json:"created_at"`
}
