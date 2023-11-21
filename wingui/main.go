package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

const VERSION = "v0.0.1"

func main() {
	go func() {
		w := app.NewWindow(app.Title("Server state"), app.MaxSize(240, 80), app.MinSize(240, 80))
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
	//	var cfg app.Config = app.Config{MinSize: image.Pt(10, 10), MaxSize: image.Pt(100, 100)}
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			//			gtx.Constraints.Min = cfg.MinSize
			//			gtx.Constraints.Max = cfg.MaxSize

			title := material.H1(th, "192.168.76.106")
			maroon := color.NRGBA{R: 0, G: 255, B: 0, A: 255}
			title.Color = maroon
			title.Alignment = text.Middle
			title.TextSize = 24.0
			title.Font.Weight = 400
			title.Layout(gtx)

			e.Frame(gtx.Ops)
		}
	}
}
