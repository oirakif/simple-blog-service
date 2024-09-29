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
	status := "ACTIVE"
	newBlogPost := model.BlogPost{
		Title:     &title,
		Content:   &content,
		AuthorID:  &authorID,
		Status:    &status,
		CreatedAt: &currentTimestamp,
		UpdatedAt: &currentTimestamp,
	}

	newBlogPostID, err := d.blogPostRepository.InsertBlogPost(newBlogPost)
	if err != nil {
		response.Error = true
		response.Message = "error occured while creating new blog post"

		return http.StatusInternalServerError, response
	}

	newBlogPost.ID = &newBlogPostID

	response.Data = newBlogPost
	response.Message = "new blog post created"

	return http.StatusCreated, response
}

func (d *BlogPostDomain) GetBlogPosts(queries *model.GetBlogPostsQueryParams) (statusCode int, response model.BlogPostResponse) {
	page, perPage := 1, 25
	sortBy, sortOrder := "created_at", "desc"
	status := "ACTIVE"
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
	if queries.Status != nil {
		status = *queries.Status
	}
	offset := utils.CalculateOffset(page, perPage)

	filterQuery.ID = queries.ID
	filterQuery.AuthorID = queries.AuthorID
	filterQuery.Title = queries.Title
	filterQuery.Status = &status
	filterQuery.Limit = perPage
	filterQuery.Offset = offset
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

func (d *BlogPostDomain) UpdateBlogPost(authorID, blogPostID int, title, content, status *string) (statusCode int, response model.BlogPostResponse) {
	currentTimestamp := time.Now()
	filterQuery := model.BlogPostFilterQuery{
		ID:       &blogPostID,
		AuthorID: &authorID,
		Limit:    1,
	}

	data, err := d.blogPostRepository.GetBlogPosts(filterQuery)
	if err != nil {
		response.Error = true
		response.Message = "error occured while getting blog post data"

		return http.StatusInternalServerError, response
	}
	if len(data) == 0 {
		response.Error = true
		response.Message = "blog post is not found"

		return http.StatusNotFound, response
	}

	updatePayload := model.BlogPost{
		Title:     title,
		Content:   content,
		AuthorID:  &authorID,
		Status:    status,
		UpdatedAt: &currentTimestamp,
	}

	err = d.blogPostRepository.UpdateBlogPost(filterQuery, updatePayload)
	if err != nil {
		response.Error = true
		response.Message = "error occured while updating new blog post"

		return http.StatusInternalServerError, response
	}

	return http.StatusNoContent, response
}

func (d *BlogPostDomain) DeleteBlogPost(authorID, blogPostID int) (statusCode int, response model.BlogPostResponse) {
	currentTimestamp := time.Now()
	filterQuery := model.BlogPostFilterQuery{
		ID:       &blogPostID,
		AuthorID: &authorID,
		Limit:    1,
	}

	data, err := d.blogPostRepository.GetBlogPosts(filterQuery)
	if err != nil {
		response.Error = true
		response.Message = "error occured while getting blog post data"

		return http.StatusInternalServerError, response
	}
	if len(data) == 0 {
		response.Error = true
		response.Message = "blog post is not found"

		return http.StatusNotFound, response
	}

	inactiveStatus := "INACTIVE"
	updatePayload := model.BlogPost{
		Status:    &inactiveStatus,
		UpdatedAt: &currentTimestamp,
	}

	err = d.blogPostRepository.UpdateBlogPost(filterQuery, updatePayload)
	if err != nil {
		response.Error = true
		response.Message = "error occured while deleting blog post"

		return http.StatusInternalServerError, response
	}

	return http.StatusNoContent, response
}
