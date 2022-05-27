package service

import (
	"go-clean-api/api/repositories"
	"go-clean-api/constant"
	"go-clean-api/infrastructure"
	"go-clean-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService struct {
	Logger          infrastructure.Logger
	UserRepository  repositories.UserRepository
	firebaseService FirebaseService
	paginationScope *gorm.DB
}

func NewUserService(logger infrastructure.Logger, r repositories.UserRepository, firebaseService FirebaseService) UserService {
	return UserService{
		Logger:          logger,
		UserRepository:  r,
		firebaseService: firebaseService,
	}
}

// WithTrx delegates transaction to repository database
func (s UserService) WithTrx(trxHandle *gorm.DB) UserService {
	s.UserRepository = s.UserRepository.WithTx(trxHandle)
	return s
}

// PaginationScope
func (s UserService) SetPaginationScope(scope func(*gorm.DB) *gorm.DB) UserService {
	s.paginationScope = s.UserRepository.WithTx(s.UserRepository.Scopes(scope)).DB
	return s
}

func (s UserService) RegisterUser(user models.UserSignupInput, password string, email string) (errss error) {
	err := s.UserRepository.CheckDublicateEmail(email)
	if err != nil {
		return err
	}
	UID, err := s.firebaseService.CreateUser(user.Email, password, false, user.Name)
	s.Logger.Zap.Info("fb created user id", UID)
	if err != nil {

		return err
	}

	claims := map[string]interface{}{
		constant.RoleIsAdmin: true,
		constant.RoleIsUser:  false,
	}
	cerr := s.firebaseService.SetClaim(UID, claims)
	if cerr != cerr {
		_ = s.firebaseService.DeleteUser(UID)
		return cerr
	}
	u := models.User{
		UUID:     UID,
		Name:     user.Name,
		Email:    user.Email,
		Address:  user.Address,
		Verified: false,
		Role:     constant.RoleIsUser,
	}

	errs := s.UserRepository.CreateUser(u)
	if errs != nil {
		_ = s.firebaseService.DeleteUser(UID)
		return errs
	}
	return nil
}

func (s UserService) GetAllUser(users *[]models.User) error {
	err := s.UserRepository.GetAllUser(users)
	return err
}
func (s UserService) GetAllUserPagination(users *[]models.User, page int, limit int) (int, error) {
	offset := (page - 1) * limit
	var count int64
	// s.UserRepository.Find(&models.User{}).Count(&count)
	err := s.UserRepository.Limit(limit).Offset(offset).Find(&users).Offset(-1).Limit(-1).Count(&count).Error
	return int(count), err
}

func (s UserService) GetAllUsers() (response map[string]interface{}, err error) {
	var users []models.User
	var count int64

	err = s.UserRepository.WithTx(s.paginationScope).Find(&users).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return nil, err
	}

	return gin.H{"data": users, "count": count}, nil
}

func (s UserService) GetUserByEmail(email string) (models.User, error) {
	return s.UserRepository.GetUserByEmail(email)
}

func (s UserService) GetUserById(id models.BINARY16) (models.User, error) {
	return s.UserRepository.GetUserById(id)
}

func (s UserService) UpdateUser(user models.User) error {
	return s.UserRepository.UpdateUser(user)
}

func (s UserService) SearchUser(search models.UserSearch, page int, limit int) ([]models.User, int64, error) {
	var users []models.User
	var query *gorm.DB = s.UserRepository.DB
	var count int64 = 0
	if search.Name != "" {
		query = query.Where("name LIKE ? ", "%"+search.Name+"%")
	}
	if search.Address != "" {
		query = query.Where("address LIKE ? ", "%"+search.Address+"%")
	}

	if page > 0 {
		offset := (page - 1) * limit
		query.Find(&users).Count(&count)
		query = query.Limit(limit).Offset(offset)

	}
	if search.SortBy != "" {
		query = query.Order(search.SortBy)
	}

	err := query.Find(&users).Error
	return users, count, err

}
func (s UserService) TotalDocumentCount() (int64, error) {
	var count int64
	err := s.UserRepository.Find(&models.User{}).Count(&count).Error
	return count, err
}
