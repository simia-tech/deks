package edkvs_test

import (
	"testing"

	"github.com/simia-tech/edkvs"
)

type environment struct {
	storeOne *edkvs.Store
	storeTwo *edkvs.Store
	tearDown func()
}

func setUpTestEnvironment(tb testing.TB) *environment {
	storeOne := edkvs.NewStore()
	storeTwo := edkvs.NewStore()
	return &environment{
		storeOne: storeOne,
		storeTwo: storeTwo,
		tearDown: func() {},
	}
}
