package main

import (
	"image"

	"context"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"fyne.io/systray"

	"git.makves.ru/test/check-server/img"
	"git.makves.ru/test/check-server/pinger"
	"git.makves.ru/test/check-server/winapi"
)

const VERSION = "v0.6.33"

const COLOR_GREEN = 0x0011aa11
const COLOR_RED = 0x000000c8
const COLOR_YELLOW = 0x0000c8c8
const COLOR_GRAYDE = 0x00dedede

var COLOR_STATES map[int]uint32 = map[int]uint32{0: COLOR_RED, 1: COLOR_GREEN}

const SERVER_MAX_COUNT = 8

// Контролируемые сервера (дефолтный список, если нет файла с конфигурацией)
var serverList []string = []string{"192.168.76.106"}

const PINGER_COUNT = 4
const PINGER_PERIOD = 60 // seconds
// API телеграм-бота
var tlgBotService = "http://192.168.76.95:8055/api/?"

// Адрес, по которому долбится прометей
var httpServer string = "192.168.76.95:8280"
var srv http.Server      // Создаем переменную для HTTP-сервера
var states map[int]int   // Состояния контролируемых серверов
var oldState map[int]int // Предыдущее состояние сервера

var quit chan os.Signal
var stateChan chan map[int]int
var spinger *pinger.SPinger
var err error
var flag = true

var win *winapi.Window

// Конфиг основного окна
var config = winapi.Config{
	Position:   image.Pt(-1, 10),
	MaxSize:    image.Pt(240, 360),
	MinSize:    image.Pt(240, 100),
	Size:       image.Pt(240, 100),
	Title:      "Доступность сервера",
	EventChan:  make(chan winapi.Event, 256),
	BorderSize: image.Pt(1, 1),
	Mode:       winapi.Windowed,
	BgColor:    COLOR_GRAYDE,
	SysMenu:    false,
}
var labelConfig = winapi.Config{
	Title:      "Child",
	EventChan:  config.EventChan,
	Size:       image.Pt(int(config.Size.X-10), int(30)),
	MinSize:    config.MinSize,
	MaxSize:    config.MaxSize,
	Position:   image.Pt(int(18), int(15)),
	Mode:       winapi.Windowed,
	BorderSize: image.Pt(0, 0),
	TextColor:  COLOR_YELLOW,
	BgColor:    config.BgColor,
	SysMenu:    false,
}

var NoClose = false

func main() {

	quit = make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	err = getConfig()

	// пингер
	oldState = make(map[int]int, SERVER_MAX_COUNT) // Предыдущее состояние сервера
	states = make(map[int]int, SERVER_MAX_COUNT)
	stateChan = make(chan map[int]int, 1)

	// systray
	go func() {
		systray.Run(onReady, onExit)
	}()

	// http-server для отдачи метрики
	go func() {
		// Создаем маршрутизатор
		mux := http.NewServeMux()
		// Наполняем его обрабатываемыми маршрутами
		mux.HandleFunc("/metrics", metrics)
		srv.Addr = httpServer
		srv.Handler = mux
		srv.ListenAndServe()
	}()

	for {
		NoClose = false
		createTask()
		if !NoClose {
			break
		}
	}
	srv.Shutdown(context.Background())

	close(config.EventChan)
}

func recreateTask() {
	NoClose = true
	winapi.SendMessage(win.Hwnd, winapi.WM_CLOSE, 0, 0)
}

func createTask() {

	spinger, err = pinger.NewPinger(serverList, PINGER_COUNT, PINGER_PERIOD, stateChan)
	if err != nil {
		panic(err)
	}

	// основное окно
	win, err = winapi.CreateNativeMainWindow(config)
	if err == nil {
		defer winapi.WinMap.Delete(win.Hwnd)

		Reconfig()

		// pinger
		go spinger.Run()
		// Обработчик событий
		go func() {
			run(win)
			if spinger.Flag {
				spinger.Stop()
			}
			winapi.SendMessage(win.Hwnd, winapi.WM_CLOSE, 0, 0)
		}()

		msg := new(winapi.Msg)
		for flag && (winapi.GetMessage(msg, 0, 0, 0) > 0) {
			winapi.TranslateMessage(msg)
			winapi.DispatchMessage(msg)
		}

	} else {
		panic(err)
	}

}

// Отдаем текущее состояние контролируемого сервера
func metrics(w http.ResponseWriter, r *http.Request) {
	var sout string = ""

	for _, w := range win.Childrens {
		_, exists := states[int(w.Id)]
		if exists {
			fState := 0.0
			if states[int(w.Id)] == 1 {
				fState = float64(states[int(w.Id)]) - float64(0.1*float64(w.Id))
			} else if states[int(w.Id)] == 0 {
				fState = float64(states[int(w.Id)]) + float64(0.1*float64(w.Id))
			}
			sout += "alive{sname=\"" + w.Config.Title + "\"} " + strconv.FormatFloat(fState, 'f', 3, 32) + "\n"
		} else {
			continue
		}
	}
	if sout == "" {
		w.Write([]byte(""))
		return
	}

	w.Write([]byte(sout))
}

