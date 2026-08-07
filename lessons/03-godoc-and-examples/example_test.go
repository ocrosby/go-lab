package shapes_test

import (
	"errors"
	"fmt"

	shapes "github.com/ocrosby/go-lab/lessons/03-godoc-and-examples"
)

// ExampleNewRectangle shows how to construct a Rectangle and handle the
// error case. The "// Output:" line below is checked by `go test` —
// change the printed value and the test fails.
func ExampleNewRectangle() {
	r, err := shapes.NewRectangle(3, 4)
	if err != nil {
		fmt.Println("unexpected error:", err)
		return
	}
	fmt.Printf("%.0f x %.0f rectangle\n", 3.0, 4.0)
	fmt.Printf("area=%.0f perimeter=%.0f\n", r.Area(), r.Perimeter())
	// Output:
	// 3 x 4 rectangle
	// area=12 perimeter=14
}

// ExampleNewRectangle_invalid shows the error path. The "_invalid" suffix
// after the function name lets one symbol have multiple examples in godoc.
func ExampleNewRectangle_invalid() {
	_, err := shapes.NewRectangle(-1, 4)
	fmt.Println(errors.Is(err, shapes.ErrNonPositive))
	// Output: true
}

// ExampleRectangle_Area attaches an example directly to the Area method.
// godoc renders it under Rectangle.Area on pkg.go.dev.
func ExampleRectangle_Area() {
	r, _ := shapes.NewRectangle(5, 6)
	fmt.Println(r.Area())
	// Output: 30
}
