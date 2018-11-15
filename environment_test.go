package edkvs_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/simia-tech/edkvs"
)

type environment struct {
	metric   edkvs.Metric
	storeOne *edkvs.Store
	nodeOne  *edkvs.Node
	storeTwo *edkvs.Store
	nodeTwo  *edkvs.Node
	tearDown func()
}

func setUpTestEnvironment(tb testing.TB) *environment {
	m := edkvs.NewMetricMock()

	storeOne := edkvs.NewStore(m)
	nodeOne, err := edkvs.NewNode(storeOne, "tcp://localhost:0")
	require.NoError(tb, err)

	storeTwo := edkvs.NewStore(m)
	nodeTwo, err := edkvs.NewNode(storeTwo, "tcp://localhost:0")
	require.NoError(tb, err)

	return &environment{
		metric:   m,
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
