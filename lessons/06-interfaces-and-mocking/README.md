# Interfaces and mocking

Small interfaces, generated mocks, and Go's rule for who defines an interface.

## Why it matters

Interfaces are how Go handles polymorphism — one function that works with many types. Two Go-specific ideas set this apart from other languages:

1. **Interfaces are satisfied structurally, not by declaration.** You never write `class Car implements Vehicle`. If a `Car` has all the methods `Vehicle` requires, it *is* a `Vehicle`, automatically. The compiler checks.
2. **Interfaces belong to the consumer, not the producer.** The package that *uses* an interface defines it — the package that provides the concrete type does not. This lets each caller ask for exactly the small set of methods they need.

Once you get these two, mocking for tests becomes trivial: any test can invent an interface that matches what the code under test needs, and pass a fake implementation.

## Prerequisites

- Lesson 05: composition and interface satisfaction basics.
- Comfort with `go test`.

## Run it

```bash
go test ./lessons/06-interfaces-and-mocking/...
```

The trailing `/...` tells Go "this package **and** all sub-packages" — necessary because the demo code lives under `cars/honda/` and `trucks/toyota/`.

Expected output (last line):

```text
ok  	github.com/ocrosby/go-lab/lessons/06-interfaces-and-mocking/cars/honda	...
ok  	github.com/ocrosby/go-lab/lessons/06-interfaces-and-mocking/trucks/toyota	...
```

## What's in this folder

| Path | What it demonstrates |
|---|---|
| [`vehicle.go`](./vehicle.go) | A `Vehicle` interface with nine methods. |
| [`car.go`](./car.go), [`truck.go`](./truck.go) | Two interfaces that embed `Vehicle`. Also carry `//go:generate` directives for mockgen. |
| [`cars/honda/`](./cars/honda/) | A concrete `Accord` implementation with tests, plus **generated mocks** (`mock_*.go`). |
| [`trucks/toyota/`](./trucks/toyota/) | A `Tundra` implementation and tests. |
| [`mocks/`](./mocks/) | Package-level generated mocks for `Car` and other interfaces. |

## Structural typing, quickly

```go
type Vehicle interface {
    TurnOn() error
    // ... eight more methods
}

type Accord struct { ... }

func (a *Accord) TurnOn() error { ... }
// ... eight more methods
```

Because `*Accord` has every method `Vehicle` requires, `*Accord` **is** a `Vehicle`. Nothing else to write.

Compare to Java: `public class Accord implements Vehicle`. In Go, the `implements` clause doesn't exist — the compiler figures it out from the method set.

**Consequence for design:** you can retroactively make an existing type satisfy a new interface. Someone else's `Vehicle` interface can be satisfied by *your* struct without either package importing the other. This is why Go interfaces are typically **small** — one to three methods, defined at the point of use.

## Generated mocks

Look at the top of [`car.go`](./car.go):

```go
//go:generate mockgen -destination=./mocks/mock_car.go -package=mocks github.com/ocrosby/go-lab/lessons/06-interfaces-and-mocking Car
```

When you run `go generate ./...`, this line becomes a shell command that generates [`mocks/mock_car.go`](./mocks/) — a fake `Car` you can use in tests to control what methods return.

Setup:

```bash
go install go.uber.org/mock/mockgen@latest
go generate ./lessons/06-interfaces-and-mocking/...
```

The generated file is committed to the repo (`mock_car.go`), so contributors don't have to install `mockgen` just to run tests. You only need to re-run `go generate` if you change the interface.

A test using the mock looks like:

```go
ctrl := gomock.NewController(t)
mockCar := mocks.NewMockCar(ctrl)
mockCar.EXPECT().TurnOn().Return(nil)  // "when TurnOn is called, return nil"

drive(mockCar)  // pass the fake into the code under test
```

See [`cars/honda/accord_test.go`](./cars/honda/accord_test.go) for real usage.

## Try it yourself

1. Add a method `GetColor() string` to `Vehicle`. Which files fail to compile? (This shows you the cost of a fat interface.)
2. Split `Vehicle` into two smaller interfaces: `Ignitable` (has `TurnOn`/`TurnOff`/`GetState`/`SetState`) and `Identifiable` (make/model/year). Which callers only need one of the two?
3. Add a `MotorcycleInterface` in a new file and write a `Harley` struct that satisfies it. No changes needed to `vehicle.go`.

## Common pitfalls

- **Fat interfaces.** `Vehicle` here has nine methods — that's a lot. Real Go code prefers small interfaces (often just one or two methods). This lesson is fat on purpose so you see the maintenance cost when you have to add a tenth.
- **Defining interfaces where the type lives.** If you're writing package `car` and defining `type Car interface`, you may have the direction backwards. Interfaces belong where they are *used*, not where the concrete type lives.
- **Forgetting `//go:generate`.** If you change the interface but don't re-run `go generate`, your mocks are stale and your tests may pass against yesterday's interface.

## You've understood this lesson when...

- You can explain in one sentence why Go doesn't need an `implements` keyword.
- You can add a new implementation of `Vehicle` in a new package without editing the interface file.
- You can wire up a mock in a test with `mockgen` without looking it up.

## Related deep-dive

- [`docs/go-build-directives.md`](../../docs/go-build-directives.md) — the full story on `//go:generate` and the rest of the `//go:` directive family.

## Next

- **Next lesson:** [07-goroutines-and-channels](../07-goroutines-and-channels/) — Go's model of concurrency.
