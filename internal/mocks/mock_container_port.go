// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/codigician/remote-code-execution/internal/rc (interfaces: ContainerPort)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	io "io"
	reflect "reflect"

	types "github.com/docker/docker/api/types"
	container "github.com/docker/docker/api/types/container"
	network "github.com/docker/docker/api/types/network"
	gomock "github.com/golang/mock/gomock"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// MockContainerPort is a mock of ContainerPort interface.
type MockContainerPort struct {
	ctrl     *gomock.Controller
	recorder *MockContainerPortMockRecorder
}

// MockContainerPortMockRecorder is the mock recorder for MockContainerPort.
type MockContainerPortMockRecorder struct {
	mock *MockContainerPort
}

// NewMockContainerPort creates a new mock instance.
func NewMockContainerPort(ctrl *gomock.Controller) *MockContainerPort {
	mock := &MockContainerPort{ctrl: ctrl}
	mock.recorder = &MockContainerPortMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContainerPort) EXPECT() *MockContainerPortMockRecorder {
	return m.recorder
}

// ContainerCreate mocks base method.
func (m *MockContainerPort) ContainerCreate(arg0 context.Context, arg1 *container.Config, arg2 *container.HostConfig, arg3 *network.NetworkingConfig, arg4 *v1.Platform, arg5 string) (container.ContainerCreateCreatedBody, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerCreate", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(container.ContainerCreateCreatedBody)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainerCreate indicates an expected call of ContainerCreate.
func (mr *MockContainerPortMockRecorder) ContainerCreate(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerCreate", reflect.TypeOf((*MockContainerPort)(nil).ContainerCreate), arg0, arg1, arg2, arg3, arg4, arg5)
}

// ContainerExecAttach mocks base method.
func (m *MockContainerPort) ContainerExecAttach(arg0 context.Context, arg1 string, arg2 types.ExecStartCheck) (types.HijackedResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerExecAttach", arg0, arg1, arg2)
	ret0, _ := ret[0].(types.HijackedResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainerExecAttach indicates an expected call of ContainerExecAttach.
func (mr *MockContainerPortMockRecorder) ContainerExecAttach(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerExecAttach", reflect.TypeOf((*MockContainerPort)(nil).ContainerExecAttach), arg0, arg1, arg2)
}

// ContainerExecCreate mocks base method.
func (m *MockContainerPort) ContainerExecCreate(arg0 context.Context, arg1 string, arg2 types.ExecConfig) (types.IDResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerExecCreate", arg0, arg1, arg2)
	ret0, _ := ret[0].(types.IDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainerExecCreate indicates an expected call of ContainerExecCreate.
func (mr *MockContainerPortMockRecorder) ContainerExecCreate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerExecCreate", reflect.TypeOf((*MockContainerPort)(nil).ContainerExecCreate), arg0, arg1, arg2)
}

// ContainerList mocks base method.
func (m *MockContainerPort) ContainerList(arg0 context.Context, arg1 types.ContainerListOptions) ([]types.Container, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerList", arg0, arg1)
	ret0, _ := ret[0].([]types.Container)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainerList indicates an expected call of ContainerList.
func (mr *MockContainerPortMockRecorder) ContainerList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerList", reflect.TypeOf((*MockContainerPort)(nil).ContainerList), arg0, arg1)
}

// ContainerStart mocks base method.
func (m *MockContainerPort) ContainerStart(arg0 context.Context, arg1 string, arg2 types.ContainerStartOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerStart", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ContainerStart indicates an expected call of ContainerStart.
func (mr *MockContainerPortMockRecorder) ContainerStart(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerStart", reflect.TypeOf((*MockContainerPort)(nil).ContainerStart), arg0, arg1, arg2)
}

// ImageList mocks base method.
func (m *MockContainerPort) ImageList(arg0 context.Context, arg1 types.ImageListOptions) ([]types.ImageSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageList", arg0, arg1)
	ret0, _ := ret[0].([]types.ImageSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImageList indicates an expected call of ImageList.
func (mr *MockContainerPortMockRecorder) ImageList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageList", reflect.TypeOf((*MockContainerPort)(nil).ImageList), arg0, arg1)
}

// ImagePull mocks base method.
func (m *MockContainerPort) ImagePull(arg0 context.Context, arg1 string, arg2 types.ImagePullOptions) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImagePull", arg0, arg1, arg2)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImagePull indicates an expected call of ImagePull.
func (mr *MockContainerPortMockRecorder) ImagePull(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImagePull", reflect.TypeOf((*MockContainerPort)(nil).ImagePull), arg0, arg1, arg2)
}
