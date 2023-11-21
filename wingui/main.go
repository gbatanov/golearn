package main

import (
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wingui/pinger"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

const VERSION = "v0.0.4"

var count = 3
var period = 60 // seconds
var tlgBotService = "http://192.168.76.95:8055/api/?"
var quit chan os.Signal
var stateChan chan bool
var spinger *pinger.SPinger
var err error

func main() {
	server := "192.168.76.106"
	quit = make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	withCaption := true
	go func() {
		var w *app.Window
		//		statusColor := color.NRGBA{R: 255, G: 255, B: 0, A: 128}
		// app.Decorated(false) - выводит окно без Caption

		w = app.NewWindow(
			app.Title("Server state"),
			app.Size(240, 80),
			app.MaxSize(240, 80),
			app.MinSize(240, 80),
			app.Decorated(withCaption))

		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	stateChan = make(chan bool, 1)
	spinger, err = pinger.NewPinger(server, count, period, stateChan)
	if err != nil {
		panic(err)
	}

	go func() {
		systray.Run(onReady, onExit)
		if spinger.Flag {
			spinger.Stop()
		}
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
	spinger.Run()
	app.Main()
}

func run(w *app.Window) error {

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th.Bg.A = 255
	th.Bg.B = 255
	th.Bg.R = 0
	th.Bg.G = 0
	var ops op.Ops
	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {

			case system.DestroyEvent:
				fmt.Println("system.DestroyEvent")
				if spinger.Flag {
					spinger.Stop()
				}

				return e.Err
			case system.FrameEvent:

				gtx := layout.NewContext(&ops, e)
				//			inset := layout.Inset{Top: 8, ...}
				title := material.H1(th, "192.168.76.106")
				maroon := color.NRGBA{R: 0, G: 255, B: 0, A: 255}

				title.Color = maroon
				title.Alignment = text.Middle
				title.TextSize = 28.0
				title.Font.Weight = 400
				// paddings
				inset := layout.Inset{Top: 20, Bottom: 8, Left: 8, Right: 8}
				inset.Layout(gtx, title.Layout)
				e.Frame(gtx.Ops)

				//			default:
				//				fmt.Println(e)
			}

			// Это просто пример использования канала для внешних событий!
			// В реале не использовать ))
		case <-quit:
			return nil

		} //select
	}
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Check Server")
	systray.SetTooltip("Check Server Health")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	mQuit.Enable()
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()

	}()
	// Sets the icon of a menu item.
	mQuit.SetIcon(icon.Data)
}

func onExit() {
	quit <- syscall.SIGTERM
}

// Send message to telegram
func sendMsg() bool {
	client := http.Client{}
	client.Timeout = 10 * time.Second
	/*
		params := url.Values{}
		params.Add("msg", "server_invalid")
		encodedData := params.Encode()
		body := strings.NewReader(encodedData)

		req, _ := http.NewRequest("POST",tlgBotService, body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Content-Length", strconv.Itoa(len(encodedData)))

		resp, err := client.Do(req)
	*/

	url := tlgBotService + "msg=server_invalid"
	resp, err := client.Get(url)

	if err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}
