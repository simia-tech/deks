package deks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/deks"
)

func TestSetInsert(t *testing.T) {
	set := deks.NewSet()
	require.NoError(t, set.Insert(testItem))
	assert.Equal(t, 1, set.Len())
	assert.Equal(t, []deks.Item{testItem}, set.Items())
}

func TestSetRemove(t *testing.T) {
	set := deks.NewSet()
	require.NoError(t, set.Insert(testItem))
	require.NoError(t, set.Remove(testItem))
	assert.Equal(t, 0, set.Len())
}
