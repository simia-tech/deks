package kea_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/simia-tech/kea"
)

type environment struct {
	metric    kea.Metric
	storeOne  *kea.Store
	serverOne *kea.Server
	storeTwo  *kea.Store
	serverTwo *kea.Server
	tearDown  func()
}

func setUpTestEnvironment(tb testing.TB) *environment {
	m := kea.NewMetricMock()

	storeOne := kea.NewStore(m)
	serverOne, err := kea.NewServer(storeOne, "tcp://localhost:0", m)
	require.NoError(tb, err)

	storeTwo := kea.NewStore(m)
	serverTwo, err := kea.NewServer(storeTwo, "tcp://localhost:0", m)
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
