package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Config holds the configuration required for the Service
type Config struct {
	FlashRate    int64  `json:"flashRate"`    // Number of flashes per KWh
	Period       int    `json:"period"`       // Cloud update period (in minutes)
	EnableMqtt   bool   `json:"enableMqtt"`   // Enable MQTT integration
	MqttHost     string `json:"mqttHost"`     // MQTT Host
	MqttUsername string `json:"mqttUsername"` // MQTT Username
	MqttPassword string `json:"mqttPassword"` // MQTT password
}

// ReadFromFile will read the configuration settings from the specified file
func (c *Config) ReadFromFile(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		b, err := ioutil.ReadFile(path)
		if err == nil {
			err = json.Unmarshal(b, &c)
		}
	}
	c.SetDefaults()
	return err
}

// WriteToFile will write the configuration settings to the specified file
func (c *Config) WriteToFile(path string) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0666)
}

// ReadFrom reads the string from the reader and deserializes it into the config values
func (c *Config) ReadFrom(r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err == nil {
		if b != nil && len(b) != 0 {
			err = json.Unmarshal(b, &c)
		}
	}
	c.SetDefaults()
	return err
}

// WriteTo serializes the config and writes it to the http response
func (c *Config) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}

// SetDefaults checks the configuration and makes sure that, if
// a value is not configured, the default value is set.
func (c *Config) SetDefaults() {
	// Set default values, if required
	if c.FlashRate <= 0 {
		c.FlashRate = 1000
	}
}
