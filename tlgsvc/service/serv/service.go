// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file
// https://github.com/SeldonIO/seldon-operator/blob/master/vendor/golang.org/x/sys/LICENSE.
//
// Run TelegramBot as service on Windows
// Copyright (c) 2023 Georgii Batanov gbatanov@yandex.ru
//
//go:build windows
// +build windows

package serv

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

const CMD = "C:\\work\\bin\\tlgsvc.exe"

var elog debug.Log

type TelegramService struct{}

func (m *TelegramService) Execute(
	args []string,
	r <-chan svc.ChangeRequest,
	changes chan<- svc.Status) (ssec bool, errno uint32) {

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	go mainProcess("TelegramBotService")

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Есть известный баг с откатом состояния, поэтому дублируется
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				testOutput := strings.Join(args, "-")
				testOutput += fmt.Sprintf("-%d [GSB]", c.Context)
				elog.Info(1, testOutput)
				break loop
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func RunService(name string) {
	var err error
	elog, err = eventlog.Open(name)
	if err != nil {
		return
	}

	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	run := svc.Run
	err = run(name, &TelegramService{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}

	elog.Info(1, fmt.Sprintf("%s service stopped", name))
}

func mainProcess(name string) {
	cmd := exec.Command(CMD)
	cmd.Run()
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
}
