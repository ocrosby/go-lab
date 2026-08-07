package mocking

//go:generate mockgen -destination=./mocks/mock_telemetry.go -package=mocks github.com/ocrosby/go-lab/lessons/07-interfaces-and-mocking TelemetryReporter

// TelemetryReporter is the port a vehicle uses to send events out to an
// external fleet-management or telemetry SaaS. It is deliberately the ONLY
// interface in this lesson that gets a mockgen-generated mock.
//
// Why this interface, and not (say) IAccordBuilder or Vehicle?
//
//   - TelemetryReporter models a system EDGE — the wire protocol between
//     our code and a service outside our deployable. The interaction itself
//     is the contract (see the "adapter tests at genuine system edges"
//     exception in rules/black-box-testing.md). Asserting "when ignition
//     happens, exactly one ReportIgnition call goes out with the VIN and
//     timestamp" is a black-box test of that contract.
//
//   - IAccordBuilder and Vehicle are same-team collaborators inside the
//     same deployable. Mocking them pins the shape of an internal
//     collaboration and breaks tests on every refactor. Use real objects
//     and small in-memory fakes for these.
//
// Prefer this pattern in every codebase: mock at the edges of your system,
// not at the edges of your classes.
type TelemetryReporter interface {
	ReportIgnition(vin string, on bool) error
}
