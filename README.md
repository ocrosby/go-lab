# Go Laboratory

Syntax of a function in Go

```go
func function_name([parameter list]) [return_types] {
    // body of the function
}
```

Four Categories of variables in Go

1. Boolean

There are two values true and false.  Bool is one of the four elementary types.

2. Numeric

These can be integers or floating-point values.  

Examples:

* uint8 (unsigned 8-bit integers from 0 to 255)
* int16 (signed 16-bit integers)
* float32
* float64
* complex62 (float32 with real and imaginary parts)
* complex128 (float64 with real and imaginary parts)
* Both int and float are elementary types

3. Strings

Once created, the string cannot be changed.
A string is also an elementary type.

1. Derived Types

These include pointer types, array types, structure types, map types, etc.

### Constants

Go doesn't allow you to mix numeric types.

Unless you explicitly give a constant a type, it's considered to be untyped, even if you give it a name.  An untyped constant that's an integer can only be used where integers are allowed.  An untyped floating-point constant can be used wherever a floating point is permitted.

All untyped constants in Go have default types.  Constants that are integers default to "int", floats default to "float64", characters to "rune", etc.  When you declare a type, the constant becomes typed.  Again, constants must be declared as their correct type, or else the program will return an error.

Declare constants without a type unless you absolutely need them;  declaring a type makes you lose flexibility of being able to mix types in an operation.

#### Rules for constant expressions in Go

* Comparison between two untyped constants results in an untyped boolean constant ("true"/"false").
* Operands of the same type result in the same type.  The expression "11/2" results in "5" rather than "5.5" because it's truncated to an integer.
* If they're not the same type, the result is the broader of the two according to this logic: "integer<rune<floating-point<complex"

### Variables

Names of variables can use letters, digits, and underscores.  Golang variables have a specific size, memory layout, range of values for the memory layout, and possible operations.  The variable definition tells the compiler how much storage to create for the variable and where to put it.

You use the var keyword to declare a variable.

### Three types of operators in Go.

#### Bitwise Operators

* AND: &
* OR: |
* XOR and COMPLEMENT: ^
* CLEAR: &^
* Arithmetic

#### Arithmetic Operators

* +, Addition
* *, Multiplication
* /, Division
* %, Percentage

#### Logical Operators

* Equal: =
* Not Equal: !=
* Less than: <
* Less than or equal to: <=
* Greater than: >
* Greater than or equal to: >=

### Strings

Strings are immutable, or read-only.  The characters represent bytes that are UTF-8 encoded.

```go
package main

import "fmt"

func main() {
    var my_words string

    my_words = "Hello World!"

    fmt.println("String: ", my_words)
}
```

### Times and Dates

In Go, times include a location that determines the date and time associated with that location.  If it's not specified, the time defaults to UTC.

If you want a timestamp, Time.Now: the signature is:

```go
func Now() time
```

For the date and time, Time.Date, the signature is in the format yyyy-mm-dd hh-mm-ss + nanoseconds:

```go
func Date (year int, month, day, hour, min, sec, nsec int, loc *Location)
```

In Go, the duration is what's elapsed in nanoseconds betweeen two instants written in int64nanosecond count.

```go
func Since(t Time) duration
```

If you want to know how long until time t, the function is:

```go
func Until(t Time) duration
```

### Control Structures

#### Conditionals

```go
package main

import(
    "fmt"
)

func main() {
    var grade = A

    x:= true

    if x {
        fmt.PrintLn(grade)
    }
}
```

```go
if condition {
    // code that will be executed if the condition is true
}
```

```go
if condition {
    // code that will be executed if the condition is true
} else {
    // code that will be executed if the condition is false
}
```

```go
if condition_1 {
    // code that will be executed if the condition_1 is true
} else if condition_2 {
    // code that will be executed if the condition_2 is true
} else {
    // code that will be executed if both conditions are false
}
```

```go
switch target {
    case case_1:
        // handle case_1
    case case_2, case_3:
        // handle cases 2/3
    default:
        // handle default case
}
```

Note: In Go, the control switch runs the first case where the condition is true and ignores the remainder.  So you don't have to explicitly break out.

#### Loops

```go
package main

import(
    "fmt"
)

func main() {
    for x:=1; x<= 10; x++ {
        fmt.PrintLn(x)
    }
}
```

