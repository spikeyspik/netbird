// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/netbirdio/netbird/client/internal/acl (interfaces: IFaceMapper)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	wgdevice "github.com/amnezia-vpn/amneziawg-go/device"

	"github.com/netbirdio/netbird/client/iface/device"
	"github.com/netbirdio/netbird/client/iface/wgaddr"
)

// MockIFaceMapper is a mock of IFaceMapper interface.
type MockIFaceMapper struct {
	ctrl     *gomock.Controller
	recorder *MockIFaceMapperMockRecorder
}

// MockIFaceMapperMockRecorder is the mock recorder for MockIFaceMapper.
type MockIFaceMapperMockRecorder struct {
	mock *MockIFaceMapper
}

// NewMockIFaceMapper creates a new mock instance.
func NewMockIFaceMapper(ctrl *gomock.Controller) *MockIFaceMapper {
	mock := &MockIFaceMapper{ctrl: ctrl}
	mock.recorder = &MockIFaceMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIFaceMapper) EXPECT() *MockIFaceMapperMockRecorder {
	return m.recorder
}

// Address mocks base method.
func (m *MockIFaceMapper) Address() wgaddr.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Address")
	ret0, _ := ret[0].(wgaddr.Address)
	return ret0
}

// Address indicates an expected call of Address.
func (mr *MockIFaceMapperMockRecorder) Address() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Address", reflect.TypeOf((*MockIFaceMapper)(nil).Address))
}

// IsUserspaceBind mocks base method.
func (m *MockIFaceMapper) IsUserspaceBind() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUserspaceBind")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsUserspaceBind indicates an expected call of IsUserspaceBind.
func (mr *MockIFaceMapperMockRecorder) IsUserspaceBind() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUserspaceBind", reflect.TypeOf((*MockIFaceMapper)(nil).IsUserspaceBind))
}

// Name mocks base method.
func (m *MockIFaceMapper) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockIFaceMapperMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockIFaceMapper)(nil).Name))
}

// SetFilter mocks base method.
func (m *MockIFaceMapper) SetFilter(arg0 device.PacketFilter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetFilter", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetFilter indicates an expected call of SetFilter.
func (mr *MockIFaceMapperMockRecorder) SetFilter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFilter", reflect.TypeOf((*MockIFaceMapper)(nil).SetFilter), arg0)
}

// GetDevice mocks base method.
func (m *MockIFaceMapper) GetDevice() *device.FilteredDevice {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDevice")
	ret0, _ := ret[0].(*device.FilteredDevice)
	return ret0
}

// GetDevice indicates an expected call of GetDevice.
func (mr *MockIFaceMapperMockRecorder) GetDevice() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDevice", reflect.TypeOf((*MockIFaceMapper)(nil).GetDevice))
}

// GetWGDevice mocks base method.
func (m *MockIFaceMapper) GetWGDevice() *wgdevice.Device {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWGDevice")
	ret0, _ := ret[0].(*wgdevice.Device)
	return ret0
}

// GetWGDevice indicates an expected call of GetWGDevice.
func (mr *MockIFaceMapperMockRecorder) GetWGDevice() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWGDevice", reflect.TypeOf((*MockIFaceMapper)(nil).GetWGDevice))
}
