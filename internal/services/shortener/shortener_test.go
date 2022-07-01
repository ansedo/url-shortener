package shortener_test

import (
	"context"
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShortener(t *testing.T) {
	ctx := context.Background()
	svc := shortener.New(ctx)
	data := []string{"https://ya.ru", "https://google.com"}

	firstID, err := svc.GenerateID(ctx)
	require.NoError(t, err)

	err = svc.Storage.Add(ctx, firstID, data[0])
	require.NoError(t, err)

	secondID, err := svc.GenerateID(ctx)
	require.NoError(t, err)

	err = svc.Storage.Add(ctx, secondID, data[1])
	require.NoError(t, err)

	_, err = svc.GenerateID(ctx)
	require.NoError(t, err)

	firstValue, err := svc.Storage.GetByShortURL(ctx, firstID)
	require.NoError(t, err)
	assert.Equal(t, firstValue, data[0])

	secondValue, err := svc.Storage.GetByShortURL(ctx, secondID)
	require.NoError(t, err)
	assert.Equal(t, secondValue, data[1])
}
