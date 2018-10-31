package edkvs_test

import (
	"testing"

	"github.com/simia-tech/edkvs"
)

type environment struct {
	storeOne *edkvs.Store
	tearDown func()
}

func setUpTestEnvironment(tb testing.TB) *environment {
	storeOne := edkvs.NewStore()
	return &environment{
		storeOne: storeOne,
		tearDown: func() {},
	}
}
