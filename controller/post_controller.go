package controller

import (
	"context"
	"net/http"

	"github.com/elangreza14/assetfindr-test/dto"
	"github.com/gin-gonic/gin"
)

type (
	IPostService interface {
		GetPosts(ctx context.Context, ids ...int) ([]dto.GetPostResponse, error)
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
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, posts)
	}
}

func (pc *PostController) CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := dto.CreateOrUpdatePostRequest{}
		err := c.ShouldBind(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		err = pc.postService.CreatePost(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusCreated, "ok")
	}
}

func (pc *PostController) UpdatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := dto.CreateOrUpdatePostRequest{}
		err := c.ShouldBind(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		uri := dto.UriPostRequest{}
		err = c.ShouldBindUri(&uri)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		err = pc.postService.UpdatePost(c, req, uri.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, req)
	}
}

func (pc *PostController) DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := dto.UriPostRequest{}
		err := c.ShouldBindUri(&uri)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		err = pc.postService.DeletePost(c, uri.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, uri)
	}
}

func (pc *PostController) GetPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := dto.UriPostRequest{}
		err := c.ShouldBindUri(&uri)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		post, err := pc.postService.GetPost(c, uri.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, post)
	}
}
