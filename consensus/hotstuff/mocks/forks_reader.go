// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import (
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"

	model "github.com/onflow/flow-go/consensus/hotstuff/model"
)

// ForksReader is an autogenerated mock type for the ForksReader type
type ForksReader struct {
	mock.Mock
}

// FinalizedBlock provides a mock function with given fields:
func (_m *ForksReader) FinalizedBlock() *model.Block {
	ret := _m.Called()

	var r0 *model.Block
	if rf, ok := ret.Get(0).(func() *model.Block); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Block)
		}
	}

	return r0
}

// FinalizedView provides a mock function with given fields:
func (_m *ForksReader) FinalizedView() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// GetProposal provides a mock function with given fields: id
func (_m *ForksReader) GetProposal(id flow.Identifier) (*model.Proposal, bool) {
	ret := _m.Called(id)

	var r0 *model.Proposal
	if rf, ok := ret.Get(0).(func(flow.Identifier) *model.Proposal); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Proposal)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(flow.Identifier) bool); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetProposalsForView provides a mock function with given fields: view
func (_m *ForksReader) GetProposalsForView(view uint64) []*model.Proposal {
	ret := _m.Called(view)

	var r0 []*model.Proposal
	if rf, ok := ret.Get(0).(func(uint64) []*model.Proposal); ok {
		r0 = rf(view)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Proposal)
		}
	}

	return r0
}

type NewForksReaderT interface {
	mock.TestingT
	Cleanup(func())
}

// NewForksReader creates a new instance of ForksReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewForksReader(t NewForksReaderT) *ForksReader {
	mock := &ForksReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
