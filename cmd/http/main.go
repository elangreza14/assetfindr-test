package main

import (
	"fmt"
	"os"

	"github.com/elangreza14/assetfindr-test/cmd/http/routes"
	"github.com/elangreza14/assetfindr-test/controller"
	"github.com/elangreza14/assetfindr-test/model"
	"github.com/elangreza14/assetfindr-test/repository"
	"github.com/elangreza14/assetfindr-test/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {

	// TODO LOGS
	// TODO GRACEFUL SHUTDOWN implement this https://medium.com/tokopedia-engineering/gracefully-shutdown-your-go-application-9e7d5c73b5ac
	// TODO MAKEFILE
	// TODO TEST CONTROLLER => http + mock service
	// TODO TEST SERVICE => service + mock db
	// TODO TEST REPOSITORY => mock sql db

	err := godotenv.Load()
	errChecker(err)

	dsn := buildPostgresConnection()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	errChecker(err)

	err = db.AutoMigrate(model.Post{}, model.Tag{})
	errChecker(err)

	postRepository := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepository)
	postController := controller.NewPostController(postService)

	router := gin.Default()
	apiGroup := router.Group("/api")
	routes.PostRoute(apiGroup, postController)

	router.Run(os.Getenv("HTTP_PORT"))
}

func errChecker(err error) {
	if err != nil {
		panic(err)
	}
}

func buildPostgresConnection() string {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOSTNAME"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL"),
	)

	return conn
}
