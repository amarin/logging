package logging

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithContextExtractors(t *testing.T) {
	var key Key = "someKey"

	ctx := context.Background()

	t.Run("no key value in context by default", func(t *testing.T) {
		keys := CurrentConfig().contextKeys(ctx)
		require.NotContains(t, keys, key)
	})
	t.Run("no key value in context without extractor", func(t *testing.T) {
		updatedCtx := key.SetToCtx(ctx, nil)
		keys := CurrentConfig().contextKeys(updatedCtx)
		require.NotContains(t, keys, key)
	})
	t.Run("load key value having extractor", func(t *testing.T) {
		Init(WithContextExtractors(key.Extractor()))
		updatedCtx := key.SetToCtx(ctx, 1)
		keys := CurrentConfig().contextKeys(updatedCtx)
		require.Contains(t, keys, key)
	})
	t.Run("no nil value returned even with extractor", func(t *testing.T) {
		Init(WithContextExtractors(key.Extractor()))
		updatedCtx := key.SetToCtx(ctx, nil)
		keys := CurrentConfig().contextKeys(updatedCtx)
		require.NotContains(t, keys, key)
	})
}