// Добавление строки с IP контролируемого сервера
func AddLabel(win *winapi.Window, lblConfig winapi.Config, id int) error {

	lblConfig.Position.Y = 10 + (lblConfig.Size.Y)*(id)
	chWin, err := winapi.CreateLabel(win, lblConfig, id)
	if err == nil {
		winapi.WinMap.Store(chWin.Hwnd, chWin)
		defer winapi.WinMap.Delete(chWin.Hwnd)
		win.Childrens = append(win.Childrens, chWin)

		return nil
	}
	return err
}

// Основной обработчик событий
// Завершение это функции инициирует отправку сообщения WM_CLOSE
func run(w *winapi.Window) error {

	for i := range serverList {
		oldState[i] = -1
		w.Childrens[i].Config.TextColor = COLOR_YELLOW
	}

	for {
		select { // выбирает либо события окна, либо общие
		case ev, ok := <-config.EventChan: // оконные события
			if !ok {
				return nil
			}
			switch ev.Source {

			case winapi.Frame: //
				switch ev.Kind {
				case winapi.Destroy: // Сообщение закрытия окна
					return nil
				}
			case winapi.Mouse:

			} // switch ev.Source

		case <-quit: // сообщение при закрытии трея
			return nil

		case stateIn, ok := <-stateChan: // сообщение от пингера
			if !ok {
				return nil
			}

			for id, state := range stateIn {
				if state == 0 || state == 1 {
					states[id] = state
					w.Childrens[id].Config.TextColor = COLOR_STATES[state]

					if oldState[id] != state {
						oldState[id] = state
						// одноразовое сообщение в телеграм
						sendMsgTlg(w.Childrens[id].Config.Title, state)

						if len(img.OkIco) > 0 && len(img.ErrIco) > 0 {
							if state == 1 {
								systray.SetIcon(img.OkIco)
							} else {
								systray.SetIcon(img.ErrIco)
							}
						}
					}
				}
			}
			w.Invalidate()
		} //select
	}
}

// трей готов к работе
func onReady() {

	if len(img.ErrIco) > 0 {
		systray.SetIcon(img.ErrIco)
		systray.SetTooltip("Состояние сервера")
	}
	systray.SetTitle("Check Server")
	mQuit := systray.AddMenuItem("Quit", "Выход")
	mQuit.Enable()
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
	systray.AddSeparator()
	mReconfig := systray.AddMenuItem("Reconfig", "Перечитать конфиг")
	mReconfig.Enable()
	go func() {
		for flag {
			<-mReconfig.ClickedCh
			res := getConfig()
			if res != nil {
				mReconfig.Disable()
			} else {
				recreateTask()
			}
		}
	}()
}

// Обработчик завершения трея
func onExit() {
	quit <- syscall.SIGTERM
	flag = false
}

// Читаем конфиг из файла
func getConfig() error {

	info, err := os.ReadFile("config.conf")
	if err != nil || len(info) < 7 {
		return err
	}
	infoS := string(info)
	infoS = strings.ReplaceAll(infoS, " ", "")
	infoS = strings.ReplaceAll(infoS, "\r", "\n")
	infoS = strings.ReplaceAll(infoS, "\n\n", "\n")
	infoS = strings.Trim(infoS, "\n ")
	serverList = strings.Split(infoS, "\n")

	return nil
}

// Отправка сообщения в телеграм  через внутренний сервис отправки телеграм-сообщений
func sendMsgTlg(checkServer string, state int) bool {
	client := http.Client{}
	client.Timeout = 10 * time.Second

	params := url.Values{}
	if state == 1 {
		params.Add("msg", checkServer+" ok")
	} else {
		params.Add("msg", checkServer+" invalid")
	}
	encodedData := params.Encode()
	body := strings.NewReader(encodedData)

	req, _ := http.NewRequest("POST", tlgBotService, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(encodedData)))

	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}

// Конфигурируем окно по конфигу
func Reconfig() {

	for id, title := range serverList {
		labelConfig.Title = title
		AddLabel(win, labelConfig, id)
	}
	win.Config.Size.Y = labelConfig.Size.Y * (len(serverList) + 2)
	win.Config.MinSize.Y = win.Config.Size.Y
	win.Config.MaxSize.Y = win.Config.Size.Y

	winapi.SetWindowPos(win.Hwnd,
		winapi.HWND_TOPMOST,
		int32(win.Config.Position.X),
		int32(win.Config.Position.Y),
		int32(win.Config.Size.X),
		int32(win.Config.Size.Y),
		winapi.SWP_NOMOVE)
	win.Invalidate()
}
