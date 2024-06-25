package controller

//go:generate mockgen -source $GOFILE -destination ../mock/controller/mock_$GOFILE -package $GOPACKAGE

import (
	"context"
	"errors"
	"net/http"

	"github.com/elangreza14/assetfindr-test/dto"
	"github.com/gin-gonic/gin"
)

type (
	IPostService interface {
		GetPosts(ctx context.Context) ([]dto.GetPostResponse, error)
		CreatePost(ctx context.Context, req dto.CreateOrUpdatePostRequest) error
		GetPost(ctx context.Context, ids int) (*dto.GetPostResponse, error)
		UpdatePost(ctx context.Context, req dto.CreateOrUpdatePostRequest, id int) error
		DeletePost(ctx context.Context, id int) error
	}

	PostController struct {
		postService IPostService
	}
)

func NewPostController(postService IPostService) *PostController {
	return &PostController{
		postService: postService,
	}
}

func (pc *PostController) GetPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		posts, err := pc.postService.GetPosts(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.NewBaseResponse(nil, err))
			return
		}

		c.JSON(http.StatusOK, dto.NewBaseResponse(posts, nil))
	}
}

func (pc *PostController) CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := dto.CreateOrUpdatePostRequest{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewBaseResponse(nil, err))
			return
		}

		err = pc.postService.CreatePost(c, req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.NewBaseResponse(nil, err))
			return
		}

		c.JSON(http.StatusCreated, dto.NewBaseResponse("created", nil))
	}
}

func (pc *PostController) UpdatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := dto.UriPostRequest{}
		err := c.ShouldBindUri(&uri)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewBaseResponse(nil, err))
			return
		}

		req := dto.CreateOrUpdatePostRequest{}
		err = c.ShouldBindJSON(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewBaseResponse(nil, err))
			return
		}

		err = pc.postService.UpdatePost(c, req, uri.ID)
		if err != nil {
			var errNotFound dto.ErrorNotFound
			if errors.As(err, &errNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, dto.NewBaseResponse(nil, err))
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.NewBaseResponse(nil, err))
			return
		}

		c.JSON(http.StatusOK, dto.NewBaseResponse("updated", nil))
	}
}

func (pc *PostController) DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := dto.UriPostRequest{}
		err := c.ShouldBindUri(&uri)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewBaseResponse(nil, err))
			return
		}

		err = pc.postService.DeletePost(c, uri.ID)
		if err != nil {
			var errNotFound dto.ErrorNotFound
			if errors.As(err, &errNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, dto.NewBaseResponse(nil, err))
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.NewBaseResponse(nil, err))
			return
		}

		c.JSON(http.StatusOK, dto.NewBaseResponse("deleted", nil))
	}
}

func (pc *PostController) GetPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := dto.UriPostRequest{}
		err := c.ShouldBindUri(&uri)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewBaseResponse(nil, err))
			return
		}

		post, err := pc.postService.GetPost(c, uri.ID)
		if err != nil {
			var errNotFound dto.ErrorNotFound
			if errors.As(err, &errNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, dto.NewBaseResponse(nil, err))
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.NewBaseResponse(nil, err))
			return
		}

		c.JSON(http.StatusOK, dto.NewBaseResponse(post, nil))
	}
}
