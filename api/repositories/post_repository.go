package repositories

import (
	"go-clean-api/infrastructure"
	"go-clean-api/models"

	"gorm.io/gorm"
)

type PostRepository struct {
	db     infrastructure.Database
	Logger infrastructure.Logger
}

func NewPostRepository(logger infrastructure.Logger, db infrastructure.Database) PostRepository {
	return PostRepository{
		Logger: logger,
		db:     db,
	}

}

func (r PostRepository) GetAllPost(post *[]models.Post) error {
	return r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Find(&post).Error
}

func (r PostRepository) CreatePost(p models.Post) error {
	return r.db.Create(&p).Error
}
func (r PostRepository) GetPostById(p *models.Post, id models.BINARY16) error {
	return r.db.Preload("User").Where("id=?", id).Find(&p).Error
}

func (r PostRepository) UpdatePost(p *models.Post) error {
	return r.db.Omit("User").Save(&p).Error
}
func (r PostRepository) GetPostByUserId(p *[]models.Post, id models.BINARY16) error {
	return r.db.Where("user_id = ?", id).Preload("User").Find(&p).Error
}
