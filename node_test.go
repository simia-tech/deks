package edkvs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeReconcilate(t *testing.T) {
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
