package main

import (
	"log"
	"time"
)

const VERSION = "v0.0.1"

func main() {
	var startTimer chan bool = make(chan bool, 1)
	go func() {
		var timer1 *time.Timer = &time.Timer{}
		for {
			select {
			case <-timer1.C:
				log.Println("Сработал таймер")
			case res, ok := <-startTimer:
				if !ok {
					return
				}
				if res {
					log.Println("Start timer")
					timer1 = time.NewTimer(15 * time.Second)
				}
			}
		}
	}()
	startTimer <- true
	time.Sleep(10 * time.Second)

	startTimer <- true
	time.Sleep(20 * time.Second)

	close(startTimer)
	log.Println("Quit")
}
