package middleware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthMiddleWare struct {
}

func NewAuthMiddlerware() AuthMiddleWare {
	return AuthMiddleWare{}
}
func (u AuthMiddleWare) SetUp() {}

func (u AuthMiddleWare) UserAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "No Authorization header found"})

		}
		token, err := ValidateToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Token", "err": err, "tokenString": authHeader})
		} else {
			claims, isTrue := token.Claims.(jwt.MapClaims)
			if !isTrue {
				c.AbortWithStatus(http.StatusUnauthorized)

			} else {
				if token.Valid {
					c.Set("userID", claims["userID"])
					fmt.Println("during authorization", claims["userID"])
				} else {
					c.AbortWithStatus(http.StatusUnauthorized)
				}

			}
		}
	}
}

func ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, err := token.Method.(*jwt.SigningMethodHMAC)
		if !err {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["zyx"])
		}
		return []byte("secretkey"), nil
	})
}
