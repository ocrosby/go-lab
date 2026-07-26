# JSON and struct tags

`encoding/json` — the standard-library package the syllabus's HTTP lessons lean on — and the **struct tags** that direct how Go values map to JSON fields.

> **Recommended before lesson 11 (HTTP) and lesson 14 (production API).** Both `json.NewDecoder(r.Body).Decode(&v)` and `json:"email"` struct tags appear in nearly every handler without any prior explanation.

## Why it matters

Almost every Go program that talks to another system uses JSON. The `encoding/json` package handles the encoding, but it needs guidance for anything non-trivial: renaming fields, omitting empty values, ignoring fields entirely, handling optional fields, catching typos. That guidance lives in struct tags — the string after the field type — and the safe decoding pattern lesson 16 uses.

## Prerequisites

- Lesson 19: types.
- Lesson 26: structs and methods.
- Lesson 23: error handling.

## Run it

```bash
go test -race ./lessons/30-json-and-struct-tags
```

Expected: 8 passes.

## What's in this folder

| File | What it demonstrates |
|---|---|
| [`json_demo.go`](./json_demo.go) | Marshal/Unmarshal, the four common struct-tag options, the safe streaming decoder pattern with `MaxBytesReader` + `DisallowUnknownFields`, and a custom `MarshalJSON` for a domain type. |
| [`json_demo_test.go`](./json_demo_test.go) | Tests exercising each shape. |

## Marshal and Unmarshal — the basics

```go
type User struct {
    ID    string
    Email string
    Age   int
}

// Struct → JSON bytes
b, err := json.Marshal(User{ID: "u1", Email: "a@b.com", Age: 36})
// b == `{"ID":"u1","Email":"a@b.com","Age":36}`

// JSON bytes → struct
var u User
err := json.Unmarshal(b, &u)
```

Notice: **field names are used as-is** by default. `ID` → `"ID"`. Public APIs usually want `"id"` — that's what struct tags fix.

## Struct tags

A struct tag is a **string** attached to a field. It looks like a comment but it's part of the language: `reflect` (which `encoding/json` uses internally) reads it.

```go
type User struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Age   int    `json:"age,omitempty"`
    Notes string `json:"-"`
}
```

- **`json:"id"`** — rename the field in JSON.
- **`json:"age,omitempty"`** — omit from output when the field has its zero value (0 for `int`, `""` for `string`, `nil` for pointers/slices/maps).
- **`json:"-"`** — never include this field in JSON.

`Marshal` of that struct produces:

```json
{"id":"u1","email":"a@b.com","age":36}
```

If `Age` is 0, `omitempty` drops it:

```json
{"id":"u1","email":"a@b.com"}
```

`Notes` never appears in either direction. Handy for secrets or internal-only fields.

## The safe decoding pattern

Naive `json.Unmarshal(body, &v)` has two production risks:

1. **Unbounded body** — a caller can send gigabytes; `Unmarshal` reads all of it. OOM vector.
2. **Silent unknown fields** — a caller sends `{"emial": "x"}` (typo) and your `Email` stays empty. Data loss with no error.

The safe pattern:

```go
// 1. Cap the body size.
r.Body = http.MaxBytesReader(w, r.Body, 1<<20)  // 1 MiB

// 2. Streaming decoder + strict-fields mode.
dec := json.NewDecoder(r.Body)
dec.DisallowUnknownFields()
if err := dec.Decode(&v); err != nil {
    return fmt.Errorf("invalid JSON: %w", err)
}

// 3. Reject requests with more than one JSON value.
if dec.More() {
    return errors.New("body must contain exactly one JSON value")
}
```

Lesson 16 uses this exact pattern in `decodeJSON`. Use it in every request handler.

## Common tag options in one table

| Tag | Effect |
|---|---|
| `json:"foo"` | Rename to `"foo"` in JSON. |
| `json:"foo,omitempty"` | Rename + omit when zero. |
| `json:",omitempty"` | Keep the field name; just omit when zero. |
| `json:"-"` | Skip this field always. |
| `json:"foo,string"` | Encode as a JSON string (useful for large ints that JS can't represent as numbers). |
| No tag | Use the exported field name as-is. |

## Custom `MarshalJSON` / `UnmarshalJSON`

When the default rules don't fit — say you want to render a `time.Duration` as `"5m"` instead of nanoseconds — implement the interfaces:

```go
type MyDur time.Duration

func (d MyDur) MarshalJSON() ([]byte, error) {
    return json.Marshal(time.Duration(d).String())
}

func (d *MyDur) UnmarshalJSON(b []byte) error {
    var s string
    if err := json.Unmarshal(b, &s); err != nil {
        return err
    }
    parsed, err := time.ParseDuration(s)
    if err != nil {
        return err
    }
    *d = MyDur(parsed)
    return nil
}
```

Every type that has these methods bypasses the reflection-driven default. Powerful but easy to over-use — reach for it only when the wire format genuinely differs from the Go type shape.

## Try it yourself

1. Add a `CreatedAt time.Time` field to `User` in `json_demo.go`. What does the JSON look like? (Answer: RFC 3339 format.)
2. Change one field's tag from `json:"email"` to `json:"email,omitempty"`. Marshal an empty-email `User`. Verify the field disappears.
3. Try `dec.DisallowUnknownFields()` off then on with a JSON body that has an extra field. Read the error you get with `on`.
4. Write a `Money` type (int cents) with a custom `MarshalJSON` that renders it as `"$12.34"`. Verify round-trip: JSON string → Money → same JSON string.
5. Try to marshal a struct with a channel field. What happens? (Channels, funcs, and complex types are unmarshalable — you'll get an error at runtime.)

## Common pitfalls

- **Typo in struct tag.** `json:"emial"` compiles fine — Go doesn't validate the field name inside the tag. Silently breaks JSON. Use `DisallowUnknownFields()` on inbound to catch caller typos, and `go vet`'s `structtag` check to catch your own.
- **Marshaling an unexported field.** Never happens — `encoding/json` only sees exported (capital-first-letter) fields. Rename or export the field if you need it in JSON.
- **Missing `&` on unmarshal.** `json.Unmarshal(b, u)` (value) vs `json.Unmarshal(b, &u)` (pointer). Only the pointer form works — the callee needs to mutate the caller's variable. Compile error if you forget on a value type, silent failure with an interface.
- **Unmarshaling numbers into `interface{}`.** They become `float64`. Reading `n.(int)` panics. Prefer typed struct fields or `json.Number` when you must accept unknown JSON.
- **Not closing `r.Body`.** After decoding an HTTP request body, still `defer r.Body.Close()` (or use MaxBytesReader which handles it). Leaks connections otherwise.
- **`omitempty` on a struct field.** `omitempty` checks a limited set of "zero" values — for a nested struct it does NOT check "all fields zero." Use a pointer (`*SubStruct`) so `nil` counts as empty.

## You've understood this lesson when...

- You can rename a Go field to snake-case in JSON with a struct tag.
- You know when to use `omitempty` and what "empty" means for each type.
- You can write the safe decoder pattern from memory (MaxBytesReader + DisallowUnknownFields + dec.More check).
- You know why `Notes string \`json:"-"\`` won't leak that field to clients.
- You can spot a `json:"emial"` typo bug in a code review.

## Next

- **Next lesson:** [31-type-assertions](../31-type-assertions/) — the `x.(T)` and `switch v := x.(type)` forms for extracting concrete types out of interface values.
