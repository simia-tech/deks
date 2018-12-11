package deks_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/simia-tech/deks"
)

type environment struct {
	metric    deks.Metric
	storeOne  *deks.Store
	serverOne *deks.Server
	storeTwo  *deks.Store
	serverTwo *deks.Server
	tearDown  func()
}

func setUpTestEnvironment(tb testing.TB) *environment {
	m := deks.NewMetricMock()

	storeOne := deks.NewStore(m)
	serverOne, err := deks.NewServer(storeOne, "tcp://localhost:0", m)
	require.NoError(tb, err)

	storeTwo := deks.NewStore(m)
	serverTwo, err := deks.NewServer(storeTwo, "tcp://localhost:0", m)
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
