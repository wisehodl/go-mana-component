package component

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

type componentKey int

const storageKey componentKey = iota

// Component identifies a named unit within a module and its position in the call hierarchy.
type Component interface {
	Module() string
	Path() []string
	PathString() string
	LogValue() slog.Value
}

type component struct {
	module string
	path   []string
}

func (c component) Module() string     { return c.module }
func (c component) Path() []string     { return slices.Clone(c.path) }
func (c component) PathString() string { return strings.Join(c.path, ".") }
func (c component) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("module", c.module),
		slog.String("path", c.PathString()),
	)
}

func insert(ctx context.Context, module string, name string, path []string) context.Context {
	return context.WithValue(ctx, storageKey, component{
		module: module,
		path:   append(path, name),
	})
}

// Get retrieves the current component from the context; returns false if none is present.
func Get(ctx context.Context) (Component, bool) {
	t, ok := ctx.Value(storageKey).(component)
	return t, ok
}

// GetFields returns the component as a string map with keys "module" and "path".
func GetFields(ctx context.Context) (map[string]string, bool) {
	c, ok := Get(ctx)
	if !ok {
		return nil, false
	}

	return map[string]string{
		"module": c.Module(),
		"path":   c.PathString(),
	}, true
}

// New sets a new component on the context, resetting any existing component.
func New(ctx context.Context, module string, name string) (context.Context, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}
	if module == "" {
		return nil, fmt.Errorf("module cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	return insert(ctx, module, name, []string{}), nil
}

// Extend appends name to the current component path, inheriting its module.
func Extend(ctx context.Context, name string) (context.Context, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}

	c, ok := Get(ctx)
	if !ok {
		return nil, fmt.Errorf("missing parent component")
	}

	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	return insert(ctx, c.Module(), name, c.Path()), nil
}

// MustNew is New but panics on error.
func MustNew(ctx context.Context, module string, name string) context.Context {
	if ctx == nil {
		panic("context is nil")
	}
	if module == "" {
		panic("module cannot be empty")
	}
	if name == "" {
		panic("name cannot be empty")
	}
	return insert(ctx, module, name, []string{})
}

// MustExtend is Extend but panics on error.
func MustExtend(ctx context.Context, name string) context.Context {
	if ctx == nil {
		panic("context is nil")
	}

	c, ok := Get(ctx)
	if !ok {
		panic("missing parent component")
	}

	if name == "" {
		panic("name cannot be empty")
	}

	return insert(ctx, c.Module(), name, c.Path())
}
