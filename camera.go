package visca

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/Zeryoshka/visca/connpool"
	"net"
	"time"
)

type Camera struct {
	conn        *connpool.Connpool
	seq         uint32
	cameraIndex byte
}

type Request struct {
	header        []byte
	payloadHeader byte
	payload       []byte
}

func (r *Request) toBytes() []byte {
	reqLen := len(r.header) + 1 + len(r.payload) + 1
	reqBytes := make([]byte, reqLen)

	copy(reqBytes[:HeaderLen], r.header)
	reqBytes[HeaderLen] = r.payloadHeader
	copy(reqBytes[HeaderLen+1:], r.payload)
	reqBytes[reqLen-1] = 0xFF

	return reqBytes
}

type requestType int

const (
	commandRequest requestType = iota
)

const HeaderLen = 8

var CommandHeaderPrefix = []byte{0x01, 0x00}

func (c *Camera) buildHeader(rtype requestType, payloadLen int) []byte {
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

func (c *Camera) getPayloadHeader() byte {
	return 0x80 | (c.cameraIndex & 0x0F)
}

func (c *Camera) buildRequestFromCommand(payload []byte) *Request {
	return &Request{
		header:        c.buildHeader(commandRequest, len(payload)),
		payloadHeader: c.getPayloadHeader(),
		payload:       payload,
	}
}

func (c *Camera) SendCommand(ctx context.Context, cmd Command) error {
	req := c.buildRequestFromCommand(cmd)
	return c.sendRequest(ctx, req)
}

func (c *Camera) sendRequest(ctx context.Context, req *Request) error {
	return c.conn.Do(ctx, func(ctx context.Context, conn net.Conn) error {
		_, err := conn.Write(req.toBytes())
		if err != nil {
			return err
		}
		return nil
	})
}

type CameraOptions struct {
	Network string
	Host    string
	Port    int
	Index   int
	Timeout time.Duration
}

// NewCamera build new Camera to manage it
func NewCamera(opts CameraOptions) (*Camera, error) {
	dialer := net.Dialer{Timeout: opts.Timeout}
	return &Camera{
		conn: connpool.NewConnpool(func(ctx context.Context) (net.Conn, error) {
			return dialer.DialContext(
				ctx,
				opts.Network,
				fmt.Sprintf("%s:%d", opts.Host, opts.Port),
			)
		}),
		cameraIndex: byte(opts.Index),
	}, nil
}
