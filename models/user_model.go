package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Base
	UUID     string `jsoin:"uuid"`
	Name     string `json:"name" form:"name" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,lte=100"`
	Address  string `json:"address" binding:"required"`
	Verified bool   `json:"verified"`
	Role     string `json:"role"`
}

func (u User) TableName() string {
	return "user"
}

func (b *User) BeforeCreate(db *gorm.DB) error {
	id, err := uuid.NewRandom()
	b.ID = BINARY16(id)
	return err
}

type UserSignupInput struct {
	Email    string `json:"email" binding:"required,email,lte=100"`
	Password string `json:"password" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Name     string `json:"name" form:"name" binding:"required"`
}
