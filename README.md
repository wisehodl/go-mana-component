# go-mana-component

Component identity propagation for layered Go libraries.

Source: https://git.wisehodl.dev/jay/go-mana-component
Mirror: https://github.com/wisehodl/go-mana-component

## What this library does

- Injects a named component identity into a `context.Context` at library boundaries
- Propagates module identity and component hierarchy across layers
- Implements `slog.LogValuer` for use as a structured log attribute, and provides a string map for generic consumers

## What this library does not do

`go-mana-component` does not log, does not define what a module or component
means semantically, and does not manage component lifecycles or dependencies.
It only carries identity. Interpretation, enforcement of naming conventions,
metrics, and tracing instrumentation all belong elsewhere.

## Installation

```bash
go get git.wisehodl.dev/jay/go-mana-component
```

If the primary repository is unavailable, use the `replace` directive in your `go.mod`:

```
replace git.wisehodl.dev/jay/go-mana-component => github.com/wisehodl/go-mana-component latest
```

## Usage

As an example, the `go-honeybee` library uses a three-layer component hierarchy
for pools, workers, and connections. To provide structured logging, each
component accepts a component-aware context and an `slog.Handler` and then
constructs a logger internally.

### At a library boundary

A top-level constructor receives a context and creates a new component
identity. Passing the component as a `slog.Any` group attribute allows the
logger to carry `module` and `path` automatically.

```go
func NewPool(ctx context.Context, id string, handler slog.Handler) (*Pool, error) {
    ctx = component.MustNew(ctx, "honeybee", "outbound_pool")

    c, _ := component.Get(ctx)
    logger := slog.New(handler).With(slog.Any("component", c), slog.String("pool_id", id))

    return &Pool{ctx: ctx, logger: logger}, nil
}
```

### Descending into a sub-component

A child constructor calls `MustExtend`, inheriting the module and extending the
path. No parent identifiers need to be passed as arguments.

```go
func NewWorker(ctx context.Context, id string, handler slog.Handler) (*Worker, error) {
    ctx = component.MustExtend(ctx, "outbound_worker")

    c, _ := component.Get(ctx)
    logger := slog.New(handler).With(slog.Any("component", c), slog.Any("peer_id", id))

    return &Worker{ctx: ctx, logger: logger}, nil
}
```

At the connection layer, another `MustExtend` call extends the path to
`outbound_pool.outbound_worker.connection` with no additional plumbing.

### Usage Notes

**Error-Returning vs Panic Variants**:

- Use `New` and `Extend` when a missing or invalid component is a recoverable
condition.
- Use `MustNew` and `MustExtend` at library boundaries where a missing
component is a programming error that should halt execution immediately.

**Generic Output**:

`GetFields` provides the component fields as a `map[string]string` for non-slog
consumers such as OpenTelemetry or Prometheus.

## Testing

```bash
go test ./...
```
