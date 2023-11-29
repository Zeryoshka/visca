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
	wConMut sync.Mutex
	conn    net.Conn

	newConn  NewConnFunc
	init     sync.Once
	requests chan *request
}

func NewConnpool(newConn NewConnFunc) *Connpool {
	return &Connpool{
		newConn:  newConn,
		requests: make(chan *request),
	}
}

func (c *Connpool) reqRespWatcher(middleReqs chan<- *request, errToResp <-chan error) {
	for req := range c.requests {

		middleReqs <- req
		var err error
		select {
		case <-req.ctx.Done():
			err = context.Canceled
			go func() { <-errToResp }() // need to clean errToResp before next call
		case err1 := <-errToResp:
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

func (c *Connpool) reqRespProcessor(middleReq <-chan *request, errToResp chan<- error) {
	for req := range middleReq {
		c.wConMut.Lock()
		if c.conn == nil {
			var err error
			c.conn, err = c.newConn(req.ctx)
			if err != nil {
				errToResp <- err
				c.wConMut.Unlock()
				continue
			}
		}
		c.wConMut.Unlock()

		err := req.task(req.ctx, c.conn)
		errToResp <- err
	}
}

func (c *Connpool) initPoolController() {
	c.init.Do(func() {
		errToResp := make(chan error)
		middleReq := make(chan *request)

		go c.reqRespWatcher(middleReq, errToResp)
		go c.reqRespProcessor(middleReq, errToResp)
	})
}

func (c *Connpool) Do(ctx context.Context, task Task) error {
	c.initPoolController()

	req := newRequest(ctx, task)
	c.requests <- req

	return <-req.resp
}
