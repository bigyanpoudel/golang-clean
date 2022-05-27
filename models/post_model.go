package models

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	Base
	Title       string   `json:"title" form:"title"`
	Description string   `json:"description" form:"description"`
	Image       string   `json:"image" form:"image"`
	UserID      BINARY16 `json:"user_id"`
	User        User     `gorm:"foreignKey:UserID;references:ID" json:"user"`
}

func (u *Post) TableName() string {
	return "post"
}

func (b *Post) BeforeCreate(db *gorm.DB) error {
	id, err := uuid.NewRandom()
	b.ID = BINARY16(id)
	return err
}

type PostWithUser []Post

func (p PostWithUser) ToMapPosts() []gin.H {
	data := []gin.H{}
	for _, p := range p {

		data = append(data, gin.H{
			"id":          p.ID,
			"image":       p.Image,
			"title":       p.Title,
			"description": p.Description,
			"created_at":  p.CreatedAt,
			"updated_at":  p.UpdatedAt,
			"user": gin.H{
				"id":   p.User.ID,
				"name": p.User.Name,
			},
		})
	}
	fmt.Println("Data", data)
	return data
}
