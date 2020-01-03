package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// PowerController handles the Web Methods for the room being monitored
type PowerController struct {
	Srv *Server
}

// AddController adds the controller routes to the router
func (c *PowerController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Methods("GET").Path("/power/get").Name("GetPower").
		Handler(Logger(c, http.HandlerFunc(c.handleGetPower)))
}

// handleGetPower will return the current power status
func (c *PowerController) handleGetPower(w http.ResponseWriter, r *http.Request) {
	rep := c.Srv.Power.GetPowerReport()

	if err := rep.WriteTo(w); err != nil {
		c.LogError("Error serializing power.", err.Error())
		http.Error(w, "Error serializing power", http.StatusInternalServerError)
	}
}

// LogInfo is used to log information messages for this controller.
func (c *PowerController) LogInfo(v ...interface{}) {
	a := fmt.Sprint(v...)
	logger.Info("PowerController: [Inf] ", a)
}

// LogError is used to log information messages for this controller.
func (c *PowerController) LogError(v ...interface{}) {
	a := fmt.Sprint(v...)
	logger.Info("PowerController: [Err] ", a)
}
