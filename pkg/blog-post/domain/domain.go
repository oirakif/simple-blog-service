package domain

import (
	"net/http"
	"oirakif/simple-blog-service/pkg/blog-post/model"
	"oirakif/simple-blog-service/pkg/blog-post/repository"
	"oirakif/simple-blog-service/pkg/utils"

	"time"
)

type BlogPostDomain struct {
	blogPostRepository repository.BlogPostRepository
}

func NewBlogPostDomain(
	blogPostRepository repository.BlogPostRepository,
) *BlogPostDomain {

	return &BlogPostDomain{
		blogPostRepository: blogPostRepository,
	}
}

func (d *BlogPostDomain) CreateBlogPost(title, content string, authorID int) (statusCode int, response model.BlogPostResponse) {
	currentTimestamp := time.Now()
	newBlogPost := model.BlogPost{
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: currentTimestamp,
		UpdatedAt: currentTimestamp,
	}
	newBlogPostID, err := d.blogPostRepository.InsertBlogPost(newBlogPost)
	if err != nil {
		response.Error = true
		response.Message = "error occured while creating new blog post"

		return http.StatusInternalServerError, response
	}
	newBlogPost.ID = newBlogPostID

	response.Data = newBlogPost
	response.Message = "new blog post created"

	return http.StatusCreated, response
}

func (d *BlogPostDomain) GetBlogPosts(queries *model.GetBlogPostsQueryParams) (statusCode int, response model.BlogPostResponse) {
	page, perPage := 1, 25
	sortBy, sortOrder := "created_at", "desc"
	var filterQuery model.BlogPostFilterQuery
	if queries.Page != nil {
		page = *queries.Page
	}
	if queries.PerPage != nil {
		perPage = *queries.PerPage
	}
	if queries.SortBy != nil {
		sortBy = *queries.SortBy
	}
	if queries.SortOrder != nil {
		sortOrder = *queries.SortOrder
	}
	offset := utils.CalculateOffset(page, perPage)

	filterQuery.Limit = perPage
	filterQuery.Offset = offset
	filterQuery.AuthorID = queries.AuthorID
	filterQuery.Title = queries.Title
	filterQuery.Status = queries.Status
	filterQuery.SortBy = &sortBy
	filterQuery.SortOrder = &sortOrder

	blogPosts, err := d.blogPostRepository.GetBlogPosts(filterQuery)
	if err != nil {
		response.Error = true
		response.Message = "error occured while getting blog posts data"

		return http.StatusInternalServerError, response
	}

	response.Data = blogPosts
	response.Message = "blog posts data"

	return http.StatusOK, response
}
