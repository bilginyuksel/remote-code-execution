package rc_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/codigician/remote-code-execution/internal/mocks"
	"github.com/codigician/remote-code-execution/internal/rc"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const _sampleDuration = time.Second

func TestNewExecRes_NilOutBuffer_UnsuccessfullAndErrBuffer(t *testing.T) {
	res := rc.NewExecRes("", "some error", _sampleDuration)

	assert.Equal(t, "some error", res.Buffer())
}

func TestNewExecRes_HaveOutBuffer_SuccessfullAndOutBuffer(t *testing.T) {
	res := rc.NewExecRes("some output", "", _sampleDuration)

	assert.Equal(t, "some output", res.Buffer())
}

func TestExec_CreateExecFailed_ReturnErr(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockContainerPort.EXPECT().
		ContainerExecCreate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(types.IDResponse{}, errors.New("err"))

	c := rc.NewClient(mockContainerPort, nil)
	_, err := c.Exec(context.TODO(), "containerID", "dir/", []string{})

	assert.NotNil(t, err)
}

func TestExec_AttachExecFailed_ReturnErr(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockContainerPort.EXPECT().
		ContainerExecCreate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(types.IDResponse{ID: "exec-id"}, nil)
	mockContainerPort.EXPECT().
		ContainerExecAttach(gomock.Any(), "exec-id", gomock.Any()).
		Return(types.HijackedResponse{}, errors.New("exec attaching err"))

	c := rc.NewClient(mockContainerPort, nil)
	_, err := c.Exec(context.TODO(), "containerID", "/dir", []string{})

	assert.NotNil(t, err)
}

func TestExec_ExecCreatedAndAttached_ReturnOutBuffer(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockConn := newMockConn(t)
	mockContainerPort.EXPECT().
		ContainerExecCreate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(types.IDResponse{ID: "exec-id"}, nil)
	mockConn.EXPECT().Read(gomock.Any()).Return(1, io.EOF)
	mockContainerPort.EXPECT().
		ContainerExecAttach(gomock.Any(), "exec-id", gomock.Any()).
		Return(types.HijackedResponse{Conn: mockConn}, nil)

	c := rc.NewClient(mockContainerPort, nil)
	_, err := c.Exec(context.TODO(), "containerID", "/dir", []string{})

	assert.Nil(t, err)
}

func newMockConn(t *testing.T) *mocks.MockConn {
	return mocks.NewMockConn(gomock.NewController(t))
}
