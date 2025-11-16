package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ocrosby/golab/jsonplaceholder/models"
	pkg2 "github.com/ocrosby/golab/jsonplaceholder/pkg/http"
	"io"
	"net/http"
)

const (
	PostsEndpoint         = "https://jsonplaceholder.typicode.com/posts"
	PostEndpoint          = "https://jsonplaceholder.typicode.com/posts/%d"
	PostsByUserIdEndpoint = "https://jsonplaceholder.typicode.com/posts?userId=%d"
)

type IPostService interface {
	GetAll() ([]*models.Post, error)
	GetByID(id int) (*models.Post, error)
	GetByUserId(userId int) ([]*models.Post, error)
	Create(post *models.Post) (*models.Post, error)
	UpdateById(id int, post *models.Post) (*models.Post, error)
	DeleteById(id int) error
	Delete(post *models.Post) error
	DeleteAll() error
	PatchById(id int, post *models.Post) (*models.Post, error)
}

type PostService struct {
	client  pkg2.IHttpClient
	context context.Context
}

// NewPostService creates a new PostService
//
// If the specified client is nil, the default http client will be used
func NewPostService(client pkg2.IHttpClient) *PostService {
	svc := new(PostService)

	if client == nil {
		client = http.DefaultClient
	}

	svc.context = context.Background()
	svc.client = client

	return svc
}

// GetAll returns all posts
//
// Example:
//
// curl -X GET https://jsonplaceholder.typicode.com/posts
func (service *PostService) GetAll() ([]*models.Post, error) {
	var posts []*models.Post

	method := http.MethodGet
	url := PostsEndpoint

	fmt.Println(method, url)

	req, err := http.NewRequestWithContext(service.context, method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			fmt.Println("closing response body encountered an error: " + err.Error())
		}
	}(res.Body)

	if err = json.NewDecoder(res.Body).Decode(&posts); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetByID returns a post by ID
//
// Example:
//
// curl -X GET https://jsonplaceholder.typicode.com/posts/1
func (service *PostService) GetByID(id int) (*models.Post, error) {
	var post *models.Post

	method := http.MethodGet
	url := fmt.Sprintf(PostEndpoint, id)

	req, err := http.NewRequestWithContext(service.context, method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// GetByUserId returns posts by user ID
//
// Example:
//
// curl -X GET https://jsonplaceholder.typicode.com/posts?userId=1
func (service *PostService) GetByUserId(userId int) ([]*models.Post, error) {
	var posts []*models.Post

	method := http.MethodGet
	url := fmt.Sprintf(PostsByUserIdEndpoint, userId)

	req, err := http.NewRequestWithContext(service.context, method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// Create creates a post
//
// Example:
//
// curl -X POST -H "Content-Type: application/json" -d '{"userId": 1, "title": "foo", "body": "bar"}' https://jsonplaceholder.typicode.com/posts
func (service *PostService) Create(post *models.Post) (*models.Post, error) {
	var createdPost *models.Post

	if post == nil {
		return nil, fmt.Errorf("post cannot be nil")
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(post)
	if err != nil {
		return nil, err
	}

	method := http.MethodPost
	url := PostsEndpoint

	req, err := http.NewRequestWithContext(service.context, method, url, &buf)
	if err != nil {
		return nil, err
	}

	res, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&createdPost)
	if err != nil {
		return nil, err
	}

	return createdPost, nil
}

// UpdateById updates a post by ID
//
// Example:
//
// curl -X PUT -H "Content-Type: application/json" -d '{"userId": 1, "title": "foo", "body": "bar"}' https://jsonplaceholder.typicode.com/posts/1
func (service *PostService) UpdateById(id int, post *models.Post) (*models.Post, error) {
	var updatedPost *models.Post

	if post == nil {
		return nil, fmt.Errorf("post cannot be nil")
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(post)
	if err != nil {
		return nil, err
	}

	postByIdEndpoint := fmt.Sprintf(PostEndpoint, id)
	req, err := http.NewRequestWithContext(service.context, http.MethodPut, postByIdEndpoint, &buf)
	if err != nil {
		return nil, err
	}

	res, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&updatedPost)
	if err != nil {
		return nil, err
	}

	return updatedPost, nil
}

// DeleteById deletes a post by ID
//
// Example:
//
// curl -X DELETE https://jsonplaceholder.typicode.com/posts/1
func (service *PostService) DeleteById(id int) error {
	postByIdEndpoint := fmt.Sprintf(PostEndpoint, id)
	pRequest, err := http.NewRequestWithContext(service.context, http.MethodDelete, postByIdEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = service.client.Do(pRequest)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a post
//
// Example:
//
// curl -X DELETE https://jsonplaceholder.typicode.com/posts/1
func (service *PostService) Delete(post *models.Post) error {
	if post == nil {
		return fmt.Errorf("post cannot be nil")
	}

	return service.DeleteById(post.Id)
}

// DeleteAll deletes all posts
func (service *PostService) DeleteAll() error {
	posts, err := service.GetAll()
	if err != nil {
		return err
	}

	for _, post := range posts {
		err = service.Delete(post)
		if err != nil {
			return err
		}
	}

	return nil
}

// PatchById patches a post by ID
//
// Example:
//
// curl -X PATCH -H "Content-Type: application/json" -d '{"title": "foo"}' https://jsonplaceholder.typicode.com/posts/1
func (service *PostService) PatchById(id int, post *models.Post) (*models.Post, error) {
	var patchedPost *models.Post

	if post == nil {
		return nil, fmt.Errorf("post cannot be nil")
	}

	// TODO: Use reflection to get the fields that have been updated

	postByIdEndpoint := fmt.Sprintf(PostEndpoint, id)
	req, err := http.NewRequestWithContext(service.context, http.MethodPatch, postByIdEndpoint, nil)
	if err != nil {
		return nil, err
	}

	res, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&patchedPost)
	if err != nil {
		return nil, err
	}

	return patchedPost, nil
}
