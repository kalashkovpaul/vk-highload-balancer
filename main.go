package main

import (
	"db_forum/app/handlers"
	"db_forum/app/middleware"
	"db_forum/app/repositories"
	"db_forum/app/usecases"
	"db_forum/pkg"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strings"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:5000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowCredentials = true

	conn, err := pgx.ParseConnectionString(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", "<your-main-server-ip>", "forum", "forum", "forum", "5432"))

	db, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     conn,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		log.Fatal("ERROR WHILE CONNECT TO DB:", err)
	}

	defer db.Close()

	// создание репозиториев
	forumRepository := repositories.MakeForumRepository(db)
	postRepository := repositories.MakePostRepository(db)
	serviceRepository := repositories.MakeServiceRepository(db)
	threadRepository := repositories.MakeThreadRepository(db)
	userRepository := repositories.MakeUserRepository(db)
	voteRepository := repositories.MakeVoteRepository(db)

	router.Use(cors.New(config))

	forumHandler := handlers.MakeForumHandler(usecases.MakeForumUseCase(forumRepository, threadRepository, userRepository))
	postHandler := handlers.MakePostHandler(usecases.MakePostUseCase(forumRepository, threadRepository, userRepository, postRepository))
	serviceHandler := handlers.MakeServiceHandler(usecases.MakeServiceUseCase(serviceRepository))
	threadHandler := handlers.MakeThreadHandler(usecases.MakeThreadUseCase(voteRepository, threadRepository, userRepository, postRepository))
	userHandler := handlers.MakeUserHandler(usecases.MakeUserUseCase(userRepository))

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Fatal(http.ListenAndServe(":9090", nil))
	}()

	router.Use(middleware.IncCounter)
	forumRoutes := router.Group(strings.Join([]string{pkg.RootRoute, pkg.ForumRoute}, ""))
	{
		forumRoutes.POST("/create", forumHandler.CreateForum)
		forumRoutes.GET("/:slug/details", forumHandler.GetForum)
		forumRoutes.POST("/:slug/create", forumHandler.CreateThread)
		forumRoutes.GET("/:slug/users", forumHandler.GetForumUsers)
		forumRoutes.GET("/:slug/:threads", forumHandler.GetForumThreads)
	}
	postRoutes := router.Group(strings.Join([]string{pkg.RootRoute, pkg.PostRoute}, ""))
	{
		postRoutes.GET("/:id/details", postHandler.GetPost)
		postRoutes.POST("/:id/details", postHandler.UpdatePost)
	}
	serviceRoutes := router.Group(strings.Join([]string{pkg.RootRoute, pkg.ServiceRoute}, ""))
	{
		serviceRoutes.POST("/clear", serviceHandler.Clear)
		serviceRoutes.GET("/status", serviceHandler.GetStatus)
	}
	threadRoutes := router.Group(strings.Join([]string{pkg.RootRoute, pkg.ThreadRoute}, ""))
	{
		threadRoutes.POST("/:slug_or_id/create", threadHandler.CreatePosts)
		threadRoutes.GET("/:slug_or_id/details", threadHandler.GetThread)
		threadRoutes.POST("/:slug_or_id/details", threadHandler.UpdateThread)
		threadRoutes.GET("/:slug_or_id/posts", threadHandler.GetThreadPosts)
		threadRoutes.POST("/:slug_or_id/vote", threadHandler.Vote)
	}
	userRoutes := router.Group(strings.Join([]string{pkg.RootRoute, pkg.UserRoute}, ""))
	{
		userRoutes.POST("/:nickname/create", userHandler.CreateUser)
		userRoutes.GET("/:nickname/profile", userHandler.GetUser)
		userRoutes.POST("/:nickname/profile", userHandler.UpdateUser)
	}

	err = router.Run(":80")
	fmt.Println("Successfull start router")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
