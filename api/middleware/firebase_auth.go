package middleware

import (
	"go-clean-api/api/repositories"
	"go-clean-api/api/service"
	"go-clean-api/infrastructure"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

type FirebaeAuth struct {
	firebaseService service.FirebaseService
	logger          infrastructure.Logger
	userRepository  repositories.UserRepository
}

func NewFirebaseAuth(logger infrastructure.Logger, firebaseService service.FirebaseService, userRepository repositories.UserRepository) FirebaeAuth {
	return FirebaeAuth{
		logger:          logger,
		firebaseService: firebaseService,
		userRepository:  userRepository,
	}
}

func (f FirebaeAuth) SetUp() {}

func (f FirebaeAuth) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := f.getTokenFromHeader(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			ctx.Abort()
		}
		f.logger.Zap.Info("..... signin provider...", token.Firebase.SignInProvider, token.Firebase.SignInProvider == "facebook.com")
		f.logger.Zap.Info("..... signin provider...", token.Claims)
		user, err := f.userRepository.GetUserByUuid(token.UID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			ctx.Abort()
		}
		if token.Firebase.SignInProvider == "facebook.com" && !token.Claims["email_verified"].(bool) {
			_, err = f.firebaseService.UpdateUser(token.UID, token.Claims["email"].(string), "", token.Claims["name"].(string), true)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				return
			}
		}
		f.logger.Zap.Info("..... next..............", user)
		ctx.Set("uid", token.UID)
		ctx.Set("id", user.ID.String())
		f.logger.Zap.Info("..... next..............")
		ctx.Next()
	}
}

func (f FirebaeAuth) getTokenFromHeader(c *gin.Context) (*auth.Token, error) {
	header := c.GetHeader("Authorization")
	idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
	token, err := f.firebaseService.VerifyToken(idToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}
