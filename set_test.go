package edkvs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/edkvs"
)

func TestSetInsert(t *testing.T) {
	set := edkvs.NewSet()
	require.NoError(t, set.Insert(testItem))
	assert.Equal(t, 1, set.Len())
	assert.Equal(t, []edkvs.Item{testItem}, set.Items())
}

func TestSetRemove(t *testing.T) {
	set := edkvs.NewSet()
	require.NoError(t, set.Insert(testItem))
	require.NoError(t, set.Remove(testItem))
	assert.Equal(t, 0, set.Len())
}
