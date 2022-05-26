package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Test struct {
	Base
	Title string `json:"title" form:"title"`
}

func (u *Test) TableName() string {
	return "test"
}

func (b *Test) BeforeCreate(db *gorm.DB) error {
	id, err := uuid.NewRandom()
	b.ID = BINARY16(id)
	return err
}
