package qsc

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap/zaptest"
)

func TestControl(t *testing.T) {
	is := is.New(t)
	log := zaptest.NewLogger(t)
	dsp.log = log
	dsp.pool.Logger = log.Sugar()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	block := "RoomCombine"
	setAndCheck := func(val float64) {
		is.NoErr(dsp.SetControl(ctx, block, val))

		cur, err := dsp.Control(ctx, block)
		is.NoErr(err)
		is.True(cur == val)
	}

	start, err := dsp.Control(ctx, block)
	is.NoErr(err)

	setAndCheck(1)
	setAndCheck(0)
	setAndCheck(start)
}
