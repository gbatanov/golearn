package main

import (
	"iservice/serv"
	"os"
	"strings"

	"golang.org/x/sys/windows/svc"
)

const VERSION = "v0.0.5"
const TEST = false
const SVC_NAME = "InteractiveService"

func main() {
	if TEST {
		serv.MainProcess()
	} else {

		inService, err := svc.IsWindowsService()
		if err != nil {
			return
		}

		if inService {
			serv.RunService(SVC_NAME)
			return
		}

		cmd := strings.ToLower(os.Args[1])
		switch cmd {
		case "install":
			err = serv.InstallService(SVC_NAME, "Interactive service "+VERSION)
		case "remove":
			err = serv.RemoveService(SVC_NAME)
		case "start":
			err = serv.StartService(SVC_NAME)
		case "stop":
			err = serv.ControlService(SVC_NAME, svc.Stop, svc.Stopped)
		}

	}
}
