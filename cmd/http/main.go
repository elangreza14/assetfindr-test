package main

import (
	"fmt"
	"os"
	"time"

	"github.com/elangreza14/assetfindr-test/cmd/http/routes"
	"github.com/elangreza14/assetfindr-test/controller"
	"github.com/elangreza14/assetfindr-test/model"
	"github.com/elangreza14/assetfindr-test/repository"
	"github.com/elangreza14/assetfindr-test/service"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// TODO LOGS
	// TODO GRACEFUL SHUTDOWN implement this https://medium.com/tokopedia-engineering/gracefully-shutdown-your-go-application-9e7d5c73b5ac
	// TODO MAKEFILE
	// TODO CI CD
	// TODO Readme
	// TODO Docs swagger

	err := godotenv.Load()
	errChecker(err)

	logger := zap.NewExample(zap.IncreaseLevel(zap.InfoLevel))
	if os.Getenv("ENV") != "DEVELOPMENT" {
		logger, err = zap.NewProduction()
		errChecker(err)
	}
	defer logger.Sync()

	logger.Level()

	dsn := buildPostgresConnection()
	db, err := gorm.Open(postgres.Open(dsn))
	errChecker(err)

	err = db.AutoMigrate(model.Post{}, model.Tag{})
	errChecker(err)

	postRepository := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepository)
	postController := controller.NewPostController(postService)

	if os.Getenv("ENV") != "DEVELOPMENT" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{fmt.Sprintf("http://localhost%s", os.Getenv("HTTP_PORT"))}
	router.Use(cors.New(config))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	router.GET("/ping", func(c *gin.Context) {
		c.Writer.Header().Add("X-Request-Id", "1234-5678-9012")
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

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
