// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	cluster "github.com/onflow/flow-go/model/cluster"
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// CollectionMetrics is an autogenerated mock type for the CollectionMetrics type
type CollectionMetrics struct {
	mock.Mock
}

// ClusterBlockFinalized provides a mock function with given fields: block
func (_m *CollectionMetrics) ClusterBlockFinalized(block *cluster.Block) {
	_m.Called(block)
}

// ClusterBlockProposed provides a mock function with given fields: block
func (_m *CollectionMetrics) ClusterBlockProposed(block *cluster.Block) {
	_m.Called(block)
}

// TransactionIngested provides a mock function with given fields: txID
func (_m *CollectionMetrics) TransactionIngested(txID flow.Identifier) {
	_m.Called(txID)
}

type mockConstructorTestingTNewCollectionMetrics interface {
	mock.TestingT
	Cleanup(func())
}

// NewCollectionMetrics creates a new instance of CollectionMetrics. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCollectionMetrics(t mockConstructorTestingTNewCollectionMetrics) *CollectionMetrics {
	mock := &CollectionMetrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
