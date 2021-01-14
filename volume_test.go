package qsc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestGetVolume(t *testing.T) {
	is := is.New(t)

	// MARB-108-DSP1
	dsp := New("10.66.160.250")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vols, err := dsp.Volumes(ctx, []string{"MARB108MEDIAGain", "MARB108MIC1Gain", "MARB108MIC2Gain"})
	fmt.Printf("vols: %v\n", vols)

	is.NoErr(err)
}

func TestGetMutes(t *testing.T) {
	is := is.New(t)

	dsp := New("10.66.160.250")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mutes, err := dsp.Mutes(ctx, []string{"MARB108MIC1Mute", "MARB108MIC2Mute", "MARB108MEDIAMute"})
	fmt.Printf("mutes: %v\n", mutes)

	is.NoErr(err)
}

func TestSetVolume(t *testing.T) {
	is := is.New(t)

	dsp := New("10.66.160.250")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := dsp.SetVolume(ctx, "MARB108MIC1Gain", 65)
	is.NoErr(err)
	vols, err := dsp.Volumes(ctx, []string{"MARB108MIC1Gain"})
	is.NoErr(err)
	is.Equal(vols["MARB108MIC1Gain"], 65)

	err = dsp.SetVolume(ctx, "MARB108MIC1Gain", 35)
	is.NoErr(err)
	vols, err = dsp.Volumes(ctx, []string{"MARB108MIC1Gain"})
	is.NoErr(err)
	is.Equal(vols["MARB108MIC1Gain"], 35)
}

func TestSetMute(t *testing.T) {
	is := is.New(t)

	dsp := New("10.66.160.250")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := dsp.SetMute(ctx, "MARB108MIC1Mute", true)
	is.NoErr(err)
	mutes, err := dsp.Mutes(ctx, []string{"MARB108MIC1Mute"})
	is.NoErr(err)
	is.Equal(mutes["MARB108MIC1Mute"], true)

	err = dsp.SetMute(ctx, "MARB108MIC1Mute", false)
	is.NoErr(err)
	mutes, err = dsp.Mutes(ctx, []string{"MARB108MIC1Mute"})
	is.NoErr(err)
	is.Equal(mutes["MARB108MIC1Mute"], false)
}
