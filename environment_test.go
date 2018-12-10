package edkvs_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/simia-tech/edkvs"
)

type environment struct {
	metric    edkvs.Metric
	storeOne  *edkvs.Store
	serverOne *edkvs.Server
	storeTwo  *edkvs.Store
	serverTwo *edkvs.Server
	tearDown  func()
}

func setUpTestEnvironment(tb testing.TB) *environment {
	m := edkvs.NewMetricMock()

	storeOne := edkvs.NewStore(m)
	serverOne, err := edkvs.NewServer(storeOne, "tcp://localhost:0", m)
	require.NoError(tb, err)

	storeTwo := edkvs.NewStore(m)
	serverTwo, err := edkvs.NewServer(storeTwo, "tcp://localhost:0", m)
	require.NoError(tb, err)

	return &environment{
		metric:    m,
		storeOne:  storeOne,
		serverOne: serverOne,
		storeTwo:  storeTwo,
		serverTwo: serverTwo,
		tearDown: func() {
			require.NoError(tb, serverOne.Close())
			require.NoError(tb, serverTwo.Close())
		},
	}
}
