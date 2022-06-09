package shortener_test

import (
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShortener(t *testing.T) {
	svc := shortener.New()
	data := []string{"https://ya.ru", "https://google.com"}

	firstID, err := svc.GenerateID()
	require.NoError(t, err)

	err = svc.Storage.Set(firstID, data[0])
	require.NoError(t, err)

	secondID, err := svc.GenerateID()
	require.NoError(t, err)

	err = svc.Storage.Set(secondID, data[1])
	require.NoError(t, err)

	_, err = svc.GenerateID()
	require.NoError(t, err)

	firstValue, err := svc.Storage.Get(firstID)
	require.NoError(t, err)
	assert.Equal(t, firstValue, data[0])

	secondValue, err := svc.Storage.Get(secondID)
	require.NoError(t, err)
	assert.Equal(t, secondValue, data[1])
}
