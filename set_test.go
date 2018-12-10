package kea_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/kea"
)

func TestSetInsert(t *testing.T) {
	set := kea.NewSet()
	require.NoError(t, set.Insert(testItem))
	assert.Equal(t, 1, set.Len())
	assert.Equal(t, []kea.Item{testItem}, set.Items())
}

func TestSetRemove(t *testing.T) {
	set := kea.NewSet()
	require.NoError(t, set.Insert(testItem))
	require.NoError(t, set.Remove(testItem))
	assert.Equal(t, 0, set.Len())
}
