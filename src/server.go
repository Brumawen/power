package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	gopifinder "github.com/brumawen/gopi-finder/src"
	"github.com/gorilla/mux"
	"github.com/kardianos/service"
	"github.com/onatm/clockwerk"
)

// Server defines the Garage web service
type Server struct {
	PortNo         int                  // Port number the server will listen on
	VerboseLogging bool                 // Verbose logging on/off
	Config         *Config              // Configuration settings
	Finder         gopifinder.Finder    // Finder client - used to find other devices
	Uploader       Uploader             // Uploader
	Power          Power                // Power information
	exit           chan struct{}        // Exit flag
	shutdown       chan struct{}        // Shutdown complete flag
	http           *http.Server         // HTTP server
	router         *mux.Router          // HTTP router
	cw             *clockwerk.Clockwerk // Clockwerk scheduler
	isregistering  bool                 // Indicates that a registration is currently ongoing
}

// Start initializes and starts the server running
func (s *Server) Start(v service.Service) error {
	s.logInfo("Service starting")
	app, err := os.Executable()
	if err != nil {
		s.logError("Error getting current executable directory.", err.Error())
	} else {
		wd, err := os.Getwd()
		if err != nil {
			s.logError("Error getting current working directory.", err.Error())
		} else {
			ad := filepath.Dir(app)
			s.logInfo("Current application path is", ad)
			if ad != wd {
				if err := os.Chdir(ad); err != nil {
					s.logError("Error changing working directory.", err.Error())
				}
			}
		}
	}

	// Create a channel that will be used to block until the Stop signal is received
	s.exit = make(chan struct{})
	go s.run()
	return nil
}

// Stop shuts the server down
func (s *Server) Stop(v service.Service) error {
	s.logInfo("Service stopping")
	// Close the channel, this will automatically release the block
	s.shutdown = make(chan struct{})
	close(s.exit)
	// Wait for the shutdown to complete
	_ = <-s.shutdown
	return nil
}

// run will start up and run the service and wait for a Stop signal
func (s *Server) run() {
	if s.PortNo < 0 {
		s.PortNo = 20515
	}

	s.logInfo("Using port no", s.PortNo)

	s.Uploader.Srv = s
	s.Finder.Logger = logger
	s.Finder.VerboseLogging = service.Interactive()

	s.logInfo("Loading Configuration")

	// Get the configuration
	if s.Config == nil {
		s.Config = &Config{}
	}
	s.Config.ReadFromFile("config.json")
	s.Power.FlashRate = s.Config.FlashRate

	s.logInfo("Configuration loaded successfully")

	// Create a router
	s.router = mux.NewRouter().StrictSlash(true)
	s.router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./html/assets"))))

	s.logInfo("Router created")

	// Add the controllers
	s.addController(new(PowerController))
	s.addController(new(LogController))

	s.logInfo("Controllers loaded")

	// Create an HTTP server
	s.http = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.PortNo),
		Handler: s.router,
	}

	// Start the web server
	go func() {
		s.logInfo("Server listening on port", s.PortNo)
		if err := s.http.ListenAndServe(); err != nil {
			msg := err.Error()
			if !strings.Contains(msg, "http: Server closed") {
				s.logError("Error starting Web Server.", msg)
			}
		}
	}()

	go func() {
		// Register service with the Finder server
		go s.RegisterService()

		// Start the scheduler
		s.logInfo("Starting schedule")
		s.startSchedule()
	}()

	// Wait for an exit signal
	_ = <-s.exit

	// Shutdown the HTTP server
	s.http.Shutdown(nil)

	// Shutdown the uploader
	s.Uploader.Close()

	s.logInfo("Shutdown complete")
	close(s.shutdown)
}

func (s *Server) startSchedule() {
	if s.Config.Period <= 0 {
		s.Config.Period = 5
	}
	if s.cw != nil {
		s.cw.Stop()
		s.cw = nil
	}
	s.cw = clockwerk.New()
	s.cw.Every(time.Duration(s.Config.Period) * time.Minute).Do(&s.Uploader)

	s.cw.Start()

	s.logDebug("Schedule set.")

	// Run immedietely
	s.Uploader.Run()
}

func (s *Server) addController(c Controller) {
	c.AddController(s.router, s)
}

// RegisterService will register the service with the devices on the network
func (s *Server) RegisterService() {
	if s.isregistering {
		return
	}
	s.isregistering = true
	isReg := false
	s.logDebug("Starting service registration.")
	for !isReg {
		s.logDebug("RegisterService: Getting device info")
		d, err := gopifinder.NewDeviceInfo()
		if err != nil {
			s.logError("Error getting device info.", err.Error())
		}
		s.logDebug("RegisterService: Creating service")
		sv := d.CreateService("Temperature")
		sv.PortNo = s.PortNo

		if sv.IPAddress == "" {
			s.logDebug("RegisterService: No IP address found.")
		} else {
			s.logDebug("RegisterService: Using IP address", sv.IPAddress)
		}

		s.logDebug("Reg: Finding devices")
		_, err = s.Finder.FindDevices()
		if err != nil {
			s.logError("RegisterService: Error getting list of devices.", err.Error())
		} else {
			if len(s.Finder.Devices) == 0 {
				s.logDebug("RegisterService: Sleeping")
				time.Sleep(15 * time.Second)
			} else {
				// Register the services with the devices
				s.logDebug("RegisterService: Registering the service.")
				s.Finder.RegisterServices([]gopifinder.ServiceInfo{sv})
				isReg = true
			}
		}
	}
	s.logDebug("Completed service registration.")
	s.isregistering = false
}

// logDebug logs a debug message to the logger
func (s *Server) logDebug(v ...interface{}) {
	if s.VerboseLogging {
		a := fmt.Sprint(v...)
		logger.Info("Server: [Dbg] ", a)
	}
}

// logInfo logs an information message to the logger
func (s *Server) logInfo(v ...interface{}) {
	a := fmt.Sprint(v...)
	logger.Info("Server: [Inf] ", a)
}

// logError logs an error message to the logger
func (s *Server) logError(v ...interface{}) {
	a := fmt.Sprint(v...)
	logger.Error("Server [Err] ", a)
}
