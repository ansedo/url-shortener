package shortener_test

import (
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"github.com/ansedo/url-shortener/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShortener(t *testing.T) {
	testData := []string{"https://ya.ru", "https://google.com"}
	testShortener := shortener.NewShortener(memory.NewStorage())

	firstID, err := testShortener.GenerateID()
	require.NoError(t, err)

	err = testShortener.Storage.Set(firstID, testData[0])
	require.NoError(t, err)

	secondID, err := testShortener.GenerateID()
	require.NoError(t, err)

	err = testShortener.Storage.Set(secondID, testData[1])
	require.NoError(t, err)

	_, err = testShortener.GenerateID()
	require.NoError(t, err)

	firstValue, err := testShortener.Storage.Get(firstID)
	require.NoError(t, err)
	assert.Equal(t, firstValue, testData[0])

	secondValue, err := testShortener.Storage.Get(secondID)
	require.NoError(t, err)
	assert.Equal(t, secondValue, testData[1])
}
