package shapes

import "errors"

// ErrNonPositive is returned when a shape is constructed with a zero
// or negative dimension. A rectangle with side <= 0 has no meaningful
// area or perimeter, so the constructor rejects it up front.
var ErrNonPositive = errors.New("shapes: dimensions must be positive")

// Rectangle is an axis-aligned rectangle described by its width and height.
// The zero value is not a valid Rectangle — use NewRectangle to construct one.
type Rectangle struct {
	width, height float64
}

// NewRectangle returns a Rectangle with the given width and height.
// It returns ErrNonPositive if either dimension is <= 0.
func NewRectangle(width, height float64) (Rectangle, error) {
	if width <= 0 || height <= 0 {
		return Rectangle{}, ErrNonPositive
	}
	return Rectangle{width: width, height: height}, nil
}

// Area returns the rectangle's area (width * height).
func (r Rectangle) Area() float64 {
	return r.width * r.height
}

// Perimeter returns the rectangle's perimeter (2 * (width + height)).
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}
