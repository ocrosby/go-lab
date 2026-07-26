// Command pointers demonstrates & and *, nil pointers, and the value-vs-
// pointer-receiver distinction.
package main

import "fmt"

// Counter has one field. Two methods below show the difference between
// value and pointer receivers.
type Counter struct {
	n int
}

// IncrementValue uses a VALUE receiver — c is a copy. Mutations to c.n
// affect only the copy and are lost when the method returns. Almost
// certainly a bug for a "counter."
func (c Counter) IncrementValue() {
	c.n++
}

// Increment uses a POINTER receiver — c points at the original struct.
// Mutations to c.n are visible to the caller. What you almost always
// want for a state-carrying method.
func (c *Counter) Increment() {
	c.n++
}

// Value returns the current count. Read-only, so a value receiver would
// technically work — but Go convention says to use pointer receivers on
// every method of a type if ANY of them mutate. Consistency prevents
// interface-satisfaction surprises.
func (c *Counter) Value() int {
	return c.n
}

// addOneValue takes an int BY VALUE. Mutations don't propagate.
func addOneValue(n int) { n++ }

// addOne takes a *int. *n dereferences to modify the caller's variable.
func addOne(n *int) { *n++ }

func main() {
	// --- & and * on primitives -------------------------------------------
	fmt.Println("--- & and * ---")
	x := 42
	p := &x
	fmt.Printf("x=%d, p=%v, *p=%d\n", x, p, *p)

	*p = 100
	fmt.Printf("after *p = 100: x=%d\n", x)

	// --- nil pointers ----------------------------------------------------
	fmt.Println("\n--- nil pointer ---")
	var np *int
	fmt.Printf("uninitialized: np = %v (nil? %v)\n", np, np == nil)
	// Uncomment the next line and re-run — you'll see a panic:
	// fmt.Println(*np)

	// --- pass by value vs pass by pointer -------------------------------
	fmt.Println("\n--- pass by value vs by pointer ---")
	n := 5
	addOneValue(n)
	fmt.Printf("after addOneValue: n = %d (unchanged)\n", n)
	addOne(&n)
	fmt.Printf("after addOne:      n = %d\n", n)

	// --- value receiver vs pointer receiver ------------------------------
	fmt.Println("\n--- value vs pointer receiver ---")
	c := Counter{}
	c.IncrementValue()
	c.IncrementValue()
	c.IncrementValue()
	fmt.Printf("after 3× IncrementValue: c.n = %d (mutations lost)\n", c.n)

	c.Increment()
	c.Increment()
	c.Increment()
	fmt.Printf("after 3× Increment:      c.n = %d\n", c.n)

	// --- new(T) vs &T{} -------------------------------------------------
	fmt.Println("\n--- new(T) vs &T{} ---")
	a := new(Counter)
	b := &Counter{n: 10}
	fmt.Printf("new(Counter):     %+v\n", *a)
	fmt.Printf("&Counter{n: 10}:  %+v\n", *b)
}
