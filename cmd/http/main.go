package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	db, err := DbConn()
	errChecker(err)

	// dependency injection
	postRepository := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepository)
	postController := controller.NewPostController(postService)

	// router
	if os.Getenv("ENV") != "DEVELOPMENT" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// cors middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{fmt.Sprintf("http://localhost%s", os.Getenv("HTTP_PORT"))}
	router.Use(cors.New(config))

	// logger middleware
	router.Use(ginzap.RecoveryWithZap(logger, true))
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// pinger
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// group api
	apiGroup := router.Group("/api")
	routes.PostRoute(apiGroup, postController)

	srv := &http.Server{
		Addr:    os.Getenv("HTTP_PORT"),
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	wait := gracefulShutdown(context.Background(), logger, time.Second*5,
		func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
		func(ctx context.Context) error {
			sqlDB, _ := db.DB()
			return sqlDB.Close()
		})

	<-wait
}

func errChecker(err error) {
	if err != nil {
		panic(err)
	}
}

func DbConn() (*gorm.DB, error) {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOSTNAME"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL"),
	)

	db, err := gorm.Open(postgres.Open(conn))
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(model.Post{}, model.Tag{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

type operation func(ctx context.Context) error

func gracefulShutdown(ctx context.Context, logger *zap.Logger, timeout time.Duration, ops ...operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		logger.Info("shutting down")

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		go func() {
			<-ctx.Done()
			logger.Info("force quit the app")
			wait <- struct{}{}
		}()

		var wg sync.WaitGroup

		for key, op := range ops {
			wg.Add(1)
			go func(key int, op operation) {
				defer wg.Done()
				processName := fmt.Sprintf("process %d", key)

				if err := op(ctx); err != nil {
					logger.Error(processName, zap.Error(err), zap.Bool("success", true))
					return
				}

				logger.Info(processName, zap.Bool("success", true))
			}(key, op)
		}

		wg.Wait()
		cancel()
		close(wait)
	}()

	return wait
}
