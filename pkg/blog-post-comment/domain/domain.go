package domain

import (
	"net/http"
	"oirakif/simple-blog-service/pkg/blog-post-comment/model"
	"oirakif/simple-blog-service/pkg/blog-post-comment/repository"
	blogPostModel "oirakif/simple-blog-service/pkg/blog-post/model"
	blogPostRepository "oirakif/simple-blog-service/pkg/blog-post/repository"

	"time"
)

type BlogPostCommentDomain struct {
	blogPostCommentRepository *repository.BlogPostCommentRepository
	blogPostRepository        *blogPostRepository.BlogPostRepository
}

func NewBlogPostCommentDomain(
	blogPostCommentRepository *repository.BlogPostCommentRepository,
	blogPostRepository *blogPostRepository.BlogPostRepository,
) *BlogPostCommentDomain {

	return &BlogPostCommentDomain{
		blogPostCommentRepository: blogPostCommentRepository,
		blogPostRepository:        blogPostRepository,
	}
}

func (d *BlogPostCommentDomain) CreateBlogPostComment(authorName string, postID *int, content string) (statusCode int, response model.BlogPostCommentResponse) {
	status := "ACTIVE"
	filterQuery := blogPostModel.BlogPostFilterQuery{
		ID:     postID,
		Status: &status,
		Limit:  1,
		Offset: 0,
	}

	blogPosts, err := d.blogPostRepository.GetBlogPosts(filterQuery)
	if err != nil {
		response.Error = true
		response.Message = "error occured while getting blog posts data"

		return http.StatusInternalServerError, response
	}
	if len(blogPosts) == 0 {
		response.Error = true
		response.Message = "blog post is not found"

		return http.StatusNotFound, response
	}

	currentTimestamp := time.Now()
	newBlogPostComment := model.BlogPostComment{
		PostID:     postID,
		AuthorName: &authorName,
		Content:    &content,
		Status:     &status,
		CreatedAt:  &currentTimestamp,
		UpdatedAt:  &currentTimestamp,
	}

	newBlogPostID, err := d.blogPostCommentRepository.InsertBlogPostComments(newBlogPostComment)
	if err != nil {
		response.Error = true
		response.Message = "error occured while creating new blog post comment"

		return http.StatusInternalServerError, response
	}

	newBlogPostComment.ID = &newBlogPostID

	response.Data = newBlogPostComment
	response.Message = "new blog post comment created"

	return http.StatusCreated, response
}
