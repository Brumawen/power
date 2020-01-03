package main

import (
	"fmt"
	"net/http"
	"os/exec"

	gopifinder "github.com/brumawen/gopi-finder/src"
	"github.com/gorilla/mux"
)

// LogController handles the Web Methods for reading log records.
type LogController struct {
	Srv *Server
}

// AddController adds the controller routes to the router
func (c *LogController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Methods("GET").Path("/log/get").Name("GetLogs").
		Handler(Logger(c, http.HandlerFunc(c.handleGetLogs)))
}

func (c *LogController) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	myInfo, err := gopifinder.NewDeviceInfo()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if myInfo.OS != "Linux" {
		http.Error(w, "Not supported.", 500)
		return
	}
	out, _ := exec.Command("journalctl", "--no-pager", "-u", "PowerMonitor", "-S", "1 hour ago").CombinedOutput()
	w.Write([]byte(out))
}

// LogInfo is used to log information messages for this controller.
func (c *LogController) LogInfo(v ...interface{}) {
	a := fmt.Sprint(v...)
	logger.Info("LogController: ", a)
}
