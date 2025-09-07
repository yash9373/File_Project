package repositories

import (
	"file_project/models"

	"gorm.io/gorm"
)

type ShareLinkRepository interface {
	Create(link *models.ShareLink) error
	FindByToken(token string) (*models.ShareLink, error)
	IncrementDownload(token string) error
	Delete(token string, createdBy uint) error
}

type shareLinkRepository struct {
	db *gorm.DB
}

func NewShareLinkRepository(db *gorm.DB) ShareLinkRepository {
	return &shareLinkRepository{db: db}
}

func (r *shareLinkRepository) Create(link *models.ShareLink) error {
	return r.db.Create(link).Error
}

func (r *shareLinkRepository) FindByToken(token string) (*models.ShareLink, error) {
	var l models.ShareLink
	if err := r.db.Preload("File").Where("token = ?", token).First(&l).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *shareLinkRepository) IncrementDownload(token string) error {
	return r.db.Model(&models.ShareLink{}).Where("token = ?", token).UpdateColumn("downloads", gorm.Expr("downloads + 1")).Error
}

func (r *shareLinkRepository) Delete(token string, createdBy uint) error {
	return r.db.Where("token = ? AND created_by_user = ?", token, createdBy).Delete(&models.ShareLink{}).Error
}
