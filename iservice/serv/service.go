//go:build windows
// +build windows

package serv

import (
	"fmt"
	"image"
	"iservice/winapi"
	"strings"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var elog debug.Log

type MakvesMemoService struct{}

var test bool = false
var apiEndChan chan bool

var win *winapi.Window

// Конфиг основного окна
var config = winapi.Config{
	Position:   image.Pt(-1, 10),
	MaxSize:    image.Pt(240, 240),
	MinSize:    image.Pt(240, 100),
	Size:       image.Pt(240, 100),
	Title:      "Доступность сервера",
	EventChan:  make(chan winapi.Event, 256),
	BorderSize: image.Pt(1, 1),
	Mode:       winapi.Windowed,
	BgColor:    0x00dedede,
	SysMenu:    false,
}

// Handler обязательно должен реализовывать метод Execute
func (m *MakvesMemoService) Execute(
	args []string,
	r <-chan svc.ChangeRequest,
	changes chan<- svc.Status) (ssec bool, errno uint32) {
	// Извещаем, что стартуем
	changes <- svc.Status{State: svc.StartPending}
	//Стартуем
	go mainProcess()
	// Извещаем, что стартовали успешно
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate: // "допросить", запрашивает текущее состояние сервиса
				// отправляем ему текущее из запроса же. TODO: возвращать реальное текущее состояние.
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
		case <-apiEndChan:
			break loop
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
	apiEndChan = make(chan bool)
	elog.Info(1, fmt.Sprintf("starting %s service", name))

	err = svc.Run(name, &MakvesMemoService{}) // крутится loop в Execute
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	// Если есть какие-то крутящиеся горутины, здесь их надо будет прибить

	elog.Info(1, fmt.Sprintf("%s service stopped", name))
}

// Для тестов без создания сервиса
func MainProcess() {
	test = true
	mainProcess()
}

func mainProcess() {

}
