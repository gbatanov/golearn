package pinger

import (
	"log"
	"time"

	ping "github.com/prometheus-community/pro-bing"
)

// library "pro-ping" does not work correctly in a loop if created before the first loop
// so I will create it again in each cycle, fortunately this is a very rare process ~1 time per minute

const MINIMAL_PERIOD = 30

type SPinger struct {
	Pinger      *ping.Pinger
	Period      int
	PacketsSent int
	PacketRecv  int
	Flag        bool
	StateChan   chan bool
	Server      string
	Count       int
}

func NewPinger(server string, count int, period int, statechan chan bool) (*SPinger, error) {
	spinger := SPinger{}
	spinger.Flag = false
	spinger.StateChan = statechan
	if period < MINIMAL_PERIOD {
		period = MINIMAL_PERIOD
	}
	spinger.Period = period
	spinger.Server = server
	spinger.Count = count
	return &spinger, nil
}

func (spinger *SPinger) createPinger() (*ping.Pinger, error) {
	pinger, err := ping.NewPinger(spinger.Server)
	if err != nil {
		return &ping.Pinger{}, err
	}
	pinger.Count = spinger.Count
	pinger.SetPrivileged(true)
	pinger.OnFinish = spinger.finish
	pinger.Timeout = 10 * time.Second // Required, otherwise it will hang when the server is turned off
	/*
		pinger.OnRecv = func(pkt *ping.Packet) {
			log.Println("packet") // arrives for each sent packet (3 times with count=3)
		}
	*/
	return pinger, nil
}

func (pinger *SPinger) finish(stats *ping.Statistics) {
	pinger.PacketsSent = stats.PacketsSent
	pinger.PacketRecv = stats.PacketsRecv
	//	log.Printf("%d:%d", stats.PacketsSent, stats.PacketsRecv)
}

func (pinger *SPinger) Run() {
	pinger.Flag = true
	var err error

	for pinger.Flag {
		pinger.Pinger, err = pinger.createPinger()
		if err != nil {
			break
		}
		pinger.PacketsSent = 0
		pinger.PacketRecv = 0

		state := false
		err = pinger.Pinger.Run() // blocks the thread until all pings have completed
		pinger.Pinger.Stop()
		if err == nil {
			state = pinger.PacketsSent > 0 && pinger.PacketRecv > 0
		}
		log.Println(state)
		pinger.StateChan <- state
		n := pinger.Period / 5
		for n > 0 && pinger.Flag {
			time.Sleep(5 * time.Second)
			n = n - 1
		}
	}
	close(pinger.StateChan)

}

func (pinger *SPinger) Stop() {
	pinger.Flag = false
}
