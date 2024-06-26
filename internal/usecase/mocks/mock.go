// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	big "math/big"
	reflect "reflect"

	entity "github.com/egor-denisov/biggest-change/internal/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockStatsOfChanging is a mock of StatsOfChanging interface.
type MockStatsOfChanging struct {
	ctrl     *gomock.Controller
	recorder *MockStatsOfChangingMockRecorder
}

// MockStatsOfChangingMockRecorder is the mock recorder for MockStatsOfChanging.
type MockStatsOfChangingMockRecorder struct {
	mock *MockStatsOfChanging
}

// NewMockStatsOfChanging creates a new mock instance.
func NewMockStatsOfChanging(ctrl *gomock.Controller) *MockStatsOfChanging {
	mock := &MockStatsOfChanging{ctrl: ctrl}
	mock.recorder = &MockStatsOfChangingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStatsOfChanging) EXPECT() *MockStatsOfChangingMockRecorder {
	return m.recorder
}

// GetAddressWithBiggestChange mocks base method.
func (m *MockStatsOfChanging) GetAddressWithBiggestChange(ctx context.Context, countOfLastBlocks uint) (*entity.BiggestChange, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddressWithBiggestChange", ctx, countOfLastBlocks)
	ret0, _ := ret[0].(*entity.BiggestChange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAddressWithBiggestChange indicates an expected call of GetAddressWithBiggestChange.
func (mr *MockStatsOfChangingMockRecorder) GetAddressWithBiggestChange(ctx, countOfLastBlocks interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddressWithBiggestChange", reflect.TypeOf((*MockStatsOfChanging)(nil).GetAddressWithBiggestChange), ctx, countOfLastBlocks)
}

// MockStatsOfChangingWebAPI is a mock of StatsOfChangingWebAPI interface.
type MockStatsOfChangingWebAPI struct {
	ctrl     *gomock.Controller
	recorder *MockStatsOfChangingWebAPIMockRecorder
}

// MockStatsOfChangingWebAPIMockRecorder is the mock recorder for MockStatsOfChangingWebAPI.
type MockStatsOfChangingWebAPIMockRecorder struct {
	mock *MockStatsOfChangingWebAPI
}

// NewMockStatsOfChangingWebAPI creates a new mock instance.
func NewMockStatsOfChangingWebAPI(ctrl *gomock.Controller) *MockStatsOfChangingWebAPI {
	mock := &MockStatsOfChangingWebAPI{ctrl: ctrl}
	mock.recorder = &MockStatsOfChangingWebAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStatsOfChangingWebAPI) EXPECT() *MockStatsOfChangingWebAPIMockRecorder {
	return m.recorder
}

// GetCurrentBlockNumber mocks base method.
func (m *MockStatsOfChangingWebAPI) GetCurrentBlockNumber(ctx context.Context) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentBlockNumber", ctx)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentBlockNumber indicates an expected call of GetCurrentBlockNumber.
func (mr *MockStatsOfChangingWebAPIMockRecorder) GetCurrentBlockNumber(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentBlockNumber", reflect.TypeOf((*MockStatsOfChangingWebAPI)(nil).GetCurrentBlockNumber), ctx)
}

// GetTransactionsByBlockNumber mocks base method.
func (m *MockStatsOfChangingWebAPI) GetTransactionsByBlockNumber(ctx context.Context, blockNumber *big.Int) ([]*entity.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionsByBlockNumber", ctx, blockNumber)
	ret0, _ := ret[0].([]*entity.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionsByBlockNumber indicates an expected call of GetTransactionsByBlockNumber.
func (mr *MockStatsOfChangingWebAPIMockRecorder) GetTransactionsByBlockNumber(ctx, blockNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionsByBlockNumber", reflect.TypeOf((*MockStatsOfChangingWebAPI)(nil).GetTransactionsByBlockNumber), ctx, blockNumber)
}
