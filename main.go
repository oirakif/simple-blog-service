package main

import (
	"database/sql"
	"fmt"
	"net/http"
	authDomain "oirakif/simple-blog-service/pkg/auth/domain"
	authHTTPHandler "oirakif/simple-blog-service/pkg/auth/handler"
	blogPostDomain "oirakif/simple-blog-service/pkg/blog-post/domain"
	blogPostHTTPHandler "oirakif/simple-blog-service/pkg/blog-post/handler"
	blogPostRepository "oirakif/simple-blog-service/pkg/blog-post/repository"
	userRepository "oirakif/simple-blog-service/pkg/user/repository"
	"oirakif/simple-blog-service/pkg/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	authV1BasicAuthUsername := os.Getenv("AUTH_V1_BASIC_AUTH_USERNAME")
	authV1BasicAuthPassword := os.Getenv("AUTH_V1_BASIC_AUTH_PASSWORD")
	if authV1BasicAuthUsername == "" || authV1BasicAuthPassword == "" {
		panic("basic auth variables are not set")
	}
	addr := fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	cfg := mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      addr,
		DBName:    os.Getenv("DB_NAME"),
		ParseTime: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	jwtUtils := utils.NewJWTUtils([]byte(jwtSecret))
	// initiate repositories
	userRepository := userRepository.NewUserRepository(db)
	blogPostRepository := blogPostRepository.NewBlogPostRepository(db)

	// initiate domains
	authDomain := authDomain.NewAuthDomain(*userRepository, *jwtUtils)
	blogPostDomain := blogPostDomain.NewBlogPostDomain(*blogPostRepository)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	mainRouter := r.Group("/simple-blog-service")

	// initiate http handlers
	authHTTPHandler := authHTTPHandler.NewAuthHTTPHandler(
		mainRouter,
		*authDomain,
		authV1BasicAuthUsername,
		authV1BasicAuthPassword,
	)
	blogPostHttpHandler := blogPostHTTPHandler.NewBlogPostHTTPHandler(mainRouter, *blogPostDomain, *jwtUtils)

	// ship up the routes
	authHTTPHandler.InitiateRoutes()
	blogPostHttpHandler.InitiateRoutes()
	r.Run()
}
