// +build windows

package main

import (
	"log"
	"os/exec"

	"github.com/kardianos/service"
	"os"
)

var logger service.Logger

type program struct {
	ExePath    string
	ConfigPath string
	SGAccel    *exec.Cmd
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	logger.Infof("Starting the Sync Gateway Accelerator Accelerator service using command: `%s %s`", p.ExePath, p.ConfigPath)

	if p.ConfigPath != "" {
		p.SGAccel = exec.Command(p.ExePath, p.ConfigPath)
	} else {
		p.SGAccel = exec.Command(p.ExePath)
	}
	err := p.SGAccel.Start()
	if err != nil {
		logger.Errorf("Failed to start Sync Gateway Accelerator due to error %v", err)
		return
	}
	err = p.SGAccel.Wait()
	if err != nil {
		logger.Errorf("Sync Gateway Accelerator exiting with %v", err)
		panic("Failed to start Sync Gateway Accelerator service.")
	}
}

func (p *program) Stop(s service.Service) error {
	logger.Infof("Stopping Sync Gateway Accelerator Accelerator service...")
	p.SGAccel.Process.Kill()
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "SGAccel",
		DisplayName: "Couchbase Sync Gateway Accelerator Accelerator",
		Description: "Couchbase Sync Gateway Accelerator Accelerator performance and scalability service.",
	}

	var exePath string
	var configPath string

	switch len(os.Args) {
	case 2:
		exePath = "C:\\Program Files (x86)\\Couchbase\\sg_accel.exe" // Uses default binary image path
		svcConfig.Arguments = []string{"start"}                      // Uses the default config
	case 3:
		exePath = "C:\\Program Files (x86)\\Couchbase\\sg_accel.exe" // Uses default binary image path
		configPath = os.Args[2]                                      // Uses custom config
		svcConfig.Arguments = []string{"start", configPath}
	case 4:
		exePath = os.Args[2]    // Uses custom binary image path
		configPath = os.Args[3] // Uses custom config
		svcConfig.Arguments = []string{"start", exePath, configPath}
	default:
		panic("Valid parameters combinations are: COMMAND [none, custom config path, or custom exe path and custom config path].")
	}

	prg := &program{
		ExePath:    exePath,
		ConfigPath: configPath,
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case os.Args[1] == "install":
		logger.Info("Installing Sync Gateway Accelerator")
		err := s.Install()
		if err != nil {
			logger.Errorf("Failed to install Sync Gateway Accelerator service: %s", err)
		}
		return
	case os.Args[1] == "uninstall":
		logger.Info("Uninstalling Sync Gateway Accelerator")
		err := s.Uninstall()
		if err != nil {
			logger.Errorf("Failed to uninstall Sync Gateway Accelerator service: %s", err)
		}
		return
	case os.Args[1] == "stop":
		err := s.Stop()
		if err != nil {
			logger.Errorf("Failed to stop Sync Gateway Accelerator service: %s", err)
		}
		return
	case os.Args[1] == "restart":
		err := s.Restart()
		if err != nil {
			logger.Errorf("Failed to restart Sync Gateway Accelerator service: %s", err)
		}
		return
	}

	err = s.Run()

	if err != nil {
		logger.Error(err)
		s.Stop()
	}
	logger.Infof("Exiting Sync Gateway Accelerator service.")
}
