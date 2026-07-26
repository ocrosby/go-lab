# Composition

How Go builds bigger types out of smaller ones — without classes and without inheritance.

## Why it matters

If you've used Java, Python, or C++, you know inheritance: `SeniorNinja extends Ninja extends Human`. **Go does not have inheritance.** Instead, Go gives you *composition* — you put one type inside another and use it. The syntax is simple, the mental model is different, and once you get it you'll design smaller, more flexible types than you would in an inheritance-heavy language.

## Prerequisites

- Lesson 02: functions and packages.
- Comfort with reading a Go `struct` definition (a named record with fields).

## Run it

This lesson doesn't have tests yet — it's a reading lesson. To confirm it compiles:

```bash
go build ./lessons/05-composition
```

Silent success means it built.

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`human.go`](./human.go) | The base type: a `Human` with a `Chill()` method, plus a `HumanInterface`. |
| [`ninja.go`](./ninja.go) | A `Ninja` that *has* a `Human` as a field, and adds an `Attack()` method. |
| [`senior_ninja.go`](./senior_ninja.go) | A `SeniorNinja` that *has* a `Ninja` as a field, and overrides `Attack()`. |

## Mental model

**Inheritance** (Java/Python): "A ninja *is* a human. A senior ninja *is* a ninja. All the human methods automatically apply."

**Composition** (Go): "A ninja *has* a human inside it. The ninja gets human behaviour by explicitly forwarding calls to the human it contains."

Look at [`ninja.go`](./ninja.go):

```go
type Ninja struct {
    human Human
}

func (n Ninja) Chill() {
    n.human.Chill()   // forward the call
}

func (n Ninja) Attack() {
    fmt.Println("Throwing ninja stars")
}
```

The `Ninja` doesn't inherit `Chill()` from `Human`. It has a `Human` field and defines its own `Chill()` method that forwards to `n.human.Chill()`. Every "inherited" behaviour is an explicit forward.

## The shortcut: embedded fields

The code in this lesson writes out the forwarding explicitly. Go has a shorthand called **embedding**: if you drop the field name and only write the type, Go automatically forwards method calls for you.

```go
type Ninja struct {
    Human   // embedded — no field name, just the type
}
// Now n.Chill() works without a Chill() method on Ninja.
```

Both styles are valid Go. Embedding is more concise; explicit fields are more, well, explicit. The rest of this repo prefers small, purpose-built types over deep embedding trees — see [lesson 06](../06-interfaces-and-mocking/) for the alternative.

## Interfaces along the way

Each file also defines an interface:

- `HumanInterface` requires `Chill()`.
- `NinjaInterface` requires `HumanInterface` methods *plus* `Attack()`.
- `SeniorNinjaInterface` is the same as `NinjaInterface`.

Notice how `NinjaInterface` embeds `HumanInterface` — that's interface composition. A type satisfies `NinjaInterface` if it has both `Chill()` and `Attack()`.

Also notice you never *declare* that `Ninja` implements `NinjaInterface`. Go checks structurally at compile time — if the methods match, the interface is satisfied. This is called **structural typing** or **duck typing at compile time**. Lesson 06 leans on this heavily.

## Try it yourself

1. Convert `Ninja` to use embedding: change the field to just `Human` (no name), and delete the `Chill()` forwarding method. Rebuild — does it still compile?
2. Add a `WithSword()` method to `SeniorNinja` that prints "Drew sword". Which struct does the method receiver need to be?
3. Write a small `main` function that creates a `SeniorNinja` and calls `Attack()`. What gets printed, and in what order? (Trace through the code before running.)

## Common pitfalls

- **Reaching for inheritance.** If you find yourself wishing you could `extends`, step back — usually the right Go answer is a smaller interface and multiple small types that satisfy it, not a deep hierarchy.
- **Value receivers vs pointer receivers.** All methods here use value receivers (`func (n Ninja)`). If a method needs to *modify* the struct, you use a pointer receiver (`func (n *Ninja)`). Mixing them on the same type is generally discouraged.
- **Method promotion doesn't cross packages the way you'd guess.** Embedded fields promote methods, but visibility still follows the export rules from lesson 02.

## You've understood this lesson when...

- You can explain in one sentence why Go's `Ninja` and `SeniorNinja` are not inheritance.
- You can convert this code to use embedding and know why it works.
- Given a struct, you can list which interfaces it satisfies without being told.

## Next

- **Next lesson:** [06-interfaces-and-mocking](../06-interfaces-and-mocking/) — how to design small interfaces at the point of use, and how to generate mocks for tests.
