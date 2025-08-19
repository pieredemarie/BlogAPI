package main

import (
	"BlogAPI/internal/app/auth"
	"BlogAPI/internal/app/post"
	"BlogAPI/pkg/jwtutils"
	"BlogAPI/pkg/storage"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	conStr := "postgres://username:password@localhost:5432/blogdb?sslmode=disable"
	db, err := storage.NewPostgresStorage(conStr)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	authHandler := &auth.Handler{Storage: db}
	postHandler := &post.PostHandler{Storage: db}

	secret := os.Getenv("JWT_SECRET")
	
	r.POST("auth/register", authHandler.RegisterHandler)
	r.POST("auth/login", authHandler.LoginHandler)

	protected := r.Group("/")
	protected.Use(jwtutils.JWTMidlleware(secret))

	protected.GET("/users/me", authHandler.GetProfile)
	protected.PUT("/users/me", authHandler.UpdateProfile)

	protected.POST("/posts", postHandler.CreatePost)
	protected.GET("/posts", postHandler.GetPosts)
	protected.GET("/posts/:id", postHandler.GetPostById)
	protected.PUT("/posts/:id", postHandler.EditPost)
	protected.DELETE("/posts/:id", postHandler.DeletePost)

	r.Run(":8080")
}