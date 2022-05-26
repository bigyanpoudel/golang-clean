package repositories

import (
	"errors"
	"fmt"
	"go-clean-api/infrastructure"
	"go-clean-api/models"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db     infrastructure.Database
	Logger infrastructure.Logger
}

func NewUserRepository(logger infrastructure.Logger, db infrastructure.Database) UserRepository {
	return UserRepository{
		Logger: logger,
		db:     db,
	}
}

func (r UserRepository) GetAllUser(users *[]models.User) error {
	err := r.db.Find(&users).Error
	if err != nil {
		return err
	}
	return nil
}

func (r UserRepository) CheckDublicateEmail(email string) (err error) {
	errs := r.db.Where("email = ?", email).Take(&models.User{})
	if errs.RowsAffected > 0 {
		fmt.Println("dublicate email")
		return errors.New("email already exist")
	}
	return nil
}

func (r UserRepository) CreateUser(u models.User) (err error) {
	errs := r.db.Create(&u).Error
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
	errs := r.db.Where("email = ?", email).First(&u).Error
	if errs != nil {
		return u, errs
	}
	return u, nil
}
func (r UserRepository) GetUserByUuid(uuid string) (user models.User, err error) {
	var u models.User
	errs := r.db.Where("uuid = ?", uuid).First(&u).Error
	if errs != nil {
		return u, errs
	}
	return u, nil
}
func (r UserRepository) GetUserById(id models.BINARY16) (user models.User, err error) {
	var u models.User
	errs := r.db.Where("id=?", id).Find(&u).Error
	if errs != nil {
		return u, errs
	}
	return u, nil
}

func (r UserRepository) UpdateUser(user models.User) (err error) {
	errs := r.db.Save(&user).Error
	if errs != nil {
		return errs
	}
	return nil
}
