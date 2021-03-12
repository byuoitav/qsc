package qsc

import (
	"os"
	"testing"
)

var dsp *DSP

func TestMain(m *testing.M) {
	dsp = New("EB-425-DSP1.byu.edu")
	os.Exit(m.Run())
}
