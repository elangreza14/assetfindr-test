package routes

import (
	"github.com/elangreza14/assetfindr-test/controller"
	"github.com/gin-gonic/gin"
)

func PostRoute(route *gin.RouterGroup, postController *controller.PostController) {
	postRoutes := route.Group("/posts")
	postRoutes.GET("", postController.GetPosts())
	postRoutes.POST("", postController.CreatePost())
	postRoutes.GET("/:id", postController.GetPost())
	postRoutes.PUT("/:id", postController.UpdatePost())
	postRoutes.DELETE("/:id", postController.DeletePost())
}
