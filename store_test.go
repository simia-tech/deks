package edkvs_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/edkvs"
)

func TestStoreSetNewValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	assert.Equal(t, 1, e.storeOne.Len())
	require.Equal(t, 1, e.storeOne.State().Len())
	assert.Equal(t, []edkvs.Item{{0xc9, 0x59, 0xbb, 0x4d, 0xee, 0x8b, 0xc0, 0xe8, 0x29, 0xda, 0x8a, 0x69, 0xb0, 0xca, 0xab, 0x87, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}}, e.storeOne.State().Items())
}

func TestStoreSetExistingValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	assert.Equal(t, 1, e.storeOne.Len())
	require.Equal(t, 1, e.storeOne.State().Len())
	assert.Equal(t, []edkvs.Item{{0x94, 0x1e, 0x71, 0xf0, 0x80, 0x5b, 0xf5, 0x6c, 0x2c, 0xeb, 0x41, 0xde, 0x78, 0xa6, 0x12, 0x24, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}}, e.storeOne.State().Items())
}

func TestStoreSetDeletedValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	require.NoError(t, e.storeOne.Delete(testKey))

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	assert.Equal(t, 1, e.storeOne.Len())
	require.Equal(t, 1, e.storeOne.State().Len())
	assert.Equal(t, []edkvs.Item{{0x5f, 0xe3, 0x26, 0x93, 0x13, 0x2b, 0x2a, 0xf1, 0x2e, 0xfc, 0xf8, 0x52, 0x41, 0x82, 0x79, 0xc0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}}, e.storeOne.State().Items())
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
	assert.Equal(t, []edkvs.Item{{0x94, 0x1e, 0x71, 0xf0, 0x80, 0x5b, 0xf5, 0x6c, 0x2c, 0xeb, 0x41, 0xde, 0x78, 0xa6, 0x12, 0x24, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}}, e.storeOne.State().Items())
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
