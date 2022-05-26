package service

import (
	"context"
	"errors"
	"go-clean-api/infrastructure"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

type FirebaseService struct {
	client *auth.Client
	logger infrastructure.Logger
	env    infrastructure.Env
}

func NewFirebaseService(client *auth.Client,
	logger infrastructure.Logger,
	env infrastructure.Env) FirebaseService {
	return FirebaseService{
		client: client,
		logger: logger,
		env:    env,
	}
}

func (fb *FirebaseService) VerifyToken(
	idToken string,
) (*auth.Token, error) {
	token, err := fb.client.VerifyIDToken(context.Background(), idToken)
	return token, err
}

// CreateUser -> creates a new user with email and password
func (fb *FirebaseService) CreateUser(
	email string,
	password string,
	verified bool,
	name string,
) (string, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		EmailVerified(verified).
		DisplayName(name).
		Disabled(false)
	u, err := fb.client.CreateUser(context.Background(), params)

	if err != nil {
		fb.logger.Zap.Error("Error while updating user", err)
		if auth.IsEmailAlreadyExists(err) {
			return "", errors.New("email already exists")
		}
		return "", err
	}

	return u.UID, err
}

func (fb *FirebaseService) VerifyEmail(email string) (string, error) {
	res, err := fb.client.EmailVerificationLink(context.Background(), email)
	fb.logger.Zap.Error("Error while updating user", err)
	return res, err
}

func (fb *FirebaseService) DeleteUser(
	uid string,
) error {

	err := fb.client.DeleteUser(context.Background(), uid)

	return err
}

func (fb *FirebaseService) GetUser(
	uid string,
) (*auth.UserRecord, error) {
	user, err := fb.client.GetUser(context.Background(), uid)
	return user, err
}

func (fb *FirebaseService) UpdateUser(
	uid string,
	email string,
	password string,
	name string,
	emailVerified bool,
) (*auth.UserRecord, error) {
	if password != "" {
		params := (&auth.UserToUpdate{}).
			Email(email).
			Password(password).
			DisplayName(name).
			EmailVerified(emailVerified)
		u, err := fb.client.UpdateUser(context.Background(), uid, params)
		if err != nil {
			return nil, err
		}

		return u, err
	} else {
		params := (&auth.UserToUpdate{}).
			Email(email).
			DisplayName(name).
			EmailVerified(emailVerified)
		u, err := fb.client.UpdateUser(context.Background(), uid, params)
		if err != nil {
			return nil, err
		}
		return u, err
	}
}

//setting claims for the user like roles and soon
func (fb *FirebaseService) SetClaim(
	uid string,
	claims gin.H,
) error {
	err := fb.client.SetCustomUserClaims(context.Background(), uid, claims)
	return err

}
