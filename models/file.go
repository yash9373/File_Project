package models

import (
	"time"

	"github.com/google/uuid"
)

// EncryptedFile metadata stored in DB; content is stored on disk in storage/ directory
// We DO NOT store the password or key; only salt/nonce are stored in the file content header.
// Path points to the file location on disk.
type EncryptedFile struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OwnerID   uint      `gorm:"not null" json:"owner_id"`
	Filename  string    `gorm:"size=255;not null" json:"filename"`
	Path      string    `gorm:"size=500;not null" json:"-"`
	Size      int64     `gorm:"not null" json:"size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
