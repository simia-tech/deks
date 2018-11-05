package edkvs_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreSetNewValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	assert.Equal(t, 1, e.storeOne.Len())
	require.Equal(t, 1, e.storeOne.State().Len())
	assert.Equal(t, [][]byte{{0x6b, 0x65, 0x79}}, e.storeOne.State().Items())
}

func TestStoreSetExistingValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	assert.Equal(t, 1, e.storeOne.Len())
	require.Equal(t, 1, e.storeOne.State().Len())
	assert.Equal(t, [][]byte{{0x6b, 0x65, 0x79, 0x01}}, e.storeOne.State().Items())
}

func TestStoreSetDeletedValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	require.NoError(t, e.storeOne.Delete(testKey))

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	assert.Equal(t, 1, e.storeOne.Len())
	require.Equal(t, 1, e.storeOne.State().Len())
	assert.Equal(t, [][]byte{{0x6b, 0x65, 0x79, 0x02}}, e.storeOne.State().Items())
}

func TestStoreGet(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	value, err := e.storeOne.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, value)
}

func TestStoreDelete(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	require.NoError(t, e.storeOne.Delete(testKey))

	assert.Equal(t, e.storeOne.Len(), 0)
	require.Equal(t, 1, e.storeOne.State().Len())
	assert.Equal(t, [][]byte{{0x6b, 0x65, 0x79, 0x01}}, e.storeOne.State().Items())
}

func TestStoreConcurrentAccess(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	wg := sync.WaitGroup{}
	for index := 0; index < 10; index++ {
		wg.Add(1)
		go func() {
			time.Sleep(10 * time.Millisecond)
			e.storeOne.Set(testKey, testValue)
			time.Sleep(10 * time.Millisecond)
			e.storeOne.Get(testKey)
			time.Sleep(10 * time.Millisecond)
			e.storeOne.Delete(testKey)
			time.Sleep(10 * time.Millisecond)
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkStoreSet(b *testing.B) {
	b.ReportAllocs()

	e := setUpTestEnvironment(b)
	defer e.tearDown()

	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		e.storeOne.Set(testKey, testValue)
	}
}

func BenchmarkStoreGet(b *testing.B) {
	b.ReportAllocs()

	e := setUpTestEnvironment(b)
	defer e.tearDown()
	e.storeOne.Set(testKey, testValue)

	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		e.storeOne.Get(testKey)
	}
}
