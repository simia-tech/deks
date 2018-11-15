package edkvs_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/simia-tech/edkvs"
)

type environment struct {
	storeOne *edkvs.Store
	nodeOne  *edkvs.Node
	storeTwo *edkvs.Store
	nodeTwo  *edkvs.Node
	tearDown func()
}

func setUpTestEnvironment(tb testing.TB) *environment {
	storeOne := edkvs.NewStore()
	nodeOne, err := edkvs.NewNode(storeOne, "tcp://localhost:0")
	require.NoError(tb, err)

	storeTwo := edkvs.NewStore()
	nodeTwo, err := edkvs.NewNode(storeTwo, "tcp://localhost:0")
	require.NoError(tb, err)

	return &environment{
		storeOne: storeOne,
		nodeOne:  nodeOne,
		storeTwo: storeTwo,
		nodeTwo:  nodeTwo,
		tearDown: func() {
			require.NoError(tb, nodeOne.Close())
			require.NoError(tb, nodeTwo.Close())
		},
	}
}
