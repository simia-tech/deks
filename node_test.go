package edkvs_test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeReconcilate(t *testing.T) {
	t.SkipNow()
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	require.NoError(t, e.nodeOne.Insert(300))
	require.NoError(t, e.nodeTwo.Insert(100))
	require.NoError(t, e.nodeTwo.Insert(200))

	n, err := e.nodeOne.Reconcilate(e.nodeTwo.Addr())
	require.NoError(t, err)
	assert.Equal(t, 1, n)

	time.Sleep(200 * time.Millisecond)

	log.Printf("node one %v", e.nodeOne.Elements())
	log.Printf("node two %v", e.nodeTwo.Elements())

	require.Equal(t, 1, e.storeTwo.Len())
	value, err := e.storeTwo.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, value)
}
