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
	"github.com/Eggi19/simple-social-media/utils"
	"github.com/gin-gonic/gin"
)

type RouterOpts struct {
	User     *handlers.UserHandler
	Tweet    *handlers.TweetHandler
	Comment  *handlers.CommentHandler
	Firebase *handlers.FirebaseHandler
}

func createRouter(con config.Config) *gin.Engine {
	// init config
	db, err := config.ConnectDB(con)
	if err != nil {
		log.Fatalf("error connecting to DB: %s", err.Error())
	}

	firebase, err := config.NewFirebaseRepository(con, "./firebase-service-account-key.json")
	if err != nil {
		log.Fatalf("error creating firebase client: %s", err.Error())
	}

	transaction := repositories.NewTransactor(db)

	//repository
	userRepository := repositories.NewUserRepositoryDb(&repositories.UserRepoOpt{
		Db:              db,
		FirestoreClient: firebase.FirestoreClient,
	})
	tweetRepository := repositories.NewTweetRepositoryDb(&repositories.TweetRepoOpt{
		Db:              db,
		FirebseDbClient: firebase.DbClient,
	})
	commentRepository := repositories.NewCommentRepositoryPostgres(&repositories.CommentRepoOpt{
		Db: db,
	})
	firebaseRepository := repositories.NewFirebaseRepositoryPostgres(&repositories.FirebaseRepoOpt{
		FirebaseMessagingClient: firebase.MessagingClient,
	})

	//usecase
	userUsecase := usecases.NewUserUsecaseImpl(&usecases.UserUsecaseOpts{
		HashAlgorithm:           utils.NewBCryptHasher(),
		AuthTokenProvider:       utils.NewJwtProvider(con),
		FirebaseMessagingClient: firebase.MessagingClient,
		Transactor:              transaction,
		UserRepository:          userRepository,
	})
	tweetUsecase := usecases.NewTweetUsecaseImpl(&usecases.TweetUsecaseOpts{
		TweetRepository: tweetRepository,
	})
	commentUsecase := usecases.NewCommentUsecaseImpl(&usecases.CommentUsecaseOpts{
		CommentRepository:       commentRepository,
		UserRepository:          userRepository,
		Transactor:              transaction,
		FirebaseMessagingClient: firebase.MessagingClient,
	})
	firebaseUsecase := usecases.NewFirebaseUsecaseImpl(&usecases.FirebaseUsecaseOpts{
		FirebaseRepository: firebaseRepository,
	})

	//handler
	userHandler := handlers.NewUserHandler(&handlers.UserHandlerOpts{
		UserUsecase: userUsecase,
	})
	tweetHandler := handlers.NewTweetHandler(&handlers.TweetHandlerOpts{
		TweetUsecase: tweetUsecase,
	})
	commentHandler := handlers.NewCommentHandler(&handlers.CommentHandlerOpts{
		CommentUsecase: commentUsecase,
	})
	firebaseHandler := handlers.NewFirebaseHandler(&handlers.FirebaseHandlerOpts{
		FirebaseUsecase: firebaseUsecase,
	})

	return NewRouter(con, &RouterOpts{
		User:     userHandler,
		Tweet:    tweetHandler,
		Comment:  commentHandler,
		Firebase: firebaseHandler,
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
	publicRouter := router.Group("/")
	publicRouter.POST("/register", handlers.User.RegisterUser)
	publicRouter.POST("/login", handlers.User.LoginUser)

	// private routers
	privateRouter := router.Group(("/"))
	privateRouter.Use(middlewares.JwtAuthMiddleware(config))
	privateRouter.POST("/tweet", handlers.Tweet.CreateTweet)
	privateRouter.POST("/comment", handlers.Comment.CreateComment)
	privateRouter.POST("/firebase/subscribe-topic", handlers.Firebase.SubscribeTopic)
	privateRouter.POST("/firebase/unsubscribe-topic", handlers.Firebase.UnsubscribeTopic)
	privateRouter.POST("/tweet/like", handlers.Tweet.LikeTweet)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, dtos.ErrResponse{Message: constants.EndpointNotFoundErrMsg})
	})

	return router
}
