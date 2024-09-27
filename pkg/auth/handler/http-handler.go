package handler

import (
	"net/http"
	"oirakif/simple-blog-service/pkg/auth/domain"
	"oirakif/simple-blog-service/pkg/auth/model"

	"github.com/gin-gonic/gin"
)

type UserHTTPHandler struct {
	router            *gin.Engine
	authDomain        domain.AuthDomain
	basicAuthUsername string
	basicAuthPassword string
}

func NewUserHTTPHandler(
	r *gin.Engine,
	authDomain domain.AuthDomain,
	basicAuthUsername string,
	basicAuthPassword string,
) (handler UserHTTPHandler) {
	handler = UserHTTPHandler{
		router:            r,
		authDomain:        authDomain,
		basicAuthUsername: basicAuthUsername,
		basicAuthPassword: basicAuthPassword,
	}

	return handler
}

func (h *UserHTTPHandler) InitiateRoutes() {
	usersV1 := h.router.Group("auth/v1",
		gin.BasicAuth(gin.Accounts{
			h.basicAuthUsername: h.basicAuthPassword,
		}),
	)

	usersV1.POST("/register", h.handleRegister)
	usersV1.POST("/login", h.handleLogin)

}

func (h *UserHTTPHandler) handleRegister(c *gin.Context) {
	var registerPayload model.RegisterHTTPPayload
	if err := c.ShouldBindJSON(&registerPayload); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request body is not specified"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	statusCode, response := h.authDomain.RegisterUser(registerPayload.Email, registerPayload.Password, registerPayload.Name)
	c.JSON(statusCode, response)
}

func (h *UserHTTPHandler) handleLogin(c *gin.Context) {
	var loginPayload model.LoginHTTPPayload
	if err := c.ShouldBindJSON(&loginPayload); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request body is not specified"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	statusCode, response := h.authDomain.Login(loginPayload.Email, loginPayload.Password)
	c.JSON(statusCode, response)
}
