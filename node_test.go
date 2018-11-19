package edkvs_test

import (
	"testing"
	"time"

	"github.com/simia-tech/edkvs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeReconcilateValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	count, err := e.nodeTwo.Reconcilate(e.nodeOne.ListenURL())
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

	count, err := e.nodeTwo.Reconcilate(e.nodeOne.ListenURL())
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
	assert.Equal(t, 1, e.storeTwo.Len())

	count, err := e.nodeTwo.Reconcilate(e.nodeOne.ListenURL())
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	require.Equal(t, 1, e.storeTwo.Len())
	value, err := e.storeTwo.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testAnotherValue, value)
}

func TestNodeStreamUpdatesToAnotherNode(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	e.nodeOne.AddPeer(e.nodeTwo.ListenURL(), time.Minute, time.Minute)
	time.Sleep(100 * time.Millisecond)

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	time.Sleep(100 * time.Millisecond)

	require.Equal(t, 1, e.storeTwo.Len())
	value, err := e.storeTwo.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, value)
}

func TestNodeStreamUpdatesToTwoOtherNodes(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	storeThree := edkvs.NewStore(e.metric)
	nodeThree, err := edkvs.NewNode(storeThree, "tcp://localhost:0", e.metric)
	require.NoError(t, err)

	e.nodeOne.AddPeer(e.nodeTwo.ListenURL(), time.Minute, time.Minute)
	e.nodeOne.AddPeer(nodeThree.ListenURL(), time.Minute, time.Minute)
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

func TestNodeStreamUpdatesToAFailingNode(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	listenURL := e.nodeTwo.ListenURL()
	require.NoError(t, e.nodeTwo.Close())

	e.nodeOne.AddPeer(listenURL, time.Minute, time.Minute)
	time.Sleep(100 * time.Millisecond)

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 0, e.storeTwo.Len())
}
