package edkvs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeReconcilateValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	count, err := e.nodeTwo.Reconcilate(e.nodeOne.Addr().Network(), e.nodeOne.Addr().String())
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	require.Equal(t, 1, e.storeTwo.Len())
	value, err := e.storeTwo.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, value)
}

func TestNodeReconcilateDeletedValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	require.NoError(t, e.storeOne.Delete(testKey))
	assert.Equal(t, 0, e.storeOne.Len())

	count, err := e.nodeTwo.Reconcilate(e.nodeOne.Addr().Network(), e.nodeOne.Addr().String())
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	require.Equal(t, 0, e.storeTwo.Len())
	value, err := e.storeTwo.Get(testKey)
	require.NoError(t, err)
	assert.Nil(t, value)
}

func TestNodeReconcilateUpdatedValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	require.NoError(t, e.storeOne.Set(testKey, testAnotherValue))
	assert.Equal(t, 1, e.storeOne.Len())

	require.NoError(t, e.storeTwo.Set(testKey, testValue))
	assert.Equal(t, 1, e.storeOne.Len())

	count, err := e.nodeTwo.Reconcilate(e.nodeOne.Addr().Network(), e.nodeOne.Addr().String())
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	require.Equal(t, 1, e.storeTwo.Len())
	value, err := e.storeTwo.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testAnotherValue, value)
}
