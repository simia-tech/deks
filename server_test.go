package kea_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/kea"
)

func TestServerReconcilateValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	count, err := e.serverTwo.Reconcilate(e.serverOne.ListenURL())
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	require.Equal(t, 1, e.storeTwo.Len())
	value, err := e.storeTwo.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, value)
}

func TestServerReconcilateDeletedValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	require.NoError(t, e.storeOne.Delete(testKey))
	assert.Equal(t, 0, e.storeOne.Len())

	count, err := e.serverTwo.Reconcilate(e.serverOne.ListenURL())
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	require.Equal(t, 0, e.storeTwo.Len())
	value, err := e.storeTwo.Get(testKey)
	require.NoError(t, err)
	assert.Nil(t, value)
}

func TestServerReconcilateUpdatedValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	require.NoError(t, e.storeOne.Set(testKey, testAnotherValue))
	assert.Equal(t, 1, e.storeOne.Len())

	require.NoError(t, e.storeTwo.Set(testKey, testValue))
	assert.Equal(t, 1, e.storeTwo.Len())

	count, err := e.serverTwo.Reconcilate(e.serverOne.ListenURL())
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	require.Equal(t, 1, e.storeTwo.Len())
	value, err := e.storeTwo.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testAnotherValue, value)
}

func TestServerStreamUpdatesToAnotherNode(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	e.serverOne.AddPeer(e.serverTwo.ListenURL(), time.Minute, time.Minute)
	time.Sleep(100 * time.Millisecond)

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	time.Sleep(100 * time.Millisecond)

	require.Equal(t, 1, e.storeTwo.Len())
	value, err := e.storeTwo.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, value)
}

func TestServerStreamUpdatesToTwoOtherNodes(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	storeThree := kea.NewStore(e.metric)
	serverThree, err := kea.NewServer(storeThree, "tcp://localhost:0", e.metric)
	require.NoError(t, err)

	e.serverOne.AddPeer(e.serverTwo.ListenURL(), time.Minute, time.Minute)
	e.serverOne.AddPeer(serverThree.ListenURL(), time.Minute, time.Minute)
	time.Sleep(100 * time.Millisecond)

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	time.Sleep(100 * time.Millisecond)

	require.Equal(t, 1, e.storeTwo.Len())
	value, err := e.storeTwo.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, value)

	require.Equal(t, 1, storeThree.Len())
	value, err = storeThree.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, value)
}

func TestServerStreamUpdatesToAFailingNode(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	listenURL := e.serverTwo.ListenURL()
	require.NoError(t, e.serverTwo.Close())

	e.serverOne.AddPeer(listenURL, time.Minute, time.Minute)
	time.Sleep(100 * time.Millisecond)

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 0, e.storeTwo.Len())
}
