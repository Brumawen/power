package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Power holds the information about the power meter
type Power struct {
	FlashRate  int64     // Number of flashes per KWh
	StartTime  time.Time // Start time
	StartPower float64   // Start power in Kwh
	PulseCount int64     // Number of pulses since start
	LastPulse  time.Time // Time of last pulse
}

// PowerReport holds details about the power that are reported
type PowerReport struct {
	StartTime    time.Time `json:"startTime"`    // Start time
	StartPower   float64   `json:"startPower"`   // Start power in Kwh
	CurrentPower float64   `json:"currentPower"` // Current power in Kwh
	PulseCount   int64     `json:"temp"`         // Number of pulses since start
	LastPulse    time.Time `json:"lastRead"`     // Time of last pulse
}

// GetPowerReport returns a sanitised version of the power data for return to the calling client
func (p *Power) GetPowerReport() PowerReport {
	return PowerReport{
		StartTime:    p.StartTime,
		StartPower:   p.StartPower,
		PulseCount:   p.PulseCount,
		LastPulse:    p.LastPulse,
		CurrentPower: p.GetCurrentPower(),
	}
}

// GetCurrentPower gets the current amount of power left
func (p *Power) GetCurrentPower() float64 {
	// Get amount of power consumed since start
	consumed := float64(p.PulseCount) / float64(p.FlashRate)
	current := p.StartPower - consumed
	return current
}

// LoadCurrentPower reads the current power from the specified file on disk
func (p *Power) LoadCurrentPower(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		b, err := ioutil.ReadFile(path)
		if err == nil {
			var current float64
			buf := bytes.NewReader(b)
			err = binary.Read(buf, binary.LittleEndian, &current)
			if err == nil {
				p.StartPower = current
			}
		}
	}
	p.StartTime = time.Now()
	p.PulseCount = 0
	return err
}

// SaveCurrentPower saves the current power to the specified file
func (p *Power) SaveCurrentPower(path string) error {
	current := p.GetCurrentPower()
	b := new(bytes.Buffer)
	if err := binary.Write(b, binary.LittleEndian, current); err != nil {
		return err
	}
	return ioutil.WriteFile(path, b.Bytes(), 0666)
}

// WriteTo serializes the entity and writes it to the http response
func (p *Power) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}

// WriteTo serializes the entity and writes it to the http response
func (p *PowerReport) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}