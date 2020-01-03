package main

import "github.com/gorilla/mux"

// Controller defines an interface for a Web Method controller
type Controller interface {
	AddController(router *mux.Router, s *Server)
	LogInfo(v ...interface{})
}
