# Design Patterns in Go

How classical design patterns are implemented and adapted for Go's unique features.

## Overview

Go is not a traditional object-oriented language, but it supports many object-oriented patterns through interfaces, embedding, and composition. This document explains how to implement design patterns effectively in Go.

## Go's Unique Features for Patterns

### Interfaces
- Implicit implementation (no `implements` keyword)
- Small, focused interfaces are preferred
- Enable polymorphism and dependency injection

### Embedding
- Composition over inheritance
- Anonymous field embedding
- Method promotion and overriding

### First-Class Functions
- Functions as values
- Closures for state encapsulation
- Strategy pattern through function types

## Creational Patterns

### Builder Pattern

**Problem**: Constructing complex objects step by step.

**Go Implementation**:
```go
// Product
type House struct {
    WindowType string
    DoorType   string
    NumFloors  int
}

// Builder interface
type IBuilder interface {
    SetWindowType(string) IBuilder
    SetDoorType(string) IBuilder  
    SetNumFloors(int) IBuilder
    Build() House
}

// Concrete builder
type ConcreteBuilder struct {
    windowType string
    doorType   string
    numFloors  int
}

func (b *ConcreteBuilder) SetWindowType(windowType string) IBuilder {
    b.windowType = windowType
    return b
}

func (b *ConcreteBuilder) SetDoorType(doorType string) IBuilder {
    b.doorType = doorType
    return b
}

func (b *ConcreteBuilder) SetNumFloors(floors int) IBuilder {
    b.numFloors = floors
    return b
}

func (b *ConcreteBuilder) Build() House {
    return House{
        WindowType: b.windowType,
        DoorType:   b.doorType,
        NumFloors:  b.numFloors,
    }
}

// Usage
func main() {
    house := &ConcreteBuilder{}
    house.SetWindowType("Wooden Window").
          SetDoorType("Wooden Door").
          SetNumFloors(2).
          Build()
}
```

**Go Advantages**:
- Method chaining through return values
- Interface-based design for multiple builders
- No complex inheritance hierarchies

### Singleton Pattern

**Problem**: Ensure only one instance exists.

**Go Implementation**:
```go
package singleton

import (
    "sync"
)

type Singleton struct {
    data string
}

var instance *Singleton
var once sync.Once

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{
            data: "This is a singleton instance",
        }
    })
    return instance
}

func (s *Singleton) GetData() string {
    return s.data
}
```

**Go Advantages**:
- `sync.Once` ensures thread-safe initialization
- No complex locking mechanisms needed
- Package-level variables for global access

### Factory Pattern

**Problem**: Create objects without specifying exact classes.

**Go Implementation**:
```go
// Product interface
type Animal interface {
    Speak() string
}

// Concrete products
type Dog struct{}
func (d Dog) Speak() string { return "Woof!" }

type Cat struct{}
func (c Cat) Speak() string { return "Meow!" }

// Factory function
func CreateAnimal(animalType string) Animal {
    switch animalType {
    case "dog":
        return &Dog{}
    case "cat":
        return &Cat{}
    default:
        return nil
    }
}

// Usage
func main() {
    dog := CreateAnimal("dog")
    fmt.Println(dog.Speak()) // Output: Woof!
}
```

## Structural Patterns

### Adapter Pattern

**Problem**: Make incompatible interfaces work together.

**Go Implementation**:
```go
// Target interface (what client expects)
type Computer interface {
    InsertIntoLightningPort()
}

// Adaptee (what we have)
type Windows struct{}

func (w *Windows) InsertIntoUSBPort() {
    fmt.Println("USB connector plugged into Windows machine")
}

// Adapter
type WindowsAdapter struct {
    windowsMachine *Windows
}

func (wa *WindowsAdapter) InsertIntoLightningPort() {
    fmt.Println("Adapter converts Lightning to USB")
    wa.windowsMachine.InsertIntoUSBPort()
}

// Client
type Mac struct{}

func (m *Mac) InsertIntoLightningPort() {
    fmt.Println("Lightning connector plugged into Mac machine")
}

// Client code
func connectDevice(com Computer) {
    com.InsertIntoLightningPort()
}

func main() {
    mac := &Mac{}
    connectDevice(mac)

    windows := &Windows{}
    windowsAdapter := &WindowsAdapter{
        windowsMachine: windows,
    }
    connectDevice(windowsAdapter)
}
```

