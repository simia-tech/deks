package edkvs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/edkvs"
)

func TestConnSetAndGet(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	conn, err := edkvs.Dial(e.nodeOne.Addr().Network(), e.nodeOne.Addr().String())
	require.NoError(t, err)

	require.NoError(t, conn.Set(testKey, testValue))

	value, err := conn.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, value)
}