package main

import (
	"fmt"
	"iservice/serv"
	"log"
	"os"
	"strings"

	"golang.org/x/sys/windows/svc"
)

const VERSION = "v0.0.1"
const TEST = false
const SVC_NAME = "InteractiveService"

func main() {
	if TEST {
		serv.MainProcess()
	} else {

		inService, err := svc.IsWindowsService()
		if err != nil {
			log.Fatalf("failed to determine if we are running in service: %v", err)
		}

		if inService {
			serv.RunService(SVC_NAME)
			return
		}
		if len(os.Args) < 2 {
			usage("no command specified")
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
		default:
			usage(fmt.Sprintf("invalid command %s", cmd))
		}
		if err != nil {
			log.Fatalf("failed to %s %s: %v", cmd, SVC_NAME, err)
		}

	}
}
func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n usage: %s install|remove|start|stop\n",
		errmsg, os.Args[0])
	os.Exit(2)
}
