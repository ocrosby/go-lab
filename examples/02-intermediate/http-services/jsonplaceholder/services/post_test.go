package services_test

import (
	"github.com/ocrosby/golab/jsonplaceholder/models"
	"github.com/ocrosby/golab/jsonplaceholder/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Post", func() {
	var service *services.PostService

	BeforeEach(func() {
		service = services.NewPostService(nil) // passing nil will use the default http client
	})

	AfterEach(func() {
		service = nil
	})

	Describe("GetAll", func() {
		It("returns all posts", func() {
			// Act
			posts, err := service.GetAll()

			// Assert
			Expect(err).ToNot(HaveOccurred())
			Expect(len(posts)).To(Equal(100))
		})
	})

	Describe("GetByID", func() {
		It("returns a post by ID", func() {
			// Act
			post, err := service.GetByID(1)

			// Assert
			Expect(err).ToNot(HaveOccurred())
			Expect(post.Id).To(Equal(1))
		})
	})

	Describe("GetByUserId", func() {
		It("returns posts by user ID", func() {
			// Act
			posts, err := service.GetByUserId(1)

			// Assert
			Expect(err).ToNot(HaveOccurred())
			Expect(len(posts)).To(Equal(10))
		})
	})

	Describe("Create", func() {
		It("creates a post", func() {
			// Arrange
			post := &models.Post{
				UserId: 1,
				Title:  "foo",
				Body:   "bar",
			}

			// Act
			createdPost, err := service.Create(post)

			// Assert
			Expect(err).ToNot(HaveOccurred())
			Expect(createdPost.Id).To(Equal(101))
		})
	})
})
