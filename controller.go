package visca

import (
	"context"
	"encoding/binary"
	"github.com/Zeryoshka/visca/connpool"
	"net"
)

type Controller struct {
	conns   map[string]*connpool.Connpool
	seq     uint32
	cameras map[viscaOverIpAddr]*Camera
}

type request struct {
	header        []byte
	payloadHeader byte
	payload       []byte
}

func (r *request) toBytes() []byte {
	reqLen := len(r.header) + 1 + len(r.payload) + 1
	reqBytes := make([]byte, reqLen)

	copy(reqBytes[:HeaderLen], r.header)
	reqBytes[HeaderLen] = r.payloadHeader
	copy(reqBytes[HeaderLen+1:], r.payload)
	reqBytes[reqLen-1] = 0xFF

	return reqBytes
}

func NewController() (*Controller, error) {
	return &Controller{
		conns:   make(map[string]*connpool.Connpool),
		cameras: make(map[viscaOverIpAddr]*Camera),
		seq:     0,
	}, nil
}

func (c *Controller) Close() error {
	for _, conn := range c.conns {
		err := conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) buildHeader(rtype requestType, payloadLen int) []byte {
	header := make([]byte, HeaderLen)
	// bytes 0, 1 - request type bytes from spec
	switch rtype {
	case commandRequest:
		copy(header[:2], CommandHeaderPrefix)
	}
	// bytes 2, 3 - len of payload (2 = 0x0 cause max len 16)
	header[2] = 0x00
	header[3] = byte(payloadLen) + 2 // basepayload(3-14) + payloadHeader(1) + 0xFF

	// bytes 4...8 - seq
	binary.LittleEndian.PutUint32(header[4:], c.seq)
	c.seq++ // with overflow (needs for devices)
	return header
}

func (c *Controller) getPayloadHeader(index byte) byte {
	return 0x80 | (index & 0x0F)
}

func (c *Controller) buildRequestFromCommand(index byte, payload []byte) *request {
	return &request{
		header:        c.buildHeader(commandRequest, len(payload)),
		payloadHeader: c.getPayloadHeader(index),
		payload:       payload,
	}
}

func (c *Controller) sendCommand(ctx context.Context, viscaAddr viscaOverIpAddr, cmd Command) error {
	req := c.buildRequestFromCommand(viscaAddr.index, cmd)
	return c.sendRequest(ctx, viscaAddr.addr, req)
}

func (c *Controller) sendRequest(ctx context.Context, addr string, req *request) error {
	return c.conns[addr].Do(ctx, func(ctx context.Context, conn net.Conn) error {
		_, err := conn.Write(req.toBytes())
		return err
	})
}
