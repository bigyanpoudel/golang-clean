package models

import (
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