### Decorator Pattern (via Embedding)

**Problem**: Add behavior to objects dynamically.

**Go Implementation**:
```go
// Component interface
type Coffee interface {
    Cost() int
    Description() string
}

// Concrete component
type SimpleCoffee struct{}

func (c SimpleCoffee) Cost() int {
    return 2
}

func (c SimpleCoffee) Description() string {
    return "Simple coffee"
}

// Decorator
type MilkDecorator struct {
    Coffee
}

func (m MilkDecorator) Cost() int {
    return m.Coffee.Cost() + 1
}

func (m MilkDecorator) Description() string {
    return m.Coffee.Description() + " with milk"
}

// Sugar decorator
type SugarDecorator struct {
    Coffee
}

func (s SugarDecorator) Cost() int {
    return s.Coffee.Cost() + 1
}

func (s SugarDecorator) Description() string {
    return s.Coffee.Description() + " with sugar"
}

// Usage
func main() {
    coffee := SimpleCoffee{}
    coffeeWithMilk := MilkDecorator{coffee}
    coffeeWithMilkAndSugar := SugarDecorator{coffeeWithMilk}
    
    fmt.Printf("%s costs $%d\n", 
        coffeeWithMilkAndSugar.Description(), 
        coffeeWithMilkAndSugar.Cost())
}
```

## Behavioral Patterns

### Strategy Pattern

**Problem**: Select algorithms at runtime.

**Go Implementation**:
```go
// Strategy interface
type PaymentStrategy interface {
    Pay(amount int) string
}

// Concrete strategies
type CreditCardStrategy struct {
    CardNumber string
}

func (c CreditCardStrategy) Pay(amount int) string {
    return fmt.Sprintf("Paid $%d using Credit Card ending in %s", 
        amount, c.CardNumber[len(c.CardNumber)-4:])
}

type PayPalStrategy struct {
    Email string
}

func (p PayPalStrategy) Pay(amount int) string {
    return fmt.Sprintf("Paid $%d using PayPal account %s", amount, p.Email)
}

// Context
type PaymentContext struct {
    strategy PaymentStrategy
}

func (p *PaymentContext) SetStrategy(strategy PaymentStrategy) {
    p.strategy = strategy
}

func (p *PaymentContext) ExecutePayment(amount int) string {
    return p.strategy.Pay(amount)
}

// Usage
func main() {
    context := &PaymentContext{}
    
    // Use credit card
    context.SetStrategy(CreditCardStrategy{CardNumber: "1234567890123456"})
    fmt.Println(context.ExecutePayment(100))
    
    // Switch to PayPal
    context.SetStrategy(PayPalStrategy{Email: "user@example.com"})
    fmt.Println(context.ExecutePayment(100))
}
```

### Observer Pattern

**Problem**: Notify multiple objects about state changes.

