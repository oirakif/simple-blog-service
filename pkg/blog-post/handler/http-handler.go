package handler

import (
	"net/http"
	"oirakif/simple-blog-service/pkg/blog-post/domain"
	"oirakif/simple-blog-service/pkg/blog-post/model"
	"oirakif/simple-blog-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

type BlogPostHTTPHandler struct {
	router         *gin.RouterGroup
	blogPostDomain domain.BlogPostDomain
	jwtUtils       utils.JWTUtils
}

func NewBlogPostHTTPHandler(
	r *gin.RouterGroup,
	blogPostDomain domain.BlogPostDomain,
	jwtUtils utils.JWTUtils,
) *BlogPostHTTPHandler {

	return &BlogPostHTTPHandler{
		router:         r,
		blogPostDomain: blogPostDomain,
		jwtUtils:       jwtUtils,
	}
}

func (h *BlogPostHTTPHandler) InitiateRoutes() {
	usersV1 := h.router.Group("/posts/v1")

	usersV1.POST("/posts", h.jwtUtils.ValidateToken, h.handleCreateBlogPost)
	usersV1.GET("/posts", h.jwtUtils.ValidateToken, h.handleGetAllBlogPost)
	usersV1.GET("/posts/:id", h.jwtUtils.ValidateToken, h.handleGetBlogPostByID)
	usersV1.PUT("/posts/:id", h.jwtUtils.ValidateToken, h.handleUpdateBlogPost)

}

func (h *BlogPostHTTPHandler) handleCreateBlogPost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "claims not found"})
		return
	}

	authorID, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "cannot get user ID from token"})
		return
	}

	if authorID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "cannot get user ID from token"})
		return
	}

	var payload model.CreateBlogPostHTTPPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request body is not specified"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusCode, response := h.blogPostDomain.
		CreateBlogPost(
			payload.Title,
			payload.Content,
			authorID,
		)

	c.JSON(statusCode, response)
}

func (h *BlogPostHTTPHandler) handleGetAllBlogPost(c *gin.Context) {
	var queryParams model.GetBlogPostsQueryParams

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusCode, response := h.blogPostDomain.GetBlogPosts(&queryParams)
	c.JSON(statusCode, response)
}

func (h *BlogPostHTTPHandler) handleGetBlogPostByID(c *gin.Context) {
	var pathParams model.BlogPostsPathParam
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var queryParams model.GetBlogPostsQueryParams
	queryParams.ID = pathParams.ID
	perPage := 1
	queryParams.PerPage = &perPage
	statusCode, response := h.blogPostDomain.GetBlogPosts(&queryParams)
	data := response.Data.([]model.BlogPost)
	response.Data = data[0]
	c.JSON(statusCode, response)
}

func (h *BlogPostHTTPHandler) handleUpdateBlogPost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "claims not found"})
		return
	}

	authorID, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "cannot get user ID from token"})
		return
	}

	if authorID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "cannot get user ID from token"})
		return
	}

	var pathParams model.BlogPostsPathParam
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payload model.UpdateBlogPostHTTPPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request body is not specified"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusCode, response := h.blogPostDomain.
		UpdateBlogPost(
			*pathParams.ID,
			authorID,
			payload.Title,
			payload.Content,
			payload.Status,
		)
	c.JSON(statusCode, response)
}
