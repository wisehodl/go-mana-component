package component

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("creates a new component", func(t *testing.T) {
		ctx, err := New(context.Background(), "mymodule", "mycomponent")
		assert.NoError(t, err)
		c, ok := Get(ctx)
		assert.True(t, ok)
		assert.Equal(t, "mymodule", c.Module())
		assert.Equal(t, []string{"mycomponent"}, c.Path())
		assert.Equal(t, "mycomponent", c.PathString())
	})

	t.Run("should return error", func(t *testing.T) {
		_, err := New(nil, "mymodule", "mycomponent")
		assert.Error(t, err)
	})

	t.Run("must variant should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNew(nil, "mymodule", "mycomponent")
		})
	})
}

func TestExtend(t *testing.T) {
	t.Run("extends existing component", func(t *testing.T) {
		ctx := MustNew(context.Background(), "mymodule", "mycomponent")
		ctx, err := Extend(ctx, "subcomponent")
		assert.NoError(t, err)
		c, ok := Get(ctx)
		assert.True(t, ok)
		assert.Equal(t, "mymodule", c.Module())
		assert.Equal(t, []string{"mycomponent", "subcomponent"}, c.Path())
		assert.Equal(t, "mycomponent.subcomponent", c.PathString())
	})

	t.Run("should return error", func(t *testing.T) {
		_, err := Extend(context.Background(), "subcomponent")
		assert.Error(t, err)
	})

	t.Run("must variant should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			MustExtend(context.Background(), "subcomponent")
		})
	})
}

func TestGet(t *testing.T) {
	_, ok := Get(context.Background())
	assert.False(t, ok)

	ctx := MustNew(context.Background(), "mymodule", "mycomponent")
	_, ok = Get(ctx)
	assert.True(t, ok)
}

func TestGetFields(t *testing.T) {
	_, ok := GetFields(context.Background())
	assert.False(t, ok)

	ctx := MustNew(context.Background(), "mymodule", "mycomponent")
	fields, ok := GetFields(ctx)
	assert.True(t, ok)
	assert.Equal(t, "mymodule", fields["module"])
	assert.Equal(t, "mycomponent", fields["path"])
}

func TestAttrs(t *testing.T) {
	ctx := MustNew(context.Background(), "mymodule", "mycomponent")
	attrs, ok := Attrs(ctx)
	assert.True(t, ok)
	assert.Equal(t, slog.String("module", "mymodule"), attrs[0])
	assert.Equal(t, slog.String("path", "mycomponent"), attrs[1])
}