**Go Implementation**:
```go
// Observer interface
type Observer interface {
    Update(string)
    GetID() string
}

// Subject interface
type Subject interface {
    Register(Observer)
    Unregister(Observer)
    Notify(string)
}

// Concrete subject
type NewsAgency struct {
    observers []Observer
    news      string
}

func (n *NewsAgency) Register(observer Observer) {
    n.observers = append(n.observers, observer)
}

func (n *NewsAgency) Unregister(observer Observer) {
    for i, obs := range n.observers {
        if obs.GetID() == observer.GetID() {
            n.observers = append(n.observers[:i], n.observers[i+1:]...)
            break
        }
    }
}

func (n *NewsAgency) Notify(news string) {
    for _, observer := range n.observers {
        observer.Update(news)
    }
}

func (n *NewsAgency) SetNews(news string) {
    n.news = news
    n.Notify(news)
}

// Concrete observer
type NewsChannel struct {
    id   string
    news string
}

func (nc *NewsChannel) Update(news string) {
    nc.news = news
    fmt.Printf("NewsChannel %s received: %s\n", nc.id, news)
}

func (nc *NewsChannel) GetID() string {
    return nc.id
}

// Usage
func main() {
    agency := &NewsAgency{}
    
    channel1 := &NewsChannel{id: "CNN"}
    channel2 := &NewsChannel{id: "BBC"}
    
    agency.Register(channel1)
    agency.Register(channel2)
    
    agency.SetNews("Breaking: Go 1.22 Released!")
}
```

## Go-Specific Patterns

### Functional Options Pattern

**Problem**: Flexible configuration of objects.

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
}

type Option func(*Server)

func WithHost(host string) Option {
    return func(s *Server) {
        s.host = host
    }
}

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func WithTimeout(timeout time.Duration) Option {
    return func(s *Server) {
        s.timeout = timeout
    }
}

func NewServer(options ...Option) *Server {
    server := &Server{
        host:    "localhost",
        port:    8080,
        timeout: 30 * time.Second,
    }
    
    for _, option := range options {
        option(server)
    }
    
    return server
}

// Usage
func main() {
    server := NewServer(
        WithHost("0.0.0.0"),
        WithPort(9090),
        WithTimeout(60*time.Second),
    )
}
```

### Worker Pool Pattern

**Problem**: Limit concurrent processing.

```go
type Job struct {
    ID   int
    Data string
}

type Result struct {
    Job Job
    Err error
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        // Simulate work
        time.Sleep(time.Second)
        
        results <- Result{
            Job: job,
            Err: nil,
        }
    }
}

func main() {
    const numWorkers = 3
    const numJobs = 10
    
    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)
    
    // Start workers
    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }
    
    // Send jobs
    for j := 1; j <= numJobs; j++ {
        jobs <- Job{ID: j, Data: fmt.Sprintf("data-%d", j)}
    }
    close(jobs)
    
    // Collect results
    for a := 1; a <= numJobs; a++ {
        result := <-results
        fmt.Printf("Job %d completed\n", result.Job.ID)
    }
}
```

## Best Practices

### Interface Design
1. **Keep interfaces small**: Prefer many small interfaces over few large ones
2. **Accept interfaces, return structs**: Be liberal in what you accept
3. **Use standard library interfaces**: `io.Reader`, `io.Writer`, etc.

### Pattern Selection
1. **Favor composition**: Use embedding over complex inheritance
2. **Use channels for communication**: Don't communicate by sharing memory
3. **Leverage Go's concurrency**: Many patterns become simpler with goroutines

### Common Mistakes
1. **Over-engineering**: Don't implement patterns you don't need
2. **Fighting the language**: Adapt patterns to Go's strengths
3. **Ignoring interfaces**: They're key to flexible design in Go

## Pattern Examples in This Repository

- **Builder**: [`/learning/advanced/patterns/creational/builder/`](../../learning/advanced/patterns/creational/builder/)
- **Singleton**: [`/learning/advanced/patterns/creational/singleton/`](../../learning/advanced/patterns/creational/singleton/)
- **Adapter**: [`/learning/advanced/patterns/structural/adapter/`](../../learning/advanced/patterns/structural/adapter/)
- **Dependency Injection**: [`/learning/advanced/dependency-injection/`](../../learning/advanced/dependency-injection/)

## Further Reading

- [Go Design Patterns Repository Examples](../../learning/advanced/patterns/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Proverbs](https://go-proverbs.github.io/)

Remember: In Go, simple solutions are often better than complex patterns. Use patterns to solve real problems, not to show off your knowledge of design patterns.