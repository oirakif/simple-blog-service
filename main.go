package main

import (
	"database/sql"
	"fmt"
	"oirakif/simple-blog-service/pkg/auth/domain"
	"oirakif/simple-blog-service/pkg/auth/handler"
	"oirakif/simple-blog-service/pkg/user/repository"
	"oirakif/simple-blog-service/pkg/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	usersV1BasicAuthUsername := os.Getenv("USER_V1_BASIC_AUTH_USERNAME")
	usersV1BasicAuthPassword := os.Getenv("USER_V1_BASIC_AUTH_PASSWORD")
	if usersV1BasicAuthUsername == "" || usersV1BasicAuthPassword == "" {
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
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	userRepository := repository.NewUserRepository(db)
	jwtUtils := utils.NewJWTUtils(os.Getenv("JWT_SECRET_KEY"))
	authDomain := domain.NewAuthDomain(userRepository, jwtUtils)
	r := gin.Default()
	userHTTPHandler := handler.NewUserHTTPHandler(
		r,
		authDomain,
		usersV1BasicAuthUsername,
		usersV1BasicAuthPassword,
	)

	userHTTPHandler.InitiateRoutes()
	r.Run()
}
