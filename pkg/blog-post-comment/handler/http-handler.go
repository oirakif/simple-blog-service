package handler

import (
	"net/http"
	"oirakif/simple-blog-service/pkg/blog-post-comment/domain"
	"oirakif/simple-blog-service/pkg/blog-post-comment/model"
	"oirakif/simple-blog-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

type BlogPostCommentHTTPHandler struct {
	router                *gin.RouterGroup
	blogPostCommentDomain domain.BlogPostCommentDomain
	jwtUtils              utils.JWTUtils
}

func NewBlogPostCommentHTTPHandler(
	r *gin.RouterGroup,
	blogPostCommentDomain domain.BlogPostCommentDomain,
	jwtUtils utils.JWTUtils,
) *BlogPostCommentHTTPHandler {

	return &BlogPostCommentHTTPHandler{
		router:                r,
		blogPostCommentDomain: blogPostCommentDomain,
		jwtUtils:              jwtUtils,
	}
}

func (h *BlogPostCommentHTTPHandler) InitiateRoutes() {
	blogPostCommentsV1 := h.router.Group("/posts/v1")

	blogPostCommentsV1.POST("/posts/:postID/comments", h.jwtUtils.ValidateToken, h.handleCreateBlogPostComment)

}

func (h *BlogPostCommentHTTPHandler) handleCreateBlogPostComment(c *gin.Context) {
	var pathParams model.BlogPostCommentPathParam
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorName, exists := c.Get("user_profile_name")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "claims not found"})
		return
	}

	authorNameParsed, ok := authorName.(string)
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "cannot get user name from token"})
		return
	}

	var payload model.CreateBlogPostCommentHTTPPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request body is not specified"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusCode, response := h.blogPostCommentDomain.CreateBlogPostComment(authorNameParsed, pathParams.PostID, payload.Content)

	c.JSON(statusCode, response)
}
