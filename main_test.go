package component

import (
	"context"
	"slices"
	"testing"
)

func assertPanics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("expected panic, but none occurred")
		}
	}()
	f()
}

func TestNew(t *testing.T) {
	t.Run("creates a new component", func(t *testing.T) {
		ctx, err := New(context.Background(), "mymodule", "mycomponent")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		c := FromContext(ctx)
		if c == nil {
			t.Fatal("expected component in context")
		}
		if c.Module() != "mymodule" {
			t.Errorf("expected module %q, got %q", "mymodule", c.Module())
		}
		if !slices.Equal(c.Path(), []string{"mycomponent"}) {
			t.Errorf("expected path %v, got %v", []string{"mycomponent"}, c.Path())
		}
		if c.PathString() != "mycomponent" {
			t.Errorf("expected path string %q, got %q", "mycomponent", c.PathString())
		}
	})

	t.Run("should return error", func(t *testing.T) {
		_, err := New(nil, "mymodule", "mycomponent")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("must variant should panic", func(t *testing.T) {
		assertPanics(t, func() {
			MustNew(nil, "mymodule", "mycomponent")
		})
	})
}

func TestExtend(t *testing.T) {
	t.Run("extends existing component", func(t *testing.T) {
		ctx := MustNew(context.Background(), "mymodule", "mycomponent")
		ctx, err := Extend(ctx, "subcomponent")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		c := FromContext(ctx)
		if c == nil {
			t.Fatal("expected component in context")
		}
		if c.Module() != "mymodule" {
			t.Errorf("expected module %q, got %q", "mymodule", c.Module())
		}
		if !slices.Equal(c.Path(), []string{"mycomponent", "subcomponent"}) {
			t.Errorf("expected path %v, got %v", []string{"mycomponent", "subcomponent"}, c.Path())
		}
		if c.PathString() != "mycomponent.subcomponent" {
			t.Errorf("expected path string %q, got %q", "mycomponent.subcomponent", c.PathString())
		}
	})

	t.Run("should return error", func(t *testing.T) {
		_, err := Extend(context.Background(), "subcomponent")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("must variant should panic", func(t *testing.T) {
		assertPanics(t, func() {
			MustExtend(context.Background(), "subcomponent")
		})
	})
}

func TestFromContext(t *testing.T) {
	if FromContext(context.Background()) != nil {
		t.Fatal("expected no component on bare context")
	}

	ctx := MustNew(context.Background(), "mymodule", "mycomponent")
	if FromContext(ctx) == nil {
		t.Fatal("expected component in context")
	}
}

func TestFieldsFromContext(t *testing.T) {
	if FieldsFromContext(context.Background()) != nil {
		t.Fatal("expected no fields on bare context")
	}

	ctx := MustNew(context.Background(), "mymodule", "mycomponent")
	fields := FieldsFromContext(ctx)
	if fields == nil {
		t.Fatal("expected fields in context")
	}
	if fields["module"] != "mymodule" {
		t.Errorf("expected module %q, got %q", "mymodule", fields["module"])
	}
	if fields["path"] != "mycomponent" {
		t.Errorf("expected path %q, got %q", "mycomponent", fields["path"])
	}
}

func TestLogValue(t *testing.T) {
	ctx := MustNew(context.Background(), "mymodule", "mycomponent")
	c := FromContext(ctx)
	v := c.LogValue()
	attrs := v.Group()
	if len(attrs) != 2 {
		t.Fatalf("expected 2 attrs, got %d", len(attrs))
	}
	if attrs[0].Key != "module" || attrs[0].Value.String() != "mymodule" {
		t.Errorf("expected module=mymodule, got %s=%s", attrs[0].Key, attrs[0].Value.String())
	}
	if attrs[1].Key != "path" || attrs[1].Value.String() != "mycomponent" {
		t.Errorf("expected path=mycomponent, got %s=%s", attrs[1].Key, attrs[1].Value.String())
	}
}
