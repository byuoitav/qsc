package qsc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/fatih/color"
	"go.uber.org/zap"
)

type Info struct {
	Hostname    string
	ModelName   string
	PowerStatus string
	IPAddress   string
}

// GetInfo is all the juicy details about the QSC that everyone is DYING to know about
func (d *DSP) GetInfo(ctx context.Context) (interface{}, error) {

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

	d.Log.Info("response", zap.Any("response", resp))
	details.ModelName = resp.Result.Platform
	details.PowerStatus = resp.Result.State

	details.IPAddress = d.Address

	return details, nil
}

// GetStatus will be getting responses for us I hope...
func (d *DSP) GetStatus(ctx context.Context) (QSCStatusGetResponse, error) {
	req := d.GetGenericStatusGetRequest(ctx)

	d.Log.Info("In GetStatus...")
	toReturn := QSCStatusGetResponse{}

	resp, err := d.SendCommand(ctx, req)
	if err != nil {
		d.Log.Info(color.HiRedString(err.Error()))
		return toReturn, err
	}

	err = json.Unmarshal(resp, &toReturn)
	if err != nil {
		d.Log.Info(color.HiRedString(err.Error()))
	}

	return toReturn, err
}
