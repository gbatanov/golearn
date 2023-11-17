// Run telegram-bot as service on Windows
// Copyright (c) 2023 Georgii Batanov gbatanov@yandex.ru
package main

import (
	"fmt"
	"log"

	"os"
	"strings"

	"git.makves.ru/test/tlgserv/serv"
	"golang.org/x/sys/windows/svc"
)

const VERSION = "0.1.3"

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n usage: %s install|remove|start|stop\n",
		errmsg, os.Args[0])
	os.Exit(2)
}

var svcName = "TlgService"

func main() {

	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in service: %v", err)
	}

	if inService {
		serv.RunService(svcName)
		return
	}

	if len(os.Args) < 2 {
		usage("no command specified")
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "install":
		err = serv.InstallService(svcName, "TlgBot service")
	case "remove":
		err = serv.RemoveService(svcName)
	case "start":
		err = serv.StartService(svcName)
	case "stop":
		err = serv.ControlService(svcName, svc.Stop, svc.Stopped)
	default:
		usage(fmt.Sprintf("invalid command %s", cmd))
	}
	if err != nil {
		log.Fatalf("failed to %s %s: %v", cmd, svcName, err)
	}

}
