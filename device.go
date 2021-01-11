package qsc

import (
	"context"
	"net"

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
		return dial.DialContext(ctx, "tcp", d.Address+":1710")
	}

	return d
}
