package main

import (
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"wingui/pinger"

	"wingui/win"

	"fyne.io/systray"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

const VERSION = "v0.0.9"

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
	fmt.Println(runtime.GOOS)
}

func main() {
	className := "WindowClass"
	classNameT := "TextClass"

	instance, err := win.GetModuleHandle()
	if err != nil {
		log.Println(err)
		return
	}

	cursor, err := win.LoadCursorResource(win.IDC_ARROW)
	if err != nil {
		log.Println(err)
		return
	}

	fn := func(hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) uintptr {
		switch msg {
		case win.CWM_CLOSE:
			win.DestroyWindow(hwnd)
		case win.CWM_DESTROY:
			win.PostQuitMessage(0)
		default:
			ret := win.DefWindowProc(hwnd, msg, wparam, lparam)
			return ret
		}
		return 0
	}
	cName, _ := syscall.UTF16PtrFromString(className)
	cNameT, _ := syscall.UTF16PtrFromString(classNameT)

	wcx := win.WndClassEx{
		LpfnWndProc:   syscall.NewCallback(fn),
		HInstance:     instance,
		HCursor:       cursor,
		HbrBackground: win.COLOR_WINDOW + 1,
		LpszClassName: cName,
	}

	wcx.CbSize = uint32(unsafe.Sizeof(wcx))

	if _, err = win.RegisterClassEx(&wcx); err != nil {
		log.Println(err)
		return
	}
	wcxT := win.WndClassEx{
		LpfnWndProc:   syscall.NewCallback(fn),
		HInstance:     instance,
		HCursor:       cursor,
		HbrBackground: win.COLOR_WINDOW + 1,
		LpszClassName: cNameT,
	}

	wcxT.CbSize = uint32(unsafe.Sizeof(wcxT))

	if _, err = win.RegisterClassEx(&wcxT); err != nil {
		log.Println(err)
	}

	mHwnd, err := win.CreateWindow(
		className,
		"Check server",
		win.CWS_VISIBLE|win.CWS_OVERLAPPEDWINDOW,
		100,
		100,
		320,
		80,
		0,
		0,
		instance,
	)
	if err != nil {
		log.Println(err)
		return
	}

	tHwnd, err := win.CreateWindow(
		classNameT,
		"Check server",
		win.CWS_VISIBLE|win.WS_CHILD,
		100,
		100,
		300,
		70,
		mHwnd,
		0,
		instance,
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	win.SetWindowText(tHwnd, server) //Пишет текст в заголовке

	for {
		msg := win.Msg{}
		gotMessage := win.GetMessage(&msg, 0, 0, 0)

		if gotMessage > 0 {
			win.TranslateMessage(&msg)
			win.DispatchMessage(&msg)
		} else {
			break
		}
	}
}

func main2() {

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

		w.Option()
		//		w.Perform(system.ActionMinimize) // сворачивает окно
		err := run(w)
		if err != nil {
			log.Fatal(err)
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
	green := color.NRGBA{R: 0, G: 255, B: 0, A: 255}    // норма
	red := color.NRGBA{R: 255, G: 0, B: 0, A: 255}      // авария
	yellow := color.NRGBA{R: 255, G: 255, B: 0, A: 255} // при старте до получения реального
	titleColor := yellow

	for {
		select { // выбирает либо события окна, либо общие
		case e := <-w.Events():
			fmt.Println(e)
			switch e := e.(type) {
			case system.DestroyEvent:
				fmt.Println("Destroy Event")
				if spinger.Flag {
					spinger.Stop()
				}

				return e.Err
			case system.FrameEvent: //
				log.Println("Frame event")
				//				ui.TransformOp{ui.Offset(f32.Point{X: 100.0, Y: 100.0})}.Add(&ops)
				gtx := layout.NewContext(&ops, e)
				gtx.Reset()

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
				log.Println("channel was closed")
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

func loadImg(path string) ([]byte, error) {
	res, err := os.ReadFile(path)
	return res, err
}
