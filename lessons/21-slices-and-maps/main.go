// Command slices-and-maps demonstrates the two collection types every
// Go program uses: slices for ordered lists, maps for key-value lookups.
package main

import "fmt"

func main() {
	// --- Slices ---------------------------------------------------------
	fmt.Println("--- slices ---")

	// Create with a literal.
	fruits := []string{"apple", "banana", "cherry"}
	fmt.Printf("initial: %v (len=%d, cap=%d)\n", fruits, len(fruits), cap(fruits))

	// Grow. ALWAYS reassign the result of append — the underlying array
	// may be reallocated to fit the new elements.
	fruits = append(fruits, "date")
	fruits = append(fruits, "elderberry", "fig")
	fmt.Printf("after append: %v (len=%d, cap=%d)\n", fruits, len(fruits), cap(fruits))

	// Index and slice.
	fmt.Println("first:", fruits[0])
	fmt.Println("middle:", fruits[2:4])

	// The "shared backing array" surprise: slicing does not copy.
	middle := fruits[2:4]
	middle[0] = "CHERRY-MODIFIED"
	fmt.Println("original fruits after mutating middle:", fruits)

	// --- Iterating a slice ----------------------------------------------
	fmt.Println("\n--- for range over slice ---")
	for i, v := range []string{"one", "two", "three"} {
		fmt.Printf("[%d] = %s\n", i, v)
	}

	// --- Maps -----------------------------------------------------------
	fmt.Println("\n--- maps ---")
	scores := map[string]int{
		"alice": 95,
		"bob":   82,
	}

	// Insert / update.
	scores["carol"] = 77
	scores["alice"] = 96 // overwrite

	// Delete.
	delete(scores, "bob")

	// Read: missing key returns the zero value (0 for int).
	fmt.Println("alice:", scores["alice"])
	fmt.Println("bob (missing):", scores["bob"])

	// The comma-ok idiom — the ONLY safe way to distinguish "key missing"
	// from "value happens to be zero."
	if v, ok := scores["missing"]; ok {
		fmt.Println("found:", v)
	} else {
		fmt.Println("missing: not in map")
	}

	// --- Iterating a map ------------------------------------------------
	// Order is randomized per iteration. Sort keys explicitly if you need
	// determinism.
	fmt.Println("\n--- for range over map ---")
	for k, v := range scores {
		fmt.Printf("%s=%d ", k, v)
	}
	fmt.Println()

	// --- Sets via maps --------------------------------------------------
	// Go has no built-in set type. map[T]struct{} is the idiomatic
	// substitute — struct{} occupies zero bytes.
	fmt.Println("\n--- set via map[T]struct{} ---")
	seen := map[string]struct{}{}
	for _, w := range []string{"alpha", "beta", "alpha", "gamma", "beta"} {
		seen[w] = struct{}{}
	}
	fmt.Println("unique words:", len(seen))
}
