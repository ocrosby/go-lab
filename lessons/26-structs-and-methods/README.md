# Structs and methods

Go's way of grouping fields into a named type, and attaching functions to that type. Underlies interfaces, composition, JSON marshaling, and about 80% of any real Go codebase.

> **Recommended before lesson 05 (composition) and lesson 06 (interfaces).** Both use `type X struct { ... }` and `func (x *X) Foo()` as their vocabulary — this lesson introduces both.

## Why it matters

Almost every real Go program defines `struct` types (data grouped together) and attaches methods to them. Interfaces are just method sets, so you can't reason about interfaces without first reasoning about methods. This lesson bridges from the primitive types (lesson 19) to the composed types the rest of the syllabus uses.

## Prerequisites

- Lesson 19: types.
- Lesson 22: pointers (structs and methods lean on `*T`).

## Run it

```bash
go test -race ./lessons/26-structs-and-methods
```

Expected: 8 passes.

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`user.go`](./user.go) | Struct definition, exported vs unexported fields, struct literals, value and pointer receivers, `Stringer` interface satisfied automatically, an anonymous struct. |
| [`user_test.go`](./user_test.go) | Tests exercising each pattern above through the exported surface. |

## Defining a struct

```go
type User struct {
    ID    string    // exported field — accessible outside the package
    Name  string    // exported
    email string    // unexported — only visible inside this package
    Roles []string  // slices, maps, other structs — all allowed as fields
}
```

The **capitalization rule** you learned in lesson 02 applies field-by-field: capital-first-letter fields are exported (public), lowercase ones are package-private. Mix freely — no keyword required.

## Struct literals

Three ways to construct a value:

```go
// 1. Named-field literal (recommended) — safe against field reordering
u := User{
    ID:   "u1",
    Name: "Ada",
}

// 2. Positional literal — brittle; if the struct gains a field, this
//    breaks silently or noisily depending on where.
u := User{"u1", "Ada", "ada@example.com", nil}

// 3. Zero value — every field at its own type's zero value
var u User            // ID: "", Name: "", email: "", Roles: nil
u := User{}           // equivalent
```

Prefer named-field literals in every case except tiny throwaway types (like `time.Duration(5 * time.Second)`-shape pairs).

## Methods

A **method** is a function with a **receiver**: a variable declared between `func` and the name, listing the type this method belongs to.

```go
type User struct{ Name string }

// Value receiver — u is a COPY of the User.
func (u User) DisplayName() string {
    return "User: " + u.Name
}

// Pointer receiver — u points at the ORIGINAL User.
// Use pointer receivers when the method mutates the receiver, or when
// the struct is large enough that copying costs, or when the struct
// contains a mutex/channel/etc that must not be copied.
func (u *User) SetName(n string) {
    u.Name = n
}
```

Callers use the method with dot syntax:

```go
u := &User{Name: "Ada"}
fmt.Println(u.DisplayName())  // "User: Ada"
u.SetName("Grace")
fmt.Println(u.DisplayName())  // "User: Grace"
```

Go automatically takes the address for you when calling a pointer-receiver method on an addressable value, so `u.SetName(...)` works whether `u` is `User` or `*User`.

### Consistency rule

If **any** method on `T` has a pointer receiver, use pointer receivers for **all** methods on `T`. Mixing value and pointer receivers is a common source of interface-satisfaction bugs — a value type doesn't satisfy an interface if the required method has a pointer receiver.

Lesson 22 covered this in the pointers context; here it applies to struct methods specifically.

> **Confused about when to use value vs pointer receivers?** Lesson 22's ["Value semantics vs pointer semantics — the deep dive"](../22-pointers/README.md#value-semantics-vs-pointer-semantics--the-deep-dive) section has the full treatment: what "value" actually copies (including the shallow-copy nuance for slice/map fields), addressability rules, the interface-satisfaction asymmetry between `T` and `*T`, when to prefer the immutable-value idiom, constructor conventions, and a decision table. This is the single most common source of confusion for people new to Go — worth 15 minutes reading if it hasn't clicked yet.

