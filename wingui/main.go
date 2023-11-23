package main

import (
	"image/color"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"wingui/pinger"

	"fyne.io/systray"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
)

const VERSION = "v0.0.11"

var server string = "192.168.76.106"
var count = 3
var period = 60 // seconds
var tlgBotService = "http://192.168.76.95:8055/api/?"
var quit chan os.Signal
var stateChan chan bool
var spinger *pinger.SPinger
var err error
var withCaption = true
var imgOk []byte
var imgErr []byte

func init() {
}

func main() {

	quit = make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	stateChan = make(chan bool, 1)
	spinger, err = pinger.NewPinger(server, count, period, stateChan)
	if err != nil {
		panic(err)
	}
	imgOk, err = loadImg("./img/check.ico")
	if err != nil {
		imgOk = make([]byte, 0)
	}
	imgErr, err = loadImg("./img/stop.ico")
	if err != nil {
		imgErr = make([]byte, 0)
	}
	// основное окно
	go func() {
		var w *app.Window
		// app.Decorated(false) - выводит окно без Caption

		w = app.NewWindow(
			app.Title("Server state"),
			app.Size(240, 80),
			app.MaxSize(240, 80),
			app.MinSize(240, 80),
			app.Decorated(withCaption))

		//		w.Perform(system.ActionMinimize) // сворачивает окно
		err := run(w)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}()

	// systray
	go func() {
		systray.Run(onReady, onExit)
		if spinger.Flag {
			spinger.Stop()
		}
	}()

	// pinger
	spinger.Run()
	// Запуск основного окна
	app.Main()
}

func run(w *app.Window) error {
	var msgSent = false // Сообщение уже отправлено
	var oldState = 1    // Предыдущее состояние сервера
	var ops op.Ops
	var title material.LabelStyle // Текст в окне (IP сервера)

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th.Bg.A = 255
	th.Bg.B = 255
	th.Bg.R = 0
	th.Bg.G = 0

	// Цвета IP сервера
	green := color.NRGBA{R: 0, G: 200, B: 0, A: 255}    // норма
	red := color.NRGBA{R: 200, G: 0, B: 0, A: 255}      // авария
	yellow := color.NRGBA{R: 200, G: 200, B: 0, A: 255} // при старте до получения реального
	titleColor := yellow

	for {
		select { // выбирает либо события окна, либо общие
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				if spinger.Flag {
					spinger.Stop()
				}
				return e.Err
			case system.FrameEvent: //
				gtx := layout.NewContext(&ops, e)

				if titleColor == green {
					paint.Fill(&ops, color.NRGBA{R: 128, G: 128, B: 20, A: 128})
				} else if titleColor == red {
					paint.Fill(&ops, color.NRGBA{R: 0, G: 128, B: 0, A: 128})
				} else {
					paint.Fill(&ops, color.NRGBA{R: 128, G: 128, B: 128, A: 128})
				}
				// register a global key listener for the escape key wrapping our entire UI.
				area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
				key.InputOp{
					Tag:  w,
					Keys: key.NameEscape,
				}.Add(gtx.Ops)

				// Выход из программы по Escape
				for _, event := range gtx.Events(w) {
					switch event := event.(type) {
					case key.Event:
						if event.Name == key.NameEscape {
							return nil
						}
					}
				}
				// render and handle UI.
				area.Pop()
				title = material.H1(th, "192.168.76.106")
				title.Color = titleColor
				title.Alignment = text.Middle
				title.TextSize = 28.0
				title.Font.Weight = 400
				// paddings
				inset := layout.Inset{Top: 20, Bottom: 8, Left: 8, Right: 8}
				inset.Layout(gtx, title.Layout)

				e.Frame(gtx.Ops)

			}

		case <-quit:
			return nil
		case state, ok := <-stateChan:
			if !ok {
				return nil
			}
			if state {
				titleColor = green

				w.Invalidate()
				oldState = 1
				msgSent = false

				if len(imgOk) > 0 {
					systray.SetIcon(imgOk)
				}
			} else {
				titleColor = red
				w.Invalidate()
				if oldState == 1 {
					oldState = 0
					if !msgSent {
						msgSent = sendMsg()
					}
					if len(imgErr) > 0 {
						systray.SetIcon(imgErr)
					}

				}
			}

		} //select
	}
}

func onReady() {

	if len(imgErr) > 0 {
		systray.SetIcon(imgErr)
		systray.SetTooltip("Check Server Health")
	}
	systray.SetTitle("Check Server")
	mQuit := systray.AddMenuItem("Quit", "Выход")
	mQuit.Enable()
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

}

func onExit() {
	quit <- syscall.SIGTERM
}

// Send message to telegram
func sendMsg() bool {
	client := http.Client{}
	client.Timeout = 10 * time.Second

	params := url.Values{}
	params.Add("msg", "server_invalid")
	encodedData := params.Encode()
	body := strings.NewReader(encodedData)

	req, _ := http.NewRequest("POST", tlgBotService, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(encodedData)))

	resp, err := client.Do(req)

	/*
		url := tlgBotService + "msg=server_invalid"
		resp, err := client.Get(url)
	*/
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}

func loadImg(path string) ([]byte, error) {
	res, err := os.ReadFile(path)
	return res, err
}
