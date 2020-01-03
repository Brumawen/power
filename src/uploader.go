package main

import (
	"fmt"
	"time"
)

// Uploader uploads the room telemetry to the various destinations
type Uploader struct {
	Srv               *Server   // Current Server
	MqttClient        *Mqtt     // MQTT client
	LastUpdateAttempt time.Time // Last time an update was attempted
	LastUpdate        time.Time // Last time the update was run
	lastValues        *Power    // Last values uploaded for Room
}

// Run is called from the scheduler (ClockWerk). This function will get the latest measurements
// and send the measurements to Thingspeak
func (u *Uploader) Run() {
	if u.MqttClient == nil {
		u.MqttClient = &Mqtt{}
		u.MqttClient.Srv = u.Srv
		u.MqttClient.Initialize()
	}

	if err := u.MqttClient.SendTelemetry(); err != nil {
		u.logError("Error sending telemetry to MQTT")
	}
}

// Close shuts down the Uploader
func (u *Uploader) Close() {
	if u.MqttClient != nil {
		u.MqttClient.Close()
	}
}

// logInfo logs an information message to the logger
func (u *Uploader) logInfo(v ...interface{}) {
	a := fmt.Sprint(v...)
	logger.Info("Uploader: [Inf] ", a)
}

// logError logs an error message to the logger
func (u *Uploader) logError(v ...interface{}) {
	a := fmt.Sprint(v...)
	logger.Error("Uploader [Err] ", a)
}
