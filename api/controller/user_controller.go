package controller

import (
	"go-clean-api/api/response"
	"go-clean-api/api/service"
	"go-clean-api/infrastructure"
	"go-clean-api/models"
	"go-clean-api/utils"
	"math"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	Logger          infrastructure.Logger
	handler         infrastructure.RequestHandler
	UserService     service.UserService
	firebaseService service.FirebaseService
	awsServce       service.AWSService
	env             infrastructure.Env
	// paginate
}

func NewUserController(logger infrastructure.Logger, handler infrastructure.RequestHandler, s service.UserService, firebaseService service.FirebaseService, awsServce service.AWSService, env infrastructure.Env) UserController {
	return UserController{
		Logger:          logger,
		handler:         handler,
		UserService:     s,
		firebaseService: firebaseService,
		awsServce:       awsServce,
		env:             env,
	}
}

func (c UserController) GetAllUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page := ctx.MustGet("page").(int64)
		limit := ctx.MustGet("limit").(int64)
		var users []models.User

		total, err := c.UserService.GetAllUserPagination(&users, int(page), int(limit))
		last_page := (float64)(total / int(limit))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"done":      200,
			"data":      users,
			"page":      page,
			"last_page": math.Ceil(last_page),
		})
	}
}

func (c UserController) GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := c.UserService.SetPaginationScope(utils.Paginate(ctx)).GetAllUsers()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})

		}
		response.JSONWithPagination(ctx, http.StatusOK, users)
	}
}

func (c UserController) GetUserProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.MustGet("id").(string)
		c.Logger.Zap.Info("......userid.....", id)
		uuid, errs := models.StringToBinary16(id)
		if errs != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   errs.Error(),
			})
			return
		}
		user, err := c.UserService.GetUserById(uuid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"done":  200,
				"error": err.Error(),
			})
			ctx.Abort()
		}
		ctx.JSON(http.StatusOK, gin.H{
			"done": 200,
			"data": user,
		})
	}
}

func (c UserController) Register() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var userSignupInput models.UserSignupInput
		if err := ctx.ShouldBindJSON(&userSignupInput); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err := c.UserService.RegisterUser(userSignupInput, userSignupInput.Password, userSignupInput.Email)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "User has been registered",
		})
	}
}

func (c UserController) LoginHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type Login struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		var loginData Login
		if err := ctx.ShouldBindJSON(&loginData); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		user, err := c.UserService.GetUserByEmail(loginData.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		u, err := c.firebaseService.GetUser(user.UUID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		// isFound := ComparePassword(user.Password, loginData.Password)
		// if !isFound {
		// 	ctx.JSON(http.StatusInternalServerError, gin.H{
		// 		"success": false,
		// 		"error":   "Credential not matched",
		// 	})
		// 	return
		// }
		// token := GenerateToken(user.ID)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"userId":  u,
		})
	}
}

func (u UserController) VerifyUser() gin.HandlerFunc {
	type VerifyInput struct {
		Email string `json:"email" binding:"required"`
	}
	return func(ctx *gin.Context) {
		var data VerifyInput
		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}
		user, err := u.UserService.GetUserByEmail(data.Email)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}
		user.Verified = true
		errs := u.UserService.UpdateUser(user)
		if errs != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "user verified",
		})
	}
}
func (u UserController) UsersByEmail() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		type UserInput struct {
			Email string `json:"email" binding:"required"`
		}
		var userInput UserInput
		if err := ctx.ShouldBindJSON(&userInput); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}
		_, err := u.UserService.GetUserByEmail(userInput.Email)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "user found",
		})
	}
}

func (u UserController) ChangePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.MustGet("uid").(string)
		user, err := u.firebaseService.GetUser(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
		}
		u.Logger.Zap.Info("....fb user info", user)
		ctx.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}

func (c UserController) FileUpload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, fileHeader, _ := ctx.Request.FormFile("file")
		url, err := c.awsServce.Upload(c.env.AwsBucketName, file, fileHeader)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"file":    err.Error(),
				"bucket":  c.env.AwsBucketName,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"file":    url,
		})
	}
}
func (c UserController) SearchUser() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var userInput models.UserSearch
		page := ctx.MustGet("page").(int64)
		limit := ctx.MustGet("limit").(int64)
		if err := ctx.ShouldBindJSON(&userInput); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"success": false,
			})
		}
		users, total, err := c.UserService.SearchUser(userInput, int(page), int(limit))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"success": false,
			})
		}
		last_page := int(total) / int(limit)
		ctx.JSON(http.StatusOK, gin.H{
			"success":   true,
			"data":      users,
			"page":      page,
			"last_page": math.Ceil(float64(last_page)),
		})
	}
}

func ComparePassword(dbPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)) == nil
}

func GenerateToken(id models.BINARY16) string {
	claims := jwt.MapClaims{
		"exp":    time.Now().Add(time.Hour * 5).Unix(),
		"iat":    time.Now().Unix(),
		"userID": id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte("secretkey"))
	return t

}
