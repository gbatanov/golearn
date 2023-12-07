package pinger

import (
	"time"

	ping "github.com/prometheus-community/pro-bing"
)

// library "pro-ping" does not work correctly in a loop if created before the first loop
// so I will create it again in each cycle, fortunately this is a very rare process ~1 time per minute

const MINIMAL_PERIOD = 30

type SPinger struct {
	Pinger      []*ping.Pinger
	Period      int
	PacketsSent int
	PacketRecv  int
	Flag        bool
	StateChan   chan map[int]int
	Server      []string
	Count       int
}

func NewPinger(server []string, count int, period int, statechan chan map[int]int) (*SPinger, error) {
	spinger := SPinger{}
	spinger.Flag = false
	spinger.StateChan = statechan
	if period < MINIMAL_PERIOD {
		period = MINIMAL_PERIOD
	}
	spinger.Pinger = make([]*ping.Pinger, len(server))
	spinger.Period = period
	spinger.Server = server
	spinger.Count = count
	return &spinger, nil
}

// Создание пингера для конкретного сервера
func (spinger *SPinger) createPinger(i int) (*ping.Pinger, error) {
	var pinger *ping.Pinger
	server := spinger.Server[i]
	pinger, err := ping.NewPinger(server)
	if err != nil {
		return nil, err
	}
	pinger.Count = spinger.Count
	pinger.SetPrivileged(true)
	pinger.OnFinish = spinger.finish
	pinger.Timeout = 10 * time.Second // Required, otherwise it will hang when the server is turned off

	return pinger, nil
}

func (pinger *SPinger) finish(stats *ping.Statistics) {
	pinger.PacketsSent = stats.PacketsSent
	pinger.PacketRecv = stats.PacketsRecv
}

func (pinger *SPinger) Run() {
	pinger.Flag = true
	var err error
	var st map[int]int = make(map[int]int, len(pinger.Server))
	for pinger.Flag {
		for i := 0; i < len(pinger.Server); i++ {
			st[i] = 0
			pinger.Pinger[i], err = pinger.createPinger(i)

			if err == nil {

				pinger.PacketsSent = 0
				pinger.PacketRecv = 0

				state := 0
				err = pinger.Pinger[i].Run() // blocks the thread until all pings have completed
				pinger.Pinger[i].Stop()
				if err == nil {
					if pinger.PacketsSent > 0 && pinger.PacketRecv > 0 {
						state = 1
					}
				}
				st[i] = state
			}
		}

		// Отправляем целиком всю пачку в обработчик
		pinger.StateChan <- st
		n := pinger.Period / 5
		for n > 0 && pinger.Flag {
			time.Sleep(5 * time.Second)
			n = n - 1
		}
	}
	//	close(pinger.StateChan)

}

func (pinger *SPinger) Stop() {
	pinger.Flag = false
}
