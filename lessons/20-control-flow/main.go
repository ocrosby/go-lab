// Command control-flow demonstrates every form of if, switch, and for
// that Go supports.
package main

import "fmt"

func main() {
	// --- if / else -------------------------------------------------------
	fmt.Println("--- if / else ---")
	for _, x := range []int{5, 0} {
		if x > 0 {
			fmt.Println(x, "is positive")
		} else if x < 0 {
			fmt.Println(x, "is negative")
		} else {
			fmt.Println(x, "is zero")
		}
	}

	// --- if with init statement -----------------------------------------
	// The `if v, ok := ...; ok {` shape scopes v and ok to the block.
	// Idiomatic for map lookups and (value, error) returns.
	fmt.Println("\n--- if with init ---")
	m := map[string]int{"answer": 42}
	if v, ok := m["answer"]; ok {
		fmt.Println("found", v, "in map")
	}

	// --- switch (with expression) ---------------------------------------
	fmt.Println("\n--- switch ---")
	for _, day := range []string{"Monday", "Saturday"} {
		switch day {
		case "Saturday", "Sunday":
			fmt.Println(day, "is a weekend")
		case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
			fmt.Println(day, "is a weekday")
		default:
			fmt.Println(day, "is unknown")
		}
	}

	// --- expression-less switch -----------------------------------------
	// Cleaner than a long if / else if chain when every branch is a
	// full boolean.
	fmt.Println("\n--- expression-less switch ---")
	age := 25
	switch {
	case age < 13:
		fmt.Printf("age %d → child\n", age)
	case age < 20:
		fmt.Printf("age %d → teen\n", age)
	case age < 65:
		fmt.Printf("age %d → adult\n", age)
	default:
		fmt.Printf("age %d → senior\n", age)
	}

	// --- for: classic three-part ----------------------------------------
	fmt.Println("\n--- for: classic C-style ---")
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// --- for: condition-only (Go's "while") -----------------------------
	fmt.Println("\n--- for: condition-only (Go's 'while') ---")
	count := 3
	fmt.Print("countdown: ")
	for count > 0 {
		fmt.Printf("%d ", count)
		count--
	}
	fmt.Println("done")

	// --- for: infinite --------------------------------------------------
	// Every "for {}" needs a break, return, or panic to exit.
	fmt.Println("\n--- for: infinite ---")
	tries := 0
	for {
		tries++
		if tries >= 3 {
			fmt.Println("break after", tries)
			break
		}
	}

	// --- for range: slice -----------------------------------------------
	// Ranging over a slice gives (index, value).
	fmt.Println("\n--- for range: slice ---")
	fruits := []string{"apple", "banana", "cherry"}
	for i, v := range fruits {
		fmt.Printf("%d:%s ", i, v)
	}
	fmt.Println()

	// --- for range: map -------------------------------------------------
	// Map iteration order is randomized per iteration — this is
	// intentional. Sort the keys first if you need determinism.
	fmt.Println("\n--- for range: map ---")
	fmt.Println("(order may vary)")
}
