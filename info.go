package qsc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/byuoitav/connpool"
	"github.com/fatih/color"
	"go.uber.org/zap"
)

type Info struct {
	Hostname    string
	ModelName   string
	PowerStatus string
	IPAddress   string
}

// Info is all the juicy details about the QSC that everyone is DYING to know about
func (d *DSP) Info(ctx context.Context) (interface{}, error) {

	// toReturn is the struct of Hardware info
	var details Info

	// get the hostname
	addr, e := net.LookupAddr(d.Address)
	if e != nil {
		details.Hostname = d.Address
	} else {
		details.Hostname = strings.Trim(addr[0], ".")
	}

	resp, err := d.GetStatus(ctx)
	if err != nil {
		return details, fmt.Errorf("There was an error getting the status: %v", err)
	}

	d.infof("response", zap.Any("response", resp))
	details.ModelName = resp.Result.Platform
	details.PowerStatus = resp.Result.State

	details.IPAddress = d.Address

	return details, nil
}

// Healthy .
func (d *DSP) Healthy(ctx context.Context) error {
	_, err := d.GetStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed health check: %s", err)
	}

	return nil
}

// GetStatus will be getting responses for us I hope...
func (d *DSP) GetStatus(ctx context.Context) (QSCStatusGetResponse, error) {
	req := d.GetGenericStatusGetRequest(ctx)

	d.infof("In GetStatus...")
	toReturn := QSCStatusGetResponse{}

	toSend, err := json.Marshal(req)
	if err != nil {
		return toReturn, err
	}

	var resp []byte
	err = d.Pool.Do(ctx, func(conn connpool.Conn) error {
		d.infof("getting status")

		conn.SetWriteDeadline(time.Now().Add(3 * time.Second))

		n, err := conn.Write(toSend)
		switch {
		case err != nil:
			return fmt.Errorf("unable to write command to get status: %v", err)
		case n != len(toSend):
			return fmt.Errorf("unable to write command to get status: wrote %v/%v bytes", n, len(toSend))
		}

		deadline, ok := ctx.Deadline()
		if !ok {
			return fmt.Errorf("no deadline set")
		}

		resp, err = conn.ReadUntil('\x00', deadline)
		if err != nil {
			return fmt.Errorf("unable to read response: %w", err)
		}

		d.debugf("Got response: %v", resp)

		return nil
	})
	if err != nil {
		return toReturn, err
	}

	err = json.Unmarshal(resp, &toReturn)
	if err != nil {
		d.infof(color.HiRedString(err.Error()))
	}

	return toReturn, err
}
