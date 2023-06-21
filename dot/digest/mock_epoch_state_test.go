// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ChainSafe/gossamer/dot/digest (interfaces: EpochState)

// Package digest is a generated GoMock package.
package digest

import (
	reflect "reflect"

	types "github.com/ChainSafe/gossamer/dot/types"
	scale "github.com/ChainSafe/gossamer/pkg/scale"
	gomock "github.com/golang/mock/gomock"
)

// MockEpochState is a mock of EpochState interface.
type MockEpochState struct {
	ctrl     *gomock.Controller
	recorder *MockEpochStateMockRecorder
}

// MockEpochStateMockRecorder is the mock recorder for MockEpochState.
type MockEpochStateMockRecorder struct {
	mock *MockEpochState
}

// NewMockEpochState creates a new mock instance.
func NewMockEpochState(ctrl *gomock.Controller) *MockEpochState {
	mock := &MockEpochState{ctrl: ctrl}
	mock.recorder = &MockEpochStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEpochState) EXPECT() *MockEpochStateMockRecorder {
	return m.recorder
}

// FinalizeBABENextConfigData mocks base method.
func (m *MockEpochState) FinalizeBABENextConfigData(arg0 *types.Header) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinalizeBABENextConfigData", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// FinalizeBABENextConfigData indicates an expected call of FinalizeBABENextConfigData.
func (mr *MockEpochStateMockRecorder) FinalizeBABENextConfigData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinalizeBABENextConfigData", reflect.TypeOf((*MockEpochState)(nil).FinalizeBABENextConfigData), arg0)
}

// FinalizeBABENextEpochData mocks base method.
func (m *MockEpochState) FinalizeBABENextEpochData(arg0 *types.Header) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinalizeBABENextEpochData", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// FinalizeBABENextEpochData indicates an expected call of FinalizeBABENextEpochData.
func (mr *MockEpochStateMockRecorder) FinalizeBABENextEpochData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinalizeBABENextEpochData", reflect.TypeOf((*MockEpochState)(nil).FinalizeBABENextEpochData), arg0)
}

// GetEpochForBlock mocks base method.
func (m *MockEpochState) GetEpochForBlock(arg0 *types.Header) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEpochForBlock", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEpochForBlock indicates an expected call of GetEpochForBlock.
func (mr *MockEpochStateMockRecorder) GetEpochForBlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEpochForBlock", reflect.TypeOf((*MockEpochState)(nil).GetEpochForBlock), arg0)
}

// HandleBABEDigest mocks base method.
func (m *MockEpochState) HandleBABEDigest(arg0 *types.Header, arg1 scale.VaryingDataType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleBABEDigest", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleBABEDigest indicates an expected call of HandleBABEDigest.
func (mr *MockEpochStateMockRecorder) HandleBABEDigest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleBABEDigest", reflect.TypeOf((*MockEpochState)(nil).HandleBABEDigest), arg0, arg1)
}
