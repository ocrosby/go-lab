// Package funcs demonstrates anonymous functions, closures, higher-order
// helpers, variadic parameters, named return values, and init().
package funcs

import "errors"

// initMessage is a package-level variable set by init() below. It's
// visible here so a test can assert init() ran.
var initMessage string

// init runs once at program startup, before main. Each file in a
// package may have any number of init() functions; each runs once.
// Here we use it to prime a package-level variable — a common use
// (registry setup, cache priming, env-var validation).
func init() {
	initMessage = "package initialized"
}

// InitMessage exposes the string init() set. The tests use it to prove
// init() ran before test code did.
func InitMessage() string { return initMessage }

// MakeAdder returns a closure that adds n to whatever it's called with.
// The returned function captures `n` by reference — each call to
// MakeAdder produces its own independent n.
func MakeAdder(n int) func(int) int {
	return func(x int) int {
		return n + x
	}
}

// MakeCounter returns a closure that increments a private counter each
// call and returns the new value. This is the canonical "stateful
// generator" pattern in Go.
func MakeCounter() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

// Sum uses a variadic parameter — accepts zero or more int arguments.
// Callers can pass individual values (Sum(1, 2, 3)) or spread a slice
// (Sum(nums...)).
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Divide returns quotient and error using NAMED return values. The
// `return` keyword with no arguments (the "naked return") returns the
// currently-declared values of quotient and err — assigned inside the
// body, or their zero values if unassigned.
//
// Use named returns sparingly — they help in short helpers but hurt
// readability in long functions.
func Divide(a, b int) (quotient int, err error) {
	if b == 0 {
		err = errors.New("divide by zero")
		return // returns quotient=0, err=the error we set
	}
	quotient = a / b
	return
}

// Apply demonstrates functions as first-class values — fn is a parameter
// of function type, called inside Apply's body. Callers pass whatever
// unary transformation they want.
func Apply(xs []int, fn func(int) int) []int {
	out := make([]int, len(xs))
	for i, x := range xs {
		out[i] = fn(x)
	}
	return out
}

// Dispatch shows functions in a map as a dispatch table — cleaner than
// a long switch when the branches are small and homogeneous.
func Dispatch(op string, a, b int) (int, bool) {
	ops := map[string]func(int, int) int{
		"+": func(a, b int) int { return a + b },
		"-": func(a, b int) int { return a - b },
		"*": func(a, b int) int { return a * b },
	}
	fn, ok := ops[op]
	if !ok {
		return 0, false
	}
	return fn(a, b), true
}
