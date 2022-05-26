package controller

import (
	"go-clean-api/api/service"
	"go-clean-api/infrastructure"
	"go-clean-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	Logger      infrastructure.Logger
	PostService service.PostService
	AWSService  service.AWSService
	Env         infrastructure.Env
}

func NewPostController(logger infrastructure.Logger, s service.PostService, AWSService service.AWSService, Env infrastructure.Env) PostController {
	return PostController{
		Logger:      logger,
		PostService: s,
		AWSService:  AWSService,
		Env:         Env,
	}
}

func (c PostController) GetAllPost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var posts []models.Post
		err := c.PostService.GetAllPost(&posts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    posts,
		})
	}
}

func (c PostController) CreatePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type CreatePost struct {
			models.Base
			Title       string          `json:"title" form:"title"`
			Description string          `json:"description" form:"description"`
			UserID      models.BINARY16 `json:"user_id"`
		}
		var post CreatePost
		if err := ctx.ShouldBindJSON(&post); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Post Unable to create",
				"Error":   err.Error(),
			})
			return
		}

		var p models.Post
		p.Title = post.Title
		p.Description = post.Description
		p.UserID = post.UserID

		createErr := c.PostService.CreatePost(p)
		if createErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   createErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Post has been created",
		})

	}
}

func (c PostController) UploadPostImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		uuid, errs := models.StringToBinary16(id)
		if errs != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   errs.Error(),
			})
			return
		}

		var post models.Post
		err := c.PostService.GetPostById(&post, uuid)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		file, fileHeader, _ := ctx.Request.FormFile("file")
		url, err := c.AWSService.Upload(c.Env.AwsBucketName, file, fileHeader)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}
		post.Image = url
		c.Logger.Zap.Info("post.......", post)
		updateError := c.PostService.UpdatePost(&post)
		if updateError != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Image has been uploaded",
		})

	}
}

func (c PostController) GetPostById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		uuid, errs := models.StringToBinary16(id)
		if errs != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   errs.Error(),
			})
			return
		}
		var post models.Post
		err := c.PostService.GetPostById(&post, uuid)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    post,
		})
	}
}

func (c PostController) GetPostByUserId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.MustGet("id").(string)
		uuid, errs := models.StringToBinary16(id)
		if errs != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   errs.Error(),
			})
			return
		}
		var post []models.Post
		err := c.PostService.GetPostByUserId(&post, uuid)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"data":    post,
		})
	}
}
