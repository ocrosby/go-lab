// Package generics demonstrates Go's type parameters and constraints via
// Map/Filter/Reduce, a generic Max with a union constraint, and a
// generic Stack[T] container using the `any` constraint.
package generics

import "cmp"

// Map applies fn to each element of xs and returns a new slice. The
// signature `[T, U any]` says "any two types, possibly different." The
// standard library's `slices.Collect` and range-over-func can replace
// hand-rolled Map in modern code — this version is here for teaching.
func Map[T, U any](xs []T, fn func(T) U) []U {
	out := make([]U, len(xs))
	for i, x := range xs {
		out[i] = fn(x)
	}
	return out
}

// Filter returns a new slice with only the elements of xs for which
// keep returns true.
func Filter[T any](xs []T, keep func(T) bool) []T {
	out := make([]T, 0, len(xs))
	for _, x := range xs {
		if keep(x) {
			out = append(out, x)
		}
	}
	return out
}

// Reduce folds xs left-to-right, starting from init. Two type parameters
// because the accumulator (U) can differ from the element type (T) —
// e.g. summing []string lengths into an int.
func Reduce[T, U any](xs []T, init U, fn func(U, T) U) U {
	acc := init
	for _, x := range xs {
		acc = fn(acc, x)
	}
	return acc
}

// Max returns the larger of two values. Uses cmp.Ordered from the
// standard library (Go 1.21+) — a constraint that covers every type
// with < and > operators, without hand-listing them.
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Contains reports whether target appears in xs. Uses `comparable` —
// the built-in constraint for types that support == and !=. Excludes
// slices, maps, and function values, which cannot be compared.
func Contains[T comparable](xs []T, target T) bool {
	for _, x := range xs {
		if x == target {
			return true
		}
	}
	return false
}

// Stack[T] is a generic LIFO stack. The `any` constraint accepts every
// type — this container makes no assumption beyond "you can put values
// in and get them out." Add Push/Pop methods on any type by adjusting
// the constraint.
type Stack[T any] struct {
	items []T
}

// Push adds x to the top of the stack.
func (s *Stack[T]) Push(x T) {
	s.items = append(s.items, x)
}

// Pop removes and returns the top item. The second return is false if
// the stack was empty (comma-ok idiom, same shape as map lookups).
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return top, true
}

// Len returns the number of items on the stack.
func (s *Stack[T]) Len() int { return len(s.items) }
