package main

import (
	"fmt"
	"github.com/ocrosby/golab/jsonplaceholder/services"
)

func main() {
	// Get a post by ID
	fmt.Println("Get a post by ID")
	postService := services.NewPostService()
	post, err := postService.GetByID(1)
	if err != nil {
		panic(err)
	}

	fmt.Println(*post)

	// Get all posts
	fmt.Println("Get all posts")
	posts, err := postService.GetAll()
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of posts:", len(posts))

	// Get posts by user ID
	fmt.Println("Get posts by user ID")
	posts, err = postService.GetByUserId(1)
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of posts:", len(posts))

}
