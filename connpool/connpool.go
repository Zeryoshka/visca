package connpool

import (
	"context"
	"net"
	"sync"
)

type NewConnFunc func(ctx context.Context) (net.Conn, error)

type Task func(ctx context.Context, conn net.Conn) error

type request struct {
	task      Task
	ctx       context.Context
	resp      chan error
	errToResp chan error
}

func newRequest(ctx context.Context, task Task) *request {
	return &request{
		task: task,
		ctx:  ctx,
		resp: make(chan error),
	}
}

type Connpool struct {
	wConMut  sync.Mutex
	closeMut sync.Mutex
	conn     net.Conn

	newConn  NewConnFunc
	init     sync.Once
	requests chan *request

	errToResp chan error
	middleReq chan *request
}

func NewConnpool(newConn NewConnFunc) *Connpool {
	return &Connpool{
		newConn:  newConn,
		requests: make(chan *request),

		errToResp: make(chan error),
		middleReq: make(chan *request),
	}
}

func (c *Connpool) reqRespWatcher() {
	for req := range c.requests {

		c.middleReq <- req
		var err error
		select {
		case <-req.ctx.Done():
			err = context.Canceled
			go func() { <-c.errToResp }() // need to clean errToResp before next call
		case err1 := <-c.errToResp:
			err = err1
		}

		if err != nil {
			c.wConMut.Lock()
			if c.conn != nil {
				_ = c.conn.Close()
				c.conn = nil
			}
			c.wConMut.Unlock()
		}

		req.resp <- err
		close(req.resp)
	}
}

func (c *Connpool) reqRespProcessor() {
	for req := range c.middleReq {
		c.wConMut.Lock()
		if c.conn == nil {
			var err error
			c.conn, err = c.newConn(req.ctx)
			if err != nil {
				c.errToResp <- err
				c.wConMut.Unlock()
				continue
			}
		}
		c.wConMut.Unlock()

		err := req.task(req.ctx, c.conn)
		c.errToResp <- err
	}
}

func (c *Connpool) initPoolController() {
	c.init.Do(func() {
		go c.reqRespWatcher()
		go c.reqRespProcessor()
	})
}

func (c *Connpool) Do(ctx context.Context, task Task) error {
	c.closeMut.Lock()
	defer c.closeMut.Unlock()

	c.initPoolController()
	req := newRequest(ctx, task)
	c.requests <- req

	return <-req.resp
}

func (c *Connpool) Close() error {
	c.closeMut.Lock()
	defer c.closeMut.Unlock()

	if !isChannelClosed(c.requests) {
		close(c.requests)
	}

	if !isChannelClosed(c.middleReq) {
		close(c.middleReq)
	}

	var err error
	if c.conn != nil {
		err = c.conn.Close()
	}

	return err
}

func isChannelClosed(ch chan *request) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
