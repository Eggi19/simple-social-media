package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Eggi19/simple-social-media/config"
	"github.com/Eggi19/simple-social-media/constants"
	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/handlers"
	"github.com/Eggi19/simple-social-media/middlewares"
	"github.com/Eggi19/simple-social-media/repositories"
	"github.com/Eggi19/simple-social-media/usecases"
	"github.com/gin-gonic/gin"
)

type RouterOpts struct {
	User *handlers.UserHandler
}

func createRouter(con config.Config) *gin.Engine {
	db, err := config.ConnectDB(con)
	if err != nil {
		log.Fatalf("error connecting to DB: %s", err.Error())
	}

	//repository
	userRepository := repositories.NewUserRepositoryPostgres(&repositories.UserRepoOpt{Db: db})

	//usecase
	userUsecase := usecases.NewUserUsecaseImpl(&usecases.UserUsecaseOpts{UserRepository: userRepository})

	//handler
	userHandler := handlers.NewUserHandler(&handlers.UserHandlerOpts{UserUsecase: userUsecase})

	return NewRouter(con, &RouterOpts{
		User: userHandler,
	})
}

func Init() {
	config, err := config.ConfigInit()
	if err != nil {
		log.Fatalf("error getting env: %s", err.Error())
	}

	router := createRouter(config)

	srv := http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%s", config.Port),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 3)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go func() {
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown: ", err)
		}
	}()

	<-ctx.Done()
	log.Println("Server exiting")
}

func NewRouter(config config.Config, handlers *RouterOpts) *gin.Engine {
	router := gin.Default()

	router.ContextWithFallback = true

	router.Use(middlewares.ErrorHandling)

	// public routers
	// publicRouter := router.Group("/")

	// private routers
	privateRouter := router.Group(("/"))
	privateRouter.Use(middlewares.JwtAuthMiddleware(config))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, dtos.ErrResponse{Message: constants.EndpointNotFoundErrMsg})
	})

	return router
}
