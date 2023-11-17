package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.makves.ru/test/check-server/pinger"
)

const VERSION = "0.1.4"

var count = 3
var period = 60 // seconds

var tlgBotService = "http://192.168.76.95:8055/api/?"

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: \n\tcheck-server.exe <serverIP>\nExample: \n\tcheck-server.exe 192.168.76.106")
		return
	}

	server := os.Args[1]

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	stateChan := make(chan bool, 1)
	pinger, err := pinger.NewPinger(server, count, period, stateChan)
	if err != nil {
		panic(err)
	}
	go func() {
		<-quit
		pinger.Stop()
	}()

	// Receiving the state of the monitored server
	go func() {
		msgSent := false
		oldState := 1
		for {
			state, ok := <-stateChan
			if !ok {
				log.Println("channel was closed")
				return
			}
			if state {
				log.Println("server alive")
				oldState = 1
				msgSent = false
			} else {
				log.Println("server don't responce")
				if oldState == 1 {
					oldState = 0
					if !msgSent {
						msgSent = sendMsg()
					}
				}
			}
		}
	}()

	//	log.Println("check run")
	pinger.Run()
	log.Println("check end")
}

// Send message to telegram
func sendMsg() bool {
	client := http.Client{}
	client.Timeout = 10 * time.Second
	url := tlgBotService + "msg=server_invalid"
	resp, err := client.Get(url)
	//	log.Println("sendMsg end")
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}
