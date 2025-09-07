package repositories

import (
	"file_project/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileRepository interface {
	Create(file *models.EncryptedFile) error
	FindByID(id uuid.UUID, ownerID uint) (*models.EncryptedFile, error)
	ListByOwner(ownerID uint) ([]models.EncryptedFile, error)
	Delete(id uuid.UUID, ownerID uint) error
	Update(file *models.EncryptedFile) error
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(file *models.EncryptedFile) error {
	return r.db.Create(file).Error
}

func (r *fileRepository) FindByID(id uuid.UUID, ownerID uint) (*models.EncryptedFile, error) {
	var f models.EncryptedFile
	if err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *fileRepository) ListByOwner(ownerID uint) ([]models.EncryptedFile, error) {
	var list []models.EncryptedFile
	if err := r.db.Where("owner_id = ?", ownerID).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *fileRepository) Delete(id uuid.UUID, ownerID uint) error {
	return r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&models.EncryptedFile{}).Error
}

func (r *fileRepository) Update(file *models.EncryptedFile) error {
	return r.db.Save(file).Error
}
