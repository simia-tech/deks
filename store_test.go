package edkvs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreSet(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	assert.Equal(t, e.storeOne.Len(), 1)
}

func TestStoreGet(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	value, err := e.storeOne.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, value)
}

func TestStoreDelete(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	require.NoError(t, e.storeOne.Delete(testKey))

	assert.Equal(t, e.storeOne.Len(), 0)
}
