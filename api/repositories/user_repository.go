package repositories

import (
	"errors"
	"fmt"
	"go-clean-api/infrastructure"
	"go-clean-api/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	infrastructure.Database
	Logger infrastructure.Logger
}

func NewUserRepository(logger infrastructure.Logger, db infrastructure.Database) UserRepository {
	return UserRepository{
		Logger:   logger,
		Database: db,
	}
}

func (r UserRepository) WithTx(txHandler *gorm.DB) UserRepository {
	if txHandler != nil {
		r.Database.DB = txHandler
	}
	return r
}

func (r UserRepository) GetAllUser(users *[]models.User) error {
	err := r.Find(&users).Error
	if err != nil {
		return err
	}
	return nil
}

func (r UserRepository) CheckDublicateEmail(email string) (err error) {
	errs := r.Where("email = ?", email).Take(&models.User{})
	if errs.RowsAffected > 0 {
		fmt.Println("dublicate email")
		return errors.New("email already exist")
	}
	return nil
}

func (r UserRepository) CreateUser(u models.User) (err error) {
	errs := r.Create(&u).Error
	if errs != nil {
		return errs
	}
	return nil
}

func (r UserRepository) HashPassword(password string) string {
	bytePass := []byte(password)
	hPass, _ := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)
	password = string(hPass)
	return password
}
func (r UserRepository) GetUserByEmail(email string) (user models.User, err error) {
	var u models.User
	errs := r.Where("email = ?", email).First(&u).Error
	if errs != nil {
		return u, errs
	}
	return u, nil
}
func (r UserRepository) GetUserByUuid(uuid string) (user models.User, err error) {
	var u models.User
	errs := r.Where("uuid = ?", uuid).First(&u).Error
	if errs != nil {
		return u, errs
	}
	return u, nil
}
func (r UserRepository) GetUserById(id models.BINARY16) (user models.User, err error) {
	var u models.User
	errs := r.Where("id=?", id).Find(&u).Error
	if errs != nil {
		return u, errs
	}
	return u, nil
}

func (r UserRepository) UpdateUser(user models.User) (err error) {
	errs := r.Save(&user).Error
	if errs != nil {
		return errs
	}
	return nil
}
