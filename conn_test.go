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

	conn, err := edkvs.Dial(e.nodeOne.ListenURL())
	require.NoError(t, err)

	require.NoError(t, conn.Set(testKey, testValue))

	value, err := conn.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, value)
}

func TestConnDelete(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	conn, err := edkvs.Dial(e.nodeOne.ListenURL())
	require.NoError(t, err)

	require.NoError(t, conn.Set(testKey, testValue))

	require.NoError(t, conn.Delete(testKey))

	keys, err := conn.Keys()
	require.NoError(t, err)
	assert.Len(t, keys, 0)
}

func TestConnKeys(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	conn, err := edkvs.Dial(e.nodeOne.ListenURL())
	require.NoError(t, err)

	require.NoError(t, conn.Set(testKey, testValue))

	keys, err := conn.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	assert.Equal(t, testKey, keys[0])
}

func TestConnPing(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	conn, err := edkvs.Dial(e.nodeOne.ListenURL())
	require.NoError(t, err)

	require.NoError(t, conn.Ping())
}
