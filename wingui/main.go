package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

const VERSION = "v0.0.2"

var quit chan os.Signal

func main() {

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
			}
			// Это просто пример использования канала для внешних событий!
			// В реале не использовать ))
		case extE := <-quit:
			log.Println(extE)
			os.Exit(0)
		} //select
	}
}
