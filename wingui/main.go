package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"os/signal"
	"syscall"

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

const VERSION = "v0.0.3"

var quit chan os.Signal

func main() {

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
	go func() {
		systray.Run(onReady, onExit)
	}()
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
