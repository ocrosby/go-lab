# Interfaces and mocking

Small interfaces, the fake-first testing style Go rewards, and the narrow case where a generated mock is the right tool: at genuine system edges.

## Why it matters

Two Go-specific ideas set interfaces apart from other languages:

1. **Interfaces are satisfied structurally, not by declaration.** You never write `class Car implements Vehicle`. If a `Car` has all the methods `Vehicle` requires, it *is* a `Vehicle`, automatically. The compiler checks.
2. **Interfaces belong to the consumer.** The package that *uses* an interface defines it — the package that provides the concrete type does not. This lets each caller ask for exactly the small set of methods they need.

Those two ideas make **fakes** — small real-looking implementations of an interface — cheap and pleasant. This lesson teaches when to use a real implementation, when to use a small fake, and when to reach for a generated mock.

## The rule this lesson teaches

**Mock at the edges of your system, not at the edges of your classes.**

- A **class edge inside your own module** (e.g. `AccordFactory` uses `AccordBuilder`, both same-team, both in this deployable): use a real implementation. Test the outer thing through its public surface; the inner thing gets covered along the way. Mocking the inner one pins the shape of an internal collaboration and breaks tests on every refactor.
- A **system edge** (e.g. a fleet-tracking SaaS, a database, an S3 bucket, the clock): mocking or a small fake is legitimate. The interaction itself is the contract, and the tests verify that contract stays honoured.

This lesson has one of each so you can see the difference on the same page.

## Prerequisites

- Lesson 06: composition and interface satisfaction basics.
- Lesson 04: `go test` basics.

## Run it

```bash
go test -race ./lessons/07-interfaces-and-mocking/...
```

Expected output: three test packages, all `ok`. **25 tests** in total across `cars/honda` and `trucks/toyota`.

## What's in this folder

| Path | What it demonstrates |
|---|---|
| [`vehicle.go`](./vehicle.go) | `Vehicle` interface — nine methods that describe any vehicle in the module. |
| [`car.go`](./car.go), [`truck.go`](./truck.go) | Small interfaces embedding `Vehicle`. No mockgen directives — they're same-team. |
| [`telemetry.go`](./telemetry.go) | `TelemetryReporter` interface — the ONE mockgen'd interface. It's a system edge (fleet-tracking SaaS). |
| [`mocks/mock_telemetry.go`](./mocks/) | Generated mock, regenerated via `go generate ./...`. |
| [`cars/honda/accord.go`](./cars/honda/accord.go) | Concrete `Accord`. `TurnOn`/`TurnOff` emit telemetry when a reporter is wired. |
| [`cars/honda/accord_builder.go`](./cars/honda/accord_builder.go), [`accord_factory.go`](./cars/honda/accord_factory.go) | Builder + Factory for constructing Accords. Both are same-team; tests use them directly. |
| [`cars/honda/*_test.go`](./cars/honda/) | Black-box stdlib tests: public methods only, real collaborators, small fake at the telemetry edge. |
| [`trucks/toyota/`](./trucks/toyota/) | A `Tundra` mirror to reinforce the pattern. |

## The two tests worth reading side by side

**Real collaborator, no mock:**

```go
// accord_factory_test.go — same-team collaborators
func newFactory() *honda.AccordFactory {
    return honda.NewAccordFactory(honda.NewAccordBuilder())
}

func TestFactory_CreateWithYearSetsYear(t *testing.T) {
    got, _ := newFactory().CreateWithYear(2020)
    if got.GetYear() != 2020 { ... }
}
```

No `MockIAccordBuilder`. No `EXPECT().Build()`. The test wires the real factory to the real builder and asserts on what the returned `Accord` reports through its own public API. Refactor the builder however you like — swap the storage, extract a helper, inline `Build` — the tests still pass because nothing in them names the builder's shape.

**Edge contract, fake reporter:**

```go
// accord_test.go — system edge
type recordingReporter struct { events []ignitionEvent }
func (r *recordingReporter) ReportIgnition(vin string, on bool) error {
    r.events = append(r.events, ignitionEvent{vin, on})
    return nil
}

func TestTurnOn_ReportsIgnitionToTelemetry(t *testing.T) {
    reporter := &recordingReporter{}
    a := honda.NewAccordWithTelemetry("VIN123", reporter)

    _ = a.TurnOn()

    if len(reporter.events) != 1 || reporter.events[0].vin != "VIN123" { ... }
}
```

Here the assertion IS on the interaction — because the interaction *is* the contract with the external service. That's the legitimate case for interaction-verification testing.

## When would you use the generated mock?

The `mockgen`-generated `MockTelemetryReporter` in `mocks/` is available for tests that need strict ordering assertions, argument matchers, or expected-call counts that a hand-rolled fake would make verbose. Use it exactly where the fake would be more code than it's worth — never for a same-team collaborator.

Regenerate the mock with:

```bash
go install go.uber.org/mock/mockgen@latest
go generate ./lessons/07-interfaces-and-mocking/...
```

## The diagnostic

*If I change how any Go code in this module is written without changing what it does, do all the tests still pass?*

For lesson 07 as it stands now: **yes**. Extract methods in `AccordFactory`, rename fields in `AccordBuilder`, inline the ignition state machine — the tests survive because none of them name the shape of the code under test.

The old version of this lesson (before this refactor) failed this diagnostic because tests reached into `builder.instance`, `accord.state`, and `factory.builder` directly, and because `accord_factory_test.go` asserted on the exact sequence of `mockedBuilder.EXPECT().Build().GetInstance()` calls. Any renamed field or reordered call would have broken tests whose behaviour the code didn't touch. That is what the black-box style prevents.

## Try it yourself

1. Rename the `year` field in `Accord` to `modelYear`. Update the two `Get`/`Set` methods to match. Run `go test` — every test still passes.
2. Change `AccordBuilder` to store the in-progress `Accord` in a slice instead of a single field. Run `go test` — still passing.
3. Add a `LogTransport` fake that wraps `TelemetryReporter` and logs every call. Wire it into an Accord and verify the logs. (Hint: the fake pattern for TelemetryReporter is already in `accord_test.go`.)
4. Add a `RegistrationService` interface for looking up a VIN against a state DMV. Should it live in this package or one that owns the DMV boundary? Should it be mocked or faked?

## Common pitfalls

- **Mocking a same-team interface** because "it makes the test simpler." It makes the test simpler *once*; it makes every future refactor of that collaboration harder. Reach for the real thing.
- **Asserting on internal method calls** with `.EXPECT().SomeMethod()` when the real assertion should be on the outcome. If the outcome is right, the internal calls don't need to be verified.
- **Reaching into private fields** (`accord.state = "on"`) in test arrange sections. Every private field you touch is a chain around a future refactor. Use the public API to set up the state you want.
- **Testing every method individually** in its own `Describe` or `TestXxxMethod` function. Behaviour is what you test, not method decomposition. If a method has no observable behaviour on its own, extract or inline it — don't write a test for it.

## You've understood this lesson when...

- You can explain in one sentence why `AccordBuilder` is not mocked in this lesson but `TelemetryReporter` is.
- You can look at a Go interface and quickly say "this deserves a mock" or "this deserves a fake or a real."
- You can convert a mock-based test into a fake-based one without losing coverage.
- You know why "the test still passes after I renamed a private field" is the diagnostic you want.

## Related deep-dive

- [Testing behavior, not implementation](https://omarcrosby.com/posts/testing-behavior-not-implementation/) — the blog post this lesson's refactor was written against.

## Next

- **Next lesson:** [08-goroutines-and-channels](../08-goroutines-and-channels/) — Go's model of concurrency.
