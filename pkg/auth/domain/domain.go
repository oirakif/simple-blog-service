package domain

import (
	"net/http"
	authModel "oirakif/simple-blog-service/pkg/auth/model"
	userModel "oirakif/simple-blog-service/pkg/user/model"
	"oirakif/simple-blog-service/pkg/user/repository"

	"oirakif/simple-blog-service/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthDomain struct {
	userRepository *repository.UserRepository
	jwtUtils       *utils.JWTUtils
}

func NewAuthDomain(
	userRepository *repository.UserRepository,
	jwtUtils *utils.JWTUtils,
) *AuthDomain {
	return &AuthDomain{
		userRepository: userRepository,
		jwtUtils:       jwtUtils,
	}

}

func (d *AuthDomain) RegisterUser(email, password, name string) (statusCode int, response authModel.RegisterResponse) {
	hashedPassword := utils.HashSHA256(password)
	filterQuery := userModel.FindUserFilterQuery{
		Email: &email,
	}
	retrievedUser, err := d.userRepository.FindUser(filterQuery)
	if err != nil {
		response.Error = true
		response.Message = "error occured while querying user data"

		return http.StatusInternalServerError, response
	}
	if retrievedUser != nil {
		response.Error = true
		response.Message = "email is already registered"

		return http.StatusConflict, response
	}
	currentTimestamp := time.Now()
	newUser := userModel.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    currentTimestamp,
		UpdatedAt:    currentTimestamp,
	}
	newUserID, err := d.userRepository.InsertUser(newUser)
	if err != nil {
		response.Error = true
		response.Message = "internal server error"

		return http.StatusInternalServerError, response
	}

	newUser.ID = newUserID
	response.Data = newUser
	response.Message = "new user has been registered"

	return http.StatusCreated, response
}

func (d *AuthDomain) Login(email, password string) (statusCode int, response authModel.LoginResponse) {
	hashedPassword := utils.HashSHA256(password)
	filterQuery := userModel.FindUserFilterQuery{
		Email:        &email,
		PasswordHash: &hashedPassword,
	}

	retrievedUser, err := d.userRepository.FindUser(filterQuery)
	if err != nil {
		response.Error = true
		response.Message = "error occured while querying user data"

		return http.StatusInternalServerError, response
	}
	if retrievedUser == nil {
		response.Error = true
		response.Message = "invalid email or password"

		return http.StatusInternalServerError, response
	}
	// Create the claims
	claims := utils.JWTClaims{
		UserID:          retrievedUser.ID,
		UserProfileName: retrievedUser.Name,
		Email:           email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token, err := d.jwtUtils.GenerateJWT(email, claims)
	if err != nil {
		response.Error = true
		response.Message = "error occured while generating token"

		return http.StatusInternalServerError, response

	}
	response.Token = token
	response.Message = "login successful"

	return http.StatusOK, response
}
