package modbus

import (
	"context"
	"github.com/jt05610/loppu"
	"github.com/jt05610/loppu/hardware"
	"go.uber.org/zap"
	"sync"
)

type Client struct {
	logger *zap.Logger
	dl     *DataLink
	mu     sync.Mutex
	resCh  chan *SerialPDU
}

func (c *Client) Open() error {
	c.resCh = make(chan *SerialPDU)
	return c.dl.serial.Open()
}

func (c *Client) Request(ctx context.Context, addr loppu.Addr, data hardware.Packet) (hardware.Packet, error) {
	s := NewSerialPDU(addr.Byte(), data)
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.dl.Send(s)
	if err != nil {
		panic(err)
	}
	go func() {
		res := &SerialPDU{}
		_, err := c.dl.Recv(res)
		if err != nil {
			panic(err)
		}
		c.resCh <- res
	}()
	select {
	case res := <-c.resCh:
		return res.PDU, nil
	case <-ctx.Done():
		return nil, hardware.ErrTimeout
	}
}

func (c *Client) Close() {
	close(c.resCh)
	c.dl.serial.Close()
}
func DefaultClient(logger *zap.Logger) hardware.Proto {
	ser, err := NewSerial(DefaultSerial, logger)
	if err != nil {
		panic(err)
	}
	return &Client{
		logger: logger,
		dl:     NewDataLink(ser),
		resCh:  make(chan *SerialPDU),
	}
}
