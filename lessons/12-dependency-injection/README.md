# Dependency injection

Passing behaviour into a type rather than hard-coding it — the way Go does it, using constructor functions and interfaces.

## Why it matters

"Dependency injection" sounds like a heavy pattern from Java frameworks, but in Go it's usually just: **a struct takes an interface as a constructor argument.** That one move makes code testable (pass a fake in tests), configurable (pass different behaviours at runtime), and decoupled (the caller decides what to plug in).

This is the same idea lesson 06 introduced with interfaces — this lesson names it and shows it applied to a small, whole program.

## Prerequisites

- Lesson 05: composition (structs holding other types).
- Lesson 06: interfaces satisfied structurally.

## Run it

```bash
go run ./lessons/12-dependency-injection
```

Expected output:

```text
Placing safeties
Placing NO safeties
Placing rock safeties
Placing ice safeties
```

Each line comes from a different **safety placer** injected into the same `RockClimber` type.

## What's in this folder

| Path | What it demonstrates |
|---|---|
| [`main.go`](./main.go) | Wiring — creates a `RockClimber` with each kind of safety placer and calls `ClumbRock()`. |
| [`pkg/rock_climber.go`](./pkg/rock_climber.go) | `RockClimber` struct that depends on a `SafetyPlacer` interface. |
| [`pkg/safety/safety_placer.go`](./pkg/safety/safety_placer.go) | The `SafetyPlacer` interface — one method. |
| [`pkg/safety/placers/`](./pkg/safety/placers/) | Four concrete implementations: concrete, none, rock, ice. |

## The pattern

Look at [`pkg/rock_climber.go`](./pkg/rock_climber.go):

```go
type RockClimber struct {
    rocksClimbed int
    sp           safety.SafetyPlacer   // interface — not a concrete type
}

func NewRockClimber(sp safety.SafetyPlacer) *RockClimber {
    return &RockClimber{sp: sp}
}

func (rc *RockClimber) ClumbRock() {
    rc.rocksClimbed++
    rc.sp.PlaceSafeties()
}
```

Three moves that add up to dependency injection:

1. **The struct field is an interface**, not a concrete type. `sp` can hold any type that has a `PlaceSafeties()` method.
2. **The constructor takes the interface as a parameter.** Callers pick which implementation to plug in.
3. **The behaviour is called through the field.** `rc.sp.PlaceSafeties()` — the `RockClimber` doesn't know *which* placer it has, just that it *has* one.

The main function then wires up different implementations for different scenarios:

```go
pkg.NewRockClimber(placers.ConcreteSafetyPlacer{}).ClumbRock()
pkg.NewRockClimber(placers.NOPSafetyPlacer{}).ClumbRock()
```

That's the pattern. In a real service you'd inject a database repository, an HTTP client, a metrics recorder, a clock — any behaviour that has multiple implementations or that you want to fake in tests.

## Why not use a framework?

Java has Spring, Python has FastAPI's `Depends`, .NET has built-in DI containers. Go has… `func NewX(deps ...) *X`. Because Go's interfaces are satisfied structurally (lesson 06), a manual constructor is enough for the vast majority of cases. There *are* Go DI libraries (`wire`, `uber-fx`, `google/wire`) for very large graphs, but reach for them only when hand-wiring becomes painful — usually deep into a production codebase.

Lesson 14 (`14-production-api`) shows hand-wired DI on a real service.

## Try it yourself

1. Write a new `RopeSafetyPlacer` that prints "Placing rope safeties". Add a line to `main.go` that uses it. No changes needed to `RockClimber`.
2. Add a `NumRocksClimbed() int` method to `RockClimber`. Modify `main.go` to print the count after all four climbs. (Note: each `NewRockClimber` call creates a *new* climber — the count is per-climber.)
3. Write a test for `RockClimber` that uses a fake `SafetyPlacer` that records how many times `PlaceSafeties` was called. Assert that calling `ClumbRock` once produces exactly one call.

## Common pitfalls

- **Injecting a concrete type instead of an interface.** If `NewRockClimber` took a `*ConcreteSafetyPlacer` directly, callers couldn't swap in a fake for testing. The interface is the seam.
- **A "fat" interface that most implementations don't fully implement.** If your interface has 12 methods and each mock has to implement all 12 (even the ones the test doesn't care about), you probably need to split the interface. See lesson 06.
- **Building the entire object graph in `main`.** That's actually the *correct* answer in small programs — but as programs grow, extracting a `wire()` function that returns the fully-wired root object keeps `main` readable.

## You've understood this lesson when...

- You can add a new `SafetyPlacer` implementation without touching `RockClimber` or any existing placer.
- You can explain in one sentence why the constructor takes an interface rather than a concrete type.
- You can write a fake implementation and use it in a test.

## Next

- **Next lesson:** [13-design-patterns](../13-design-patterns/) — classic OO design patterns rendered in idiomatic Go.