## Methods can be defined on any type in the same package

Not just structs. Any type you defined here can have methods:

```go
type Celsius float64

func (c Celsius) Fahrenheit() Celsius {
    return c*9/5 + 32
}
```

This is why Go doesn't need "extension methods" — you just declare a named type and attach the methods.

## The `Stringer` interface (implicit)

`fmt` looks for a `String() string` method on any value it prints. If your struct has one, `fmt.Println(user)` uses it:

```go
func (u User) String() string {
    return fmt.Sprintf("User(%s: %s)", u.ID, u.Name)
}

fmt.Println(User{ID: "u1", Name: "Ada"})  // prints "User(u1: Ada)"
```

You didn't `implements Stringer` anywhere — the interface is satisfied structurally (lesson 06 has the full story on this).

## Anonymous structs

Sometimes you need a one-off shape and don't want to name it — for a config, a JSON payload, a test fixture:

```go
config := struct {
    Host string
    Port int
}{
    Host: "localhost",
    Port: 8080,
}
```

Common in test setup and inline JSON handling. Don't overuse them — named types read better once the shape is used in more than one place.

## Struct comparison and hashability

Two struct values are `==`-comparable when **every** field is comparable. That means:

- All-primitive structs are comparable and usable as map keys.
- A struct containing a slice, map, or function is not comparable — no `==`, cannot be a map key.

```go
type Point struct{ X, Y int }
p1 == p2                     // OK
m := map[Point]string{}       // OK — Point works as a key

type WithSlice struct{ IDs []string }
w1 == w2                     // compile error
```

Use `reflect.DeepEqual` (slow) or a hand-written comparison for structs with non-comparable fields.

## Try it yourself

1. Add an `Age int` field to `User`. Update the test fixture. Notice which tests you have to change (spoiler: none, if you used named-field literals).
2. Add a value-receiver method `IsActive() bool` alongside the existing pointer-receiver methods. Run `go vet` — it may warn about mixed receivers.
3. Try comparing two `User` values with `==`. Now add a `Roles []string` field and try again. Read the compiler error.
4. Add a `String() string` method to `User` so `fmt.Println(u)` prints something nicer than the default `{u1 Ada ...}`.
5. Define an anonymous struct inside a function for a temporary config. Print its fields.

## Common pitfalls

- **Mixed value/pointer receivers.** Half your methods take `T`, half take `*T`. Only `*T` satisfies interfaces that need both. The rule: pick one, stick to it, per type.
- **Passing large structs by value.** Every function call copies. If your struct is 500 bytes and you're passing it through five layers, that's a real cost. Use a pointer.
- **Comparing structs with non-comparable fields.** `s1 == s2` fails to compile if the struct contains a slice, map, or func. Non-obvious until you try.
- **Copying a struct with a mutex.** `sync.Mutex` should never be copied — its zero value is a fresh usable mutex, but copying it leaves both copies pointing at the same state. Always use `*Mutex` embedded in a struct, and use pointer receivers on that struct.
- **Struct tags with typos.** `json:"emial"` compiles fine. It silently fails to marshal/unmarshal the field. Lesson 30 covers struct tags in detail.

## You've understood this lesson when...

- You can write a `type Foo struct { ... }` with mixed exported and unexported fields.
- You can write value and pointer receivers and know when to use which.
- You can predict which structs will fail to compile when compared with `==`.
- You can add a `String() string` method to any type and know why `fmt` picks it up automatically.
- You can spot mixed value/pointer receivers in a code review.

## Next

- **Next lesson (recommended):** [27-functions-and-closures](../27-functions-and-closures/) — first-class functions, closures, and variadic parameters — the last language-mechanics lesson before you re-enter the main track.
