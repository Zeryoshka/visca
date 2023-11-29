package visca

import (
	"context"
	"github.com/Zeryoshka/visca/connpool"
	"github.com/Zeryoshka/visca/connpool/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net"
	"testing"
)

func newCamera(newConn connpool.NewConnFunc, opts CameraOptions) *Camera {
	return &Camera{
		conn:        connpool.NewConnpool(newConn),
		cameraIndex: byte(opts.Index),
	}
}

func TestSendCommnd(t *testing.T) {
	ctrl := gomock.NewController(t)

	header := []byte{
		0x01, 0x00, 0x00, 0x06,
		0x00, 0x00, 0x00, 0x00,
	}
	payload := []byte{
		0x81,
		0x01, 0x04, 0x35, 0x00,
		0xFF,
	}
	reqLine := make([]byte, len(header)+len(payload))
	copy(reqLine[:len(header)], header)
	copy(reqLine[len(header):], payload)

	camera := newCamera(func(ctx context.Context) (net.Conn, error) {
		conn := mock.NewMockConn(ctrl)

		conn.EXPECT().Write(gomock.Eq(reqLine)).Times(1).Return(len(reqLine), nil)

		return conn, nil
	}, CameraOptions{Index: 1})

	err := camera.SendCommand(context.TODO(), CamWbCommand(WbAuto1))
	assert.Nil(t, err)

	ctrl.Finish()
}
