---
name: gin-lambda-api-service-patterns
description: Guides Go API services built with Gin and deployed as a serverless function (AWS Lambda via an API-Gateway adapter) alongside a local HTTP-server dev mode. Use when touching a service's entrypoint, its dual local/serverless run paths, its build-and-deploy pipeline, or reviewing why a change that works locally behaves differently once deployed.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "go"
---

# Gin + Lambda Dual-Runtime API Service

A common shape for a Go API service targeting serverless deployment: one
Gin engine, wired to two different runtimes depending on environment — a
plain `http.Server` for local development, and a Lambda/API-Gateway
adapter (e.g. `github.com/apex/gateway` or `github.com/aws/aws-lambda-go`'s
`httpadapter`) for the deployed path. The two paths share the same route
handlers but are genuinely different code, not the same server started two
ways — this skill is about the places that difference actually matters.

## The dual entrypoint

```go
func main() {
    router := setupRouter()

    if os.Getenv("APP_ENV") == "production" {
        // Serverless path: no listener, no graceful-shutdown signal handling —
        // the Lambda runtime owns the process lifecycle.
        log.Fatal(gateway.ListenAndServe(port, router))
    } else {
        // Local dev path: a real listener, real graceful shutdown.
        srv := &http.Server{Addr: ":" + port, Handler: router}
        go func() {
            if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                log.Fatal(err)
            }
        }()
        waitForShutdownSignal(srv)
    }
}
```

Treat these as two genuinely different runtime paths, not one server
started two ways: the Lambda-adapter path typically has no graceful
shutdown or OS signal handling (the platform owns process lifecycle, cold
starts and freezes replace long-running-process concerns), while the local
path does. A change to shutdown behavior, connection draining, or
long-lived background goroutines needs to be reasoned about separately for
each path — testing only the local path can miss a change that behaves
differently once actually deployed as a function.

## Build and deploy shape

A Lambda-target Go service is typically built as a static cross-compiled
binary and zipped, not built as a container image:

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ./bootstrap ./cmd
zip main.zip bootstrap
aws lambda update-function-code --function-name <fn> --zip-file fileb://./main.zip
```

If the repo also has a `Dockerfile`, confirm what it's actually for before
assuming it drives the deployed artifact — it may only be used for local
development parity (a container matching the Lambda runtime for local
testing) while the real deploy path is the zip-and-update-function-code
flow above, or vice versa. Don't assume container-based deployment
tooling (ECS/Kubernetes-style health checks, readiness probes) applies
just because a `Dockerfile` exists in a Lambda-targeting repo.

## Version-of-truth drift

A Go version can be declared in several places that silently disagree:
`go.mod`'s `go` directive, a `Dockerfile`'s base image tag, a build
pipeline's install step, and a README. Before "fixing" a build issue by
changing one of these, check all of them — a project can genuinely be
running an older toolchain in CI than `go.mod` declares, and changing only
the file that happens to be open risks fixing the symptom in one place
while leaving the actual build environment unchanged.

## Profiling and instrumentation availability differs by path

A `pprof` HTTP endpoint, or any other instrumentation that depends on a
long-lived listener, is typically only reachable in the local
`http.Server` path — the Lambda/serverless path has no equivalent
always-on listener to attach it to. Don't assume a debugging tool
available locally is also available in the deployed environment without
checking whether it depends on the listener that only exists in the local
runtime branch.

## Gotchas

- **The two runtime paths are not interchangeable** — a fix or change
  tested only through the local `http.Server` path can behave differently
  (or not run at all, e.g. graceful-shutdown logic) once actually
  deployed through the serverless adapter path.
- **A `Dockerfile`'s presence doesn't imply container-based deployment** —
  confirm what actually builds and ships the deployed artifact
  (cross-compiled zip vs. container image) before assuming which
  deployment model applies.
- **Version-of-truth can disagree across `go.mod`, `Dockerfile`, CI
  config, and README simultaneously** — check all of them before treating
  any single one as authoritative when diagnosing a build issue.
- **Local-only tooling (pprof, live-reload) silently doesn't exist in the
  deployed path** if it depends on a listener that only the local runtime
  branch creates — don't assume parity between "works locally" and
  "available once deployed."
- **A serverless entrypoint typically has no graceful-shutdown/signal
  handling by design** — the platform owns process lifecycle; don't add
  local-style shutdown logic to that path expecting it to run, and don't
  assume the absence of shutdown handling there is an oversight to fix.

## Real-world grounding

The Gin-plus-Lambda-adapter shape (`apex/gateway`, or the AWS-maintained
`aws-lambda-go-api-proxy`) is a standard, widely-used pattern for running
an existing Go HTTP router on Lambda without rewriting handlers against a
Lambda-native request/response type — the adapter translates API Gateway
events into the same `http.Request`/`http.ResponseWriter` interface the
router already expects. The tradeoff this pattern accepts deliberately is
exactly the loss of long-running-process assumptions (graceful shutdown,
a persistent listener for profiling) that a conventional server has, in
exchange for the platform managing scaling and lifecycle.

## Verification

- [ ] Changes to shutdown, connection draining, or background-goroutine
      behavior have been reasoned about for both the local and serverless
      runtime paths, not just the one under active development
- [ ] The actual build/deploy artifact (zip vs. container image) is
      confirmed before assuming which deployment model governs the
      service
- [ ] Go version is checked across every place it's declared
      (`go.mod`, `Dockerfile`, CI config, README) before "fixing" a
      version-related build issue in only one of them
- [ ] Debugging/profiling tooling is confirmed available in the actual
      deployed runtime path before relying on it, not assumed present
      because it works locally
