package service

import (
	"go-clean-api/api/repositories"
	"go-clean-api/infrastructure"
	"go-clean-api/models"
)

type PostService struct {
	Logger         infrastructure.Logger
	PostRepository repositories.PostRepository
}

func NewPostService(logger infrastructure.Logger, r repositories.PostRepository) PostService {
	return PostService{
		Logger:         logger,
		PostRepository: r,
	}
}

func (s PostService) GetAllPost(post *[]models.Post) error {
	return s.PostRepository.GetAllPost(post)
}

func (s PostService) CreatePost(p models.Post) error {
	return s.PostRepository.CreatePost(p)
}
func (s PostService) GetPostById(p *models.Post, id models.BINARY16) error {
	return s.PostRepository.GetPostById(p, id)
}
func (s PostService) UpdatePost(p *models.Post) error {
	return s.PostRepository.Save(&p).Error
}

func (s PostService) GetPostByUserId(p *[]models.Post, id models.BINARY16) error {
	return s.PostRepository.GetPostByUserId(p, id)
}
