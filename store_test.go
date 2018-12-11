package deks_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/deks"
)

func TestStoreSetNewValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	assert.Equal(t, 1, e.storeOne.Len())
	require.Equal(t, 1, e.storeOne.State().Len())
	assert.Equal(t, []deks.Item{{0xa6, 0x2f, 0x22, 0x25, 0xbf, 0x70, 0xbf, 0xac, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}}, e.storeOne.State().Items())
}

func TestStoreSetExistingValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	assert.Equal(t, 1, e.storeOne.Len())
	require.Equal(t, 1, e.storeOne.State().Len())
	assert.Equal(t, []deks.Item{{0xa6, 0x2f, 0x22, 0x25, 0xbf, 0x70, 0xbf, 0xac, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}}, e.storeOne.State().Items())
}

func TestStoreSetDeletedValue(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	require.NoError(t, e.storeOne.Delete(testKey))

	require.NoError(t, e.storeOne.Set(testKey, testValue))

	assert.Equal(t, 1, e.storeOne.Len())
	require.Equal(t, 1, e.storeOne.State().Len())
	assert.Equal(t, []deks.Item{{0xa6, 0x2f, 0x22, 0x25, 0xbf, 0x70, 0xbf, 0xac, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2}}, e.storeOne.State().Items())
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
	assert.Equal(t, []deks.Item{{0xa6, 0x2f, 0x22, 0x25, 0xbf, 0x70, 0xbf, 0xac, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}}, e.storeOne.State().Items())
}

func TestStoreEach(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	require.Equal(t, 1, e.storeOne.Len())

	err := e.storeOne.Each(func(key, value []byte) error {
		assert.Equal(t, testKey, key)
		assert.Equal(t, testValue, value)
		return nil
	})
	require.NoError(t, err)
}

func TestStoreTidy(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	require.NoError(t, e.storeOne.Set(testKey, testValue))
	require.NoError(t, e.storeOne.Delete(testKey))
	require.Equal(t, 0, e.storeOne.Len())

	require.Equal(t, 1, e.storeOne.DeletedLen())

	require.NoError(t, e.storeOne.Tidy())

	assert.Equal(t, 0, e.storeOne.DeletedLen())
	assert.Equal(t, 0, e.storeOne.Len())
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
