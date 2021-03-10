// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/psewda/typing/pkg/signin/userinfo (interfaces: Userinfo)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	userinfo "github.com/psewda/typing/pkg/signin/userinfo"
	reflect "reflect"
)

// MockUserinfo is a mock of Userinfo interface
type MockUserinfo struct {
	ctrl     *gomock.Controller
	recorder *MockUserinfoMockRecorder
}

// MockUserinfoMockRecorder is the mock recorder for MockUserinfo
type MockUserinfoMockRecorder struct {
	mock *MockUserinfo
}

// NewMockUserinfo creates a new mock instance
func NewMockUserinfo(ctrl *gomock.Controller) *MockUserinfo {
	mock := &MockUserinfo{ctrl: ctrl}
	mock.recorder = &MockUserinfoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserinfo) EXPECT() *MockUserinfoMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockUserinfo) Get() (*userinfo.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(*userinfo.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockUserinfoMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserinfo)(nil).Get))
}