package visca

import (
	"context"
	"github.com/Zeryoshka/visca/connpool"
	"net"
	"time"
)

const ViscaPort = 52381

type Camera struct {
	viscaAddr  viscaOverIpAddr
	controller *Controller
}

func (c *Controller) AddCamera(address string, index int, timeout time.Duration) (*Camera, error) {
	if _, ok := c.conns[address]; !ok {
		dialer := net.Dialer{Timeout: timeout}

		c.conns[address] = connpool.NewConnpool(
			func(ctx context.Context) (net.Conn, error) {
				return dialer.DialContext(ctx, "udp", address)
			},
		)
	}

	if (index <= 0) || (index >= 8) {
		return nil, IncorrectDeviceIndexErr
	}
	viscaAddr := viscaOverIpAddr{address, byte(index)}
	if _, busy := c.cameras[viscaAddr]; busy {
		return nil, IncorrectDeviceIndexErr
	}

	camera := &Camera{viscaAddr, c}
	c.cameras[viscaAddr] = camera

	return camera, nil
}

func (c *Controller) RemoveCamera(viscaAddr viscaOverIpAddr) error {
	camera, exists := c.cameras[viscaAddr]
	if !exists {
		return CameraNotFoundErr
	}
	if err := camera.controller.conns[viscaAddr.addr].Close(); err != nil {
		return err
	}
	delete(camera.controller.conns, viscaAddr.addr)
	delete(c.cameras, viscaAddr)
	return nil
}

func (c *Camera) SendCommand(ctx context.Context, cmd Command) error {
	return c.controller.sendCommand(ctx, c.viscaAddr, cmd)
}

type viscaOverIpAddr struct {
	addr  string
	index byte
}
