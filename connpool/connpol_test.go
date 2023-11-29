package connpool

import (
	"context"
	"fmt"
	"github.com/Zeryoshka/visca/connpool/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net"
	"testing"
)

//go:generate mockgen -destination=mock/mock_conn.go -package=mock net Conn Dialer
//go:generate mockgen -destination=mock/mock_context.go -package=mock context Context

var SimpleReadedBytes = []byte{0x01, 0x02, 0x03}
var SimpleWrittenBytes = []byte{0x03, 0x05, 0x06}

func TestSimpleTask(t *testing.T) {

	tests := []int{1, 3}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d-call", tc), func(t *testing.T) {
			ctrl := gomock.NewController(t)

			newConnCounter := 0
			pool := NewConnpool(func(ctx context.Context) (net.Conn, error) {
				newConnCounter++
				conn := mock.NewMockConn(ctrl)

				conn.EXPECT().Read(gomock.Any()).DoAndReturn(func(out []byte) (int, error) {
					out = make([]byte, len(SimpleReadedBytes))
					copy(out, SimpleReadedBytes)
					return len(out), nil
				}).Times(tc).SetArg(0, SimpleReadedBytes)

				conn.EXPECT().Write(gomock.Eq(SimpleWrittenBytes)).DoAndReturn(func(in []byte) (int, error) {
					return len(SimpleWrittenBytes), nil
				}).Times(tc)

				return conn, nil
			})

			for i := 0; i < tc; i++ {
				err := pool.Do(context.TODO(), func(ctx context.Context, conn net.Conn) error {

					in := make([]byte, len(SimpleWrittenBytes))
					copy(in, SimpleWrittenBytes)
					n, err := conn.Write(in)
					assert.Nil(t, err)
					assert.Equal(t, n, len(SimpleWrittenBytes))
					assert.Equal(t, in, SimpleWrittenBytes)

					out := make([]byte, 3)
					n, err = conn.Read(out)
					assert.Nil(t, err)
					assert.Equal(t, n, len(SimpleReadedBytes))
					assert.Equal(t, out, SimpleReadedBytes)

					return nil
				})
				assert.Nil(t, err)
			}
			assert.Equal(t, newConnCounter, 1)
			ctrl.Finish()

		})
	}
}

func TestSimpleContextDone(t *testing.T) {
	ctrl := gomock.NewController(t)

	expectedNewConnCounter := 2
	newConnCounter := 0
	pool := NewConnpool(func(ctx context.Context) (net.Conn, error) {
		newConnCounter++
		conn := mock.NewMockConn(ctrl)

		conn.EXPECT().Read(gomock.Any()).DoAndReturn(func(out []byte) (int, error) {
			out = make([]byte, len(SimpleReadedBytes))
			copy(out, SimpleReadedBytes)
			return len(out), nil
		}).Times(1).SetArg(0, SimpleReadedBytes)

		conn.EXPECT().Write(gomock.Eq(SimpleWrittenBytes)).DoAndReturn(func(in []byte) (int, error) {
			return len(SimpleWrittenBytes), nil
		}).Times(1)

		conn.EXPECT().Close().Times(1)

		return conn, nil
	})

	for i := 0; i < 2; i++ {

		ctx, cancel := context.WithCancel(context.TODO())
		err := pool.Do(ctx, func(ctx context.Context, conn net.Conn) error {
			out := make([]byte, 3)
			n, err := conn.Read(out)
			assert.Nil(t, err)
			assert.Equal(t, n, len(SimpleReadedBytes))
			assert.Equal(t, out, SimpleReadedBytes)

			in := make([]byte, len(SimpleWrittenBytes))
			copy(in, SimpleWrittenBytes)
			n, err = conn.Write(in)
			assert.Nil(t, err)
			assert.Equal(t, n, len(SimpleWrittenBytes))
			assert.Equal(t, in, SimpleWrittenBytes)

			cancel()
			return nil
		})
		if assert.Error(t, err) {
			assert.Equal(t, err, context.Canceled)
		}
	}

	assert.Equal(t, newConnCounter, expectedNewConnCounter)

	ctrl.Finish()
}
