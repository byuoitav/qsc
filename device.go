package qsc

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/connpool"
)

type DSP struct {
	Address string
	Pool    *connpool.Pool
	Log     Logger
}

const _kTimeoutInSeconds = 2.0

func New(addr string, opts ...Option) *DSP {
	options := options{
		ttl:   _defaultTTL,
		delay: _defaultDelay,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	d := &DSP{
		Address: addr,
		Pool: &connpool.Pool{
			TTL:    options.ttl,
			Delay:  options.delay,
			Logger: options.logger,
		},
		Log: options.logger,
	}

	d.Pool.NewConnection = func(ctx context.Context) (net.Conn, error) {
		dial := net.Dialer{}
		conn, err := dial.DialContext(ctx, "tcp", d.Address+":1710")
		if err != nil {
			return nil, err
		}

		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(5 * time.Second)
		}

		conn.SetDeadline(deadline)

		buf := make([]byte, 64)
		for i := range buf {
			buf[i] = 0x01
		}
		for !bytes.Contains(buf, []byte{0x00}) {
			_, err := conn.Read(buf)
			if err != nil {
				conn.Close()
				return nil, fmt.Errorf("unable to read new connection prompt: %w", err)
			}
		}

		return conn, nil
	}

	return d
}
