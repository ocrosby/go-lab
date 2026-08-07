# Communicating Sequential Processes and Go

Go's approach to concurrency is often summarized in a single sentence: *"Do not communicate by sharing memory; share memory by communicating."* That slogan is not a marketing line — it names a specific tradition of concurrency theory called **Communicating Sequential Processes (CSP)**, first described by Tony Hoare in 1978. Go's goroutines, channels, and `select` statement are the language's practical realization of ideas that were originally a piece of formal mathematics.

This document explains what CSP is in plain language, why the Go designers chose it as the foundation for the language's concurrency model, and how the CSP ideas map to the Go constructs you use every day. It also honestly notes where Go departs from Hoare's original theory, and when the shared-memory-plus-locks model is still the right tool.

## Table of contents

- [What CSP is](#what-csp-is)
- [Why Go adopted CSP](#why-go-adopted-csp)
- [How CSP maps to Go](#how-csp-maps-to-go)
- [The slogan, unpacked](#the-slogan-unpacked)
- [Where Go diverges from Hoare's CSP](#where-go-diverges-from-hoares-csp)
- [Practical CSP-style patterns in Go](#practical-csp-style-patterns-in-go)
- [When shared memory is still the right tool](#when-shared-memory-is-still-the-right-tool)
- [Historical lineage](#historical-lineage)
- [References](#references)

## What CSP is

CSP is a way of describing systems made of many small, independent components that coordinate by **passing messages** rather than by reading and writing shared variables. The full theory is a formal algebra with rules for composing components, proving properties, and reasoning about deadlock and liveness. In plain terms, the core ideas are:

1. **A system is a collection of sequential processes.** Each process is a program that runs top-to-bottom, one step at a time. Nothing exotic — just ordinary sequential code.
2. **Processes are independent.** No process reaches into another's memory. Each has its own local state.
3. **Processes communicate over named channels.** A channel is a conduit for messages between exactly one sender and one receiver at a time (in the original theory).
4. **Communication is synchronization.** Sending a message and receiving that message happen together, at the same moment — a *rendezvous*. Neither side proceeds until both are ready.
5. **Composition is first-class.** Systems are built by running processes in parallel, in sequence, or by offering a choice between several possible next actions.

Under this model, all coordination is visible in the code as a message exchange. There is no invisible action-at-a-distance where one component changes a shared variable and another silently observes the change. The message passing *is* the coordination, and the pattern of message flow *is* the design.

Hoare introduced this model in a 1978 paper in *Communications of the ACM* and refined it in a 1985 book. Both remain in print and are the definitive references: [*Communicating Sequential Processes* (CACM 1978)](https://dl.acm.org/doi/10.1145/359576.359585); [*Communicating Sequential Processes* (Hoare, 1985, free PDF)](https://www.usingcsp.com/cspbook.pdf).

## Why Go adopted CSP

The Go team could have built concurrency on the model most languages use: threads plus shared memory plus mutexes. They deliberately chose not to. The reasons, drawn from talks and papers by Rob Pike and the Go team, are:

**Threads-and-locks does not scale in the mind of the programmer.** Reasoning about which lock protects which variable, in what order locks must be acquired, and where race conditions can hide, becomes intractable as a program grows. The bugs it produces (deadlocks, data races, torn reads) are among the hardest to reproduce and diagnose.

**Concurrent design should be composable.** A well-designed concurrent component in a threads-and-locks system is not automatically safe to combine with another well-designed component — the two may deadlock at the seam. CSP-style components, because their only interface is a channel, compose more predictably.

**The world is already concurrent.** Modern programs juggle network sockets, file I/O, timers, cancellation signals, and background work. A language whose primary tool for this is "spawn a thread and grab a mutex" makes every problem harder than it needs to be. Go was designed with server programming in mind, and the CSP model fits that shape naturally.

**A small number of primitives should cover most cases.** Go's designers wanted the concurrency toolkit to be small and orthogonal: goroutines for independent execution, channels for communication, `select` for choice, `context` for cancellation. Together these cover the vast majority of concurrent designs without introducing dozens of specialized synchronization primitives.

Rob Pike's 2012 talk [*Concurrency Is Not Parallelism*](https://go.dev/blog/waza-talk) is the canonical statement of the design philosophy. The related talk [*Go Concurrency Patterns*](https://go.dev/talks/2012/concurrency.slide) walks through the primitives with worked examples.

## How CSP maps to Go

The CSP-to-Go correspondence is direct and worth committing to memory.

| CSP concept | Go construct |
|---|---|
| Sequential process | Goroutine |
| Named channel | Value of type `chan T` |
| Message send (`c ! v`) | `c <- v` |
| Message receive (`c ? x`) | `x := <-c` |
| Synchronous rendezvous | Unbuffered channel |
| External choice (`P □ Q`) | `select` statement |
| Parallel composition (`P ‖ Q`) | Two `go` statements |
| Process termination | Goroutine returns; `close(c)` signals downstream |

A minimal example makes the mapping concrete:

```go
package main

import "fmt"

func main() {
    ch := make(chan string)         // create a channel
    go func() {                     // spawn a process
        ch <- "hello from process"  // send a message
    }()
    msg := <-ch                     // receive a message
    fmt.Println(msg)
}
```

Read this as CSP: two processes (the `main` goroutine and the one started by `go func`) share a channel `ch`. The child sends the string `"hello from process"`. The main process receives it. Because the channel is unbuffered, the send blocks until the receive happens — a textbook synchronous rendezvous. Once the message is exchanged, both sides proceed independently.

The `select` statement is Go's implementation of CSP's external choice: a process offers several possible communications and takes whichever one becomes ready first.

```go
select {
case msg := <-events:
    handle(msg)
case <-time.After(5 * time.Second):
    return errors.New("timed out waiting for event")
case <-ctx.Done():
    return ctx.Err()
}
```

Read as CSP: this process is willing to (a) receive an event, or (b) receive a timeout signal, or (c) receive a cancellation signal — whichever arrives first. The process does not busy-wait; the runtime parks the goroutine until one of the alternatives is ready.

## The slogan, unpacked

*"Do not communicate by sharing memory; share memory by communicating."*

The two halves refer to two different disciplines for coordination between concurrent components:

- **Communicating by sharing memory** — the traditional model. Two threads read and write a shared variable, and use a mutex to make the reads and writes atomic. The variable is the medium of communication; the lock exists to hide the fact that the medium is not safe.
- **Sharing memory by communicating** — the CSP model. A value moves from one goroutine to another over a channel. Ownership of the value passes with the message. At any moment, exactly one goroutine has the right to touch it.

The point of the slogan is that both approaches share memory in some sense — a Go channel is, underneath, a lock-protected queue in shared memory. What differs is which discipline the *programmer* uses. In the CSP model, you never write code that says "acquire this lock, mutate this field, release this lock." You write code that says "send this value to that goroutine." The synchronization comes for free from the channel operation itself.

This is why Go can have both `chan T` *and* `sync.Mutex` in the same standard library without contradiction. The channel is the preferred, higher-level tool; the mutex is available for the small number of cases where it genuinely fits better (see [When shared memory is still the right tool](#when-shared-memory-is-still-the-right-tool)).

## Where Go diverges from Hoare's CSP

Go was inspired by CSP; it is not a faithful implementation of the formal theory. The important differences:

**Channels are first-class values, not fixed names.** In Hoare's original notation, channels are static names in the program text — you write to `c` and `c` is a specific channel. In Go, channels are ordinary values that can be created, stored in variables, passed as arguments, returned from functions, and sent through other channels. This is a substantial ergonomic upgrade and enables patterns (like sending a reply-channel with a request) that pure CSP does not directly support.

**Buffered channels are an extension.** Hoare's CSP is fundamentally synchronous — every message exchange is a rendezvous. Go allows `make(chan T, N)` to create a channel with a buffer of size `N`; sends succeed without a matching receive as long as the buffer has room. This is convenient but subtly changes the reasoning. A buffered channel is closer to an asynchronous mailbox than to a CSP rendezvous, and it invalidates some of the deadlock-freedom arguments the pure theory allows.

**Goroutines are anonymous.** In CSP, processes have names and the parallel composition operator explicitly lists them. In Go, `go f()` spawns an anonymous unit of execution; there is no first-class handle for it. The `context` package and channels are how you keep track of goroutines you care about.

**No formal algebra.** Hoare's CSP comes with rewrite rules, trace semantics, and refinement checking (used by tools like [FDR](https://cocotec.io/fdr/)) that can prove properties like deadlock-freedom. Go has no equivalent. It borrows the *ideas* — sequential processes, channel communication, choice — without the mathematical machinery.

**Selective receive is choice on channels, not on senders.** In some CSP dialects (and in Erlang) a process can pattern-match on the *content* of incoming messages. Go's `select` chooses among channel operations, not among message contents. Content-based dispatch is done with an ordinary `switch` after the receive.

None of these are bugs. They are the tradeoffs of picking CSP-inspired concurrency for a pragmatic systems language rather than for a proof assistant.

## Practical CSP-style patterns in Go

The concurrency section of this repo has runnable examples of the patterns below. This section describes each in one paragraph so the theory-to-practice link is visible; the file `rules/go-concurrency.md` in the parent claude-config repo has more depth on each.

**Generator.** A function launches a goroutine and returns a receive-only channel. Callers iterate the channel with `for v := range ch`. This is a lazy sequence expressed as message passing — the producer runs concurrently with the consumer and hands values over one at a time.

**Fan-in.** Multiple input channels merge into one output channel by spawning one forwarding goroutine per input. This is CSP's parallel composition applied to a common problem: "combine several event streams into one."

**Fan-out worker pool.** A fixed number of worker goroutines read jobs from a shared channel and write results to another. Closing the input channel signals all workers to exit after draining. This is how CSP describes a bounded degree of parallelism.

**Pipeline.** Stages are connected by channels; each stage is a goroutine reading from its input channel and writing to its output. The whole pipeline is composed of small, independently-reasoned stages.

**Cancellation via `context.Context`.** The `Done()` method of a `Context` returns a channel that is closed when the context is cancelled. A goroutine selects on this channel alongside its normal work, giving cancellation exactly the same shape as any other CSP communication.

**Quit signal.** A `quit` channel passed to a goroutine, checked in a `select`, is the direct CSP idiom for orderly shutdown. The context pattern above is a more structured version of the same idea.

For deeper coverage of each pattern — with anti-patterns, closing rules, and leak-prevention checklists — see `rules/go-concurrency.md`.

## When shared memory is still the right tool

The CSP style is Go's default, not its law. The standard library ships `sync.Mutex`, `sync.RWMutex`, `sync.WaitGroup`, `sync.Once`, `sync/atomic`, and `sync.Map` precisely because there are workloads where they are the better fit:

- **A single struct with a handful of protected fields.** A `sync.Mutex` immediately above the fields it guards is often clearer than routing every access through a manager goroutine.
- **High-frequency read-mostly caches.** `sync.RWMutex` or `atomic.Pointer` is often significantly faster than channel-mediated access when the critical section is a nanosecond-scale read — a channel send/receive involves goroutine scheduling overhead that dwarfs a single atomic load.
- **Reference counting, one-shot initialization, atomic counters.** These map directly to `sync/atomic` and `sync.Once`. Wrapping them in a goroutine would add latency and complexity for no gain.
- **Existing shared-memory data structures from third-party libraries.** Wrapping them in channels for the sake of purity is usually not worth the code cost.

The Go proverb here is: *"Channels orchestrate; mutexes serialize."* Use channels when the design is fundamentally about passing values or coordinating stages of work. Use mutexes when the design is fundamentally about protecting a small piece of state from concurrent access.

## Historical lineage

The line from Hoare's paper to Go's `chan` runs through a specific series of languages, most of them from Bell Labs, and most of them involving Rob Pike:

- **1978** — Hoare publishes [*Communicating Sequential Processes*](https://dl.acm.org/doi/10.1145/359576.359585) in CACM.
- **1980s** — [Occam](https://en.wikipedia.org/wiki/Occam_(programming_language)), designed to run on the [Transputer](https://en.wikipedia.org/wiki/Transputer), is the first mainstream language built directly on CSP.
- **1985** — Hoare publishes the CSP book ([free PDF](https://www.usingcsp.com/cspbook.pdf)), formalizing the process algebra.
- **Mid-1980s** — Rob Pike creates [Newsqueak](https://en.wikipedia.org/wiki/Newsqueak) at Bell Labs, applying CSP ideas to interactive programs on the Blit terminal.
- **1990s** — Phil Winterbottom designs [Alef](https://en.wikipedia.org/wiki/Alef_(programming_language)) for Plan 9, drawing on Pike's Newsqueak work to bring CSP-style concurrency into a systems-programming language.
- **Late 1990s** — [Limbo](https://en.wikipedia.org/wiki/Limbo_(programming_language)), for the Inferno operating system, refines the model further with garbage collection and safe types.
- **2007–2009** — Robert Griesemer, Rob Pike, and Ken Thompson design Go at Google. The concurrency model is a direct descendant of the Newsqueak → Alef → Limbo line, with buffered channels, first-class channel values, and `select` inherited and refined.

Erlang, developed independently at Ericsson starting in the 1980s, is a parallel case: it uses message passing between isolated processes but with an asynchronous mailbox model closer to the [Actor model](https://en.wikipedia.org/wiki/Actor_model) than to pure CSP. Go and Erlang share design goals but arrived by different routes.

## References

### Foundational

- Hoare, C. A. R. *Communicating Sequential Processes*. CACM 21(8), 1978. [ACM Digital Library](https://dl.acm.org/doi/10.1145/359576.359585)
- Hoare, C. A. R. *Communicating Sequential Processes*. Prentice-Hall, 1985. [Free PDF](https://www.usingcsp.com/cspbook.pdf)

### Go-specific

- Rob Pike. *Concurrency Is Not Parallelism*. Waza 2012. [Blog post and video](https://go.dev/blog/waza-talk)
- Rob Pike. *Go Concurrency Patterns*. Google I/O 2012. [Slides](https://go.dev/talks/2012/concurrency.slide) · [Video](https://www.youtube.com/watch?v=f6kdp27TYZs)
- Sameer Ajmani. *Advanced Go Concurrency Patterns*. Google I/O 2013. [Slides](https://go.dev/talks/2013/advconc.slide) · [Video](https://www.youtube.com/watch?v=QDDwwePbDtw)
- The Go Blog. *Share Memory By Communicating*. [go.dev/doc/codewalk/sharemem](https://go.dev/doc/codewalk/sharemem/)
- The Go Memory Model. [go.dev/ref/mem](https://go.dev/ref/mem)

### Related theory

- Milner, R. *Communication and Concurrent Systems* (CCS). Prentice-Hall, 1989. A parallel formalism to CSP.
- The Actor model. [Wikipedia](https://en.wikipedia.org/wiki/Actor_model) — the message-passing model that inspired Erlang, contrasted with CSP.
- FDR: refinement checker for CSP. [cocotec.io/fdr](https://cocotec.io/fdr/) — used to prove properties of CSP designs.

### In this repo

- `rules/go-concurrency.md` in the parent claude-config repo — practical patterns, closing rules, and anti-patterns.
- `lessons/08-goroutines-and-channels/` — worked examples of goroutines, channels, `select`, and the primitive patterns.
