package service

import (
	"go-clean-api/api/repositories"
	"go-clean-api/constant"
	"go-clean-api/infrastructure"
	"go-clean-api/models"
)

type UserService struct {
	Logger          infrastructure.Logger
	UserRepository  repositories.UserRepository
	firebaseService FirebaseService
}

func NewUserService(logger infrastructure.Logger, r repositories.UserRepository, firebaseService FirebaseService) UserService {
	return UserService{
		Logger:          logger,
		UserRepository:  r,
		firebaseService: firebaseService,
	}
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

func (s UserService) GetUserByEmail(email string) (models.User, error) {
	return s.UserRepository.GetUserByEmail(email)
}

func (s UserService) GetUserById(id models.BINARY16) (models.User, error) {
	return s.UserRepository.GetUserById(id)
}

func (s UserService) UpdateUser(user models.User) error {
	return s.UserRepository.UpdateUser(user)
}
