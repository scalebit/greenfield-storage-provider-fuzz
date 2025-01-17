// Code generated by MockGen. DO NOT EDIT.
// Source: ./piece_store_client.go

// Package client is a generated GoMock package.
package client

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockPieceStoreAPI is a mock of PieceStoreAPI interface.
type MockPieceStoreAPI struct {
	ctrl     *gomock.Controller
	recorder *MockPieceStoreAPIMockRecorder
}

// MockPieceStoreAPIMockRecorder is the mock recorder for MockPieceStoreAPI.
type MockPieceStoreAPIMockRecorder struct {
	mock *MockPieceStoreAPI
}

// NewMockPieceStoreAPI creates a new mock instance.
func NewMockPieceStoreAPI(ctrl *gomock.Controller) *MockPieceStoreAPI {
	mock := &MockPieceStoreAPI{ctrl: ctrl}
	mock.recorder = &MockPieceStoreAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPieceStoreAPI) EXPECT() *MockPieceStoreAPIMockRecorder {
	return m.recorder
}

// DeletePiece mocks base method.
func (m *MockPieceStoreAPI) DeletePiece(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePiece", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePiece indicates an expected call of DeletePiece.
func (mr *MockPieceStoreAPIMockRecorder) DeletePiece(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePiece", reflect.TypeOf((*MockPieceStoreAPI)(nil).DeletePiece), ctx, key)
}

// GetPiece mocks base method.
func (m *MockPieceStoreAPI) GetPiece(ctx context.Context, key string, offset, limit int64) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPiece", ctx, key, offset, limit)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPiece indicates an expected call of GetPiece.
func (mr *MockPieceStoreAPIMockRecorder) GetPiece(ctx, key, offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPiece", reflect.TypeOf((*MockPieceStoreAPI)(nil).GetPiece), ctx, key, offset, limit)
}

// PutPiece mocks base method.
func (m *MockPieceStoreAPI) PutPiece(key string, value []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutPiece", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutPiece indicates an expected call of PutPiece.
func (mr *MockPieceStoreAPIMockRecorder) PutPiece(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutPiece", reflect.TypeOf((*MockPieceStoreAPI)(nil).PutPiece), key, value)
}
