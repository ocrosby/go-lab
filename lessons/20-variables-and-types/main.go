// Command variables-and-types demonstrates Go's variable declaration
// forms, zero values, constants, iota, and explicit type conversion.
package main

import "fmt"

// Package-level constants — must use `const`, not `:=`.
const Pi = 3.14159

// A const block with iota for enum-like values. Compiles to Sunday=0,
// Monday=1, Tuesday=2, ...
const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func main() {
	// --- Zero values -----------------------------------------------------
	// A variable declared with `var` and no initializer gets its type's
	// zero value. There is no "uninitialized" garbage.
	var (
		zeroInt    int
		zeroFloat  float64
		zeroBool   bool
		zeroString string
	)
	fmt.Println("--- primitive zero values ---")
	fmt.Printf("int:     %d\n", zeroInt)
	fmt.Printf("float64: %g\n", zeroFloat)
	fmt.Printf("bool:    %v\n", zeroBool)
	fmt.Printf("string:  %q\n", zeroString)

	// --- Short declaration form ------------------------------------------
	// The := form declares AND assigns. Type is inferred from the value.
	// Only allowed inside functions.
	name := "Ada"
	age := 36
	weight := 62.3
	licensed := true
	fmt.Println("\n--- declared with := ---")
	fmt.Printf("name: %s, age: %d, weight: %g, licensed: %v\n",
		name, age, weight, licensed)

	// --- Constants -------------------------------------------------------
	fmt.Println("\n--- constants ---")
	fmt.Printf("Pi ≈ %g\n", Pi)

	// --- iota ------------------------------------------------------------
	fmt.Println("--- iota ---")
	fmt.Printf("Sunday=%d Monday=%d Tuesday=%d\n", Sunday, Monday, Tuesday)

	// --- Explicit type conversion ----------------------------------------
	// Go NEVER auto-converts numeric types. `int(42)` is required even
	// to widen a smaller int; `float64(i)` for int→float. No implicit
	// truncation, no implicit widening.
	i := 42
	f := float64(i)
	s := fmt.Sprintf("%d", i) // number → string via fmt (or strconv.Itoa)
	fmt.Println("\n--- type conversion ---")
	fmt.Printf("i=%d, f=%g, s=%s\n", i, f, s)
}
