package rc_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/codigician/remote-code-execution/internal/mocks"
	"github.com/codigician/remote-code-execution/internal/rc"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateImageIfNotExists_ImageExists_DontCreateImage(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockContainerPort.EXPECT().
		ImageList(gomock.Any(), gomock.Any()).
		Return([]types.ImageSummary{{RepoTags: []string{"my-image"}}, {RepoTags: []string{"already-there"}}}, nil)

	client := rc.NewClient(mockContainerPort, nil)
	err := client.CreateImageIfNotExists(context.Background(), "my-image")

	assert.Nil(t, err)
}

func TestCreateImageIfNotExists_ImageDoesNotExists_CreateImage(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockIOReadCloser := newMockReadCloser(t)
	mockContainerPort.EXPECT().
		ImageList(gomock.Any(), gomock.Any()).
		Return([]types.ImageSummary{{RepoTags: []string{"already-there"}}}, nil)
	mockIOReadCloser.EXPECT().
		Read(gomock.Any()).
		Return(1, io.EOF).AnyTimes()
	mockIOReadCloser.EXPECT().
		Close().
		Return(errors.New("close error"))
	mockContainerPort.EXPECT().
		ImagePull(gomock.Any(), "my-image", gomock.Any()).
		Return(mockIOReadCloser, nil)

	client := rc.NewClient(mockContainerPort, nil)
	err := client.CreateImageIfNotExists(context.Background(), "my-image")

	assert.Nil(t, err)
}

func TestCreateImageIfNotExists_PullImageFailed_ReturnErr(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockIOReadCloser := newMockReadCloser(t)
	mockContainerPort.EXPECT().
		ImageList(gomock.Any(), gomock.Any()).
		Return([]types.ImageSummary{{RepoTags: []string{"already-there"}}}, nil)
	mockIOReadCloser.EXPECT().
		Read(gomock.Any()).
		Return(1, io.EOF).AnyTimes()
	mockContainerPort.EXPECT().
		ImagePull(gomock.Any(), "my-image", gomock.Any()).
		Return(mockIOReadCloser, errors.New("err happened when pulling"))

	client := rc.NewClient(mockContainerPort, nil)
	err := client.CreateImageIfNotExists(context.Background(), "my-image")

	assert.NotNil(t, err)
}

func TestIsImageExists_ImageListFailed_ReturnErr(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockContainerPort.EXPECT().
		ImageList(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("some err"))

	client := rc.NewClient(mockContainerPort, nil)
	exists := client.IsImageExists(context.Background(), "my-image")

	assert.False(t, exists)
}

func newMockReadCloser(t *testing.T) *mocks.MockReadCloser {
	return mocks.NewMockReadCloser(gomock.NewController(t))
}
