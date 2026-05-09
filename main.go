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

type Component interface {
	Module() string
	Path() []string
	PathString() string
}

type component struct {
	module string
	path   []string
}

func (c component) Module() string     { return c.module }
func (c component) Path() []string     { return slices.Clone(c.path) }
func (c component) PathString() string { return strings.Join(c.path, ".") }

func insert(ctx context.Context, module string, name string, path []string) context.Context {
	return context.WithValue(ctx, storageKey, component{
		module: module,
		path:   append(path, name),
	})
}

func Get(ctx context.Context) (Component, bool) {
	t, ok := ctx.Value(storageKey).(component)
	return t, ok
}

func MustStart(ctx context.Context, module string, name string) context.Context {
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

func MustNext(ctx context.Context, name string) context.Context {
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

func Start(ctx context.Context, module string, name string) (context.Context, error) {
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

func Next(ctx context.Context, name string) (context.Context, error) {
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

func Attrs(ctx context.Context) ([]slog.Attr, bool) {
	fields, ok := GetFields(ctx)
	if !ok {
		return nil, false
	}

	return []slog.Attr{
		slog.String("module", fields["module"]),
		slog.String("path", fields["path"]),
	}, true
}
