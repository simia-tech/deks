package edkvs_test

import "github.com/simia-tech/edkvs"

var (
	testKey          = []byte("key")
	testValue        = []byte("value")
	testAnotherValue = []byte("another value")
	testItem         = edkvs.Item{0x01, 0x02, 0x03, 0x04}
)
