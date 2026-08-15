---
name: golang-testing
description: Use when asked to benchmark Go code, compare Go implementations for performance, add a fuzz test, or investigate allocations on a hot path - covers go test -bench/-benchmem, sub-benchmarks, benchstat for noise-resistant comparison, and native fuzzing (testing.F). For basic table-driven unit tests use go-service-idioms or test-driven-development instead; this skill is specifically about benchmark-driven optimization and property-style fuzz testing.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Go Benchmark-Driven Optimization and Fuzzing

Table-driven unit tests answer "is this correct." This skill answers two
different questions: "is this fast enough, and did my change actually help"
(benchmarks), and "does this hold for inputs I didn't think to write a test
case for" (fuzzing). Reach for `go-service-idioms` or
`test-driven-development` for ordinary correctness tests first — only bring
in benchmarks once a performance question is on the table, and fuzzing once
a function parses or validates untrusted/varied input.

## Writing a benchmark

```go
func BenchmarkProcess(b *testing.B) {
    data := generateTestData(1000)
    b.ResetTimer() // exclude setup cost above from the measured loop
    for i := 0; i < b.N; i++ {
        Process(data)
    }
}
```

Run with `-benchmem` always — throughput (`ns/op`) without allocation counts
(`B/op`, `allocs/op`) hides regressions where a change got faster per-op but
started allocating more, which shows up as GC pressure under real load:

```bash
go test -bench=BenchmarkProcess -benchmem ./...
# BenchmarkProcess-8   10000   105234 ns/op   4096 B/op   10 allocs/op
```

## Sub-benchmarks to compare alternatives directly

Use `b.Run` to benchmark several implementations side by side in one
invocation, so the comparison runs under identical conditions:

```go
func BenchmarkJoin(b *testing.B) {
    parts := []string{"hello", "world", "foo", "bar", "baz"}

    b.Run("plus_concat", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var s string
            for _, p := range parts {
                s += p
            }
            _ = s
        }
    })

    b.Run("strings_builder", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var sb strings.Builder
            for _, p := range parts {
                sb.WriteString(p)
            }
            _ = sb.String()
        }
    })
}
```

This is how a claim like "`strings.Builder` is faster than `+=` in a loop"
(see `golang-patterns`) gets verified rather than taken on faith — run both,
read the numbers, don't guess.

## benchstat: comparing runs without fooling yourself

A single benchmark run is noisy — CPU frequency scaling, other processes,
and GC timing all shift `ns/op` run to run. Never compare two single runs
by eye and declare a winner; use
[`golang.org/x/perf/cmd/benchstat`](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
to get a statistical comparison with confidence:

```bash
go test -bench=. -benchmem -count=10 ./... > old.txt
# make the change
go test -bench=. -benchmem -count=10 ./... > new.txt
benchstat old.txt new.txt
```

`-count=10` (or higher) is required — `benchstat` needs multiple samples per
side to compute variance; feeding it single-run output gives a meaningless
delta with no confidence indication.

## Native fuzzing (Go 1.18+)

Fuzzing generates inputs the developer didn't think to write by hand,
seeded from a small corpus, looking for panics, crashes, or a property
violation you assert explicitly inside `f.Fuzz`.

```go
func FuzzParseJSON(f *testing.F) {
    f.Add(`{"name": "test"}`) // seed corpus — real examples to start from
    f.Add(`[]`)
    f.Add(`""`)

    f.Fuzz(func(t *testing.T, input string) {
        var result map[string]any
        if err := json.Unmarshal([]byte(input), &result); err != nil {
            return // invalid JSON is an expected outcome, not a failure
        }
        // Property: anything that unmarshals successfully must re-marshal.
        if _, err := json.Marshal(result); err != nil {
            t.Errorf("Marshal failed after successful Unmarshal: %v", err)
        }
    })
}
```

```bash
go test -fuzz=FuzzParseJSON -fuzztime=30s
```

Always seed with `f.Add` using real, representative examples — an unseeded
fuzz target wastes early cycles rediscovering the input format from
scratch instead of exploring edge cases around valid input.

Fuzzing asserts a **property** (parse-then-marshal round-trips, a compare
function is antisymmetric, a parser never panics), not a specific expected
output — there's no table of expected values because the inputs aren't
enumerated by the developer. Any input that makes the fuzz function panic
or fail its assertion is automatically saved under `testdata/fuzz/<name>/`
as a permanent regression test.

## When to reach for which

| Question | Tool |
|---|---|
| Is the output correct for this input? | table-driven test (`go-service-idioms`) |
| Is it fast enough, and did my change help? | `go test -bench -benchmem` + `benchstat` |
| Does a property hold across inputs I didn't write by hand? | `testing.F` fuzz test |

## Gotchas

- `b.N` is chosen by the testing framework, not the author — it scales up
  automatically until the benchmark runs long enough to measure reliably.
  Never hardcode a loop count instead of using `b.N`, and always call
  `b.ResetTimer()` after any per-benchmark setup so setup cost isn't
  counted in the measured `ns/op`.
- A benchmark that never uses its result (`process(data)` with the return
  value discarded) risks the compiler optimizing the call away entirely
  under future compiler versions — assign to a package-level variable or
  use `testing.B.Elapsed`-style patterns defensively if the result is
  otherwise unused, so the benchmark keeps measuring real work.
- Corpus files under `testdata/fuzz/` are regression tests, not scratch
  files — commit them; deleting a saved failing corpus entry after "fixing"
  the bug removes the proof the fix actually holds against that input.
- `go test -fuzz=Name` runs **only** the named fuzz target and ignores
  everything else in the package during the fuzzing phase — run the full
  `go test ./...` separately to confirm nothing else regressed.

## Real-world grounding

Go's native fuzzing (`testing.F`, shipped in Go 1.18) evolved directly from
Dmitry Vyukov's external `go-fuzz` project, which found real bugs in the Go
standard library itself (`image/png`, `encoding/json`, `compress/flate`)
years before fuzzing was built into `go test` — the tool exists because
hand-written table-driven cases systematically miss the malformed,
adversarial inputs that fuzzing generates automatically.
