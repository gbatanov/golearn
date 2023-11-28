package main

import (
	"fmt"
	"image"
	"log"
	"sync"

	"github.com/gbatanov/golearn/wingui3/winapi"
)

const VERSION = "v0.0.5"

var mouseX, mouseY int = 0, 0
var startMove bool = false

func main() {
	var wg sync.WaitGroup
	var config winapi.Config
	config.Decorated = true
	config.Position = image.Pt(20, 20)
	config.MaxSize = image.Pt(800, 600)
	config.MinSize = image.Pt(100, 100)
	config.Size = image.Pt(320, 120)
	config.Title = "GsbTest"
	config.Wg = &wg
	config.EventChan = make(chan winapi.Event, 128)

	go func() {
		for {
			ev, ok := <-config.EventChan
			if !ok {
				return
			}
			switch ev.Source {
			case winapi.Mouse:
				MouseHandler(ev)
			}

		}
	}()

	err := winapi.CreateNativeMainWindow(&config)
	if err == nil {
		close(config.EventChan)
		fmt.Println("Quit")
	} else {
		panic(err.Error())
	}

}

// События мыши
func MouseHandler(ev winapi.Event) {
	switch ev.Kind {
	case winapi.Move:
		log.Println("Mouse move ", ev.Position)
		if startMove {

			ev.SWin.Config.Position.X += (ev.Position.X)
			ev.SWin.Config.Position.Y += (ev.Position.Y)
			winapi.SetWindowPos(
				ev.SWin.Hwnd,
				0, int32(ev.SWin.Config.Position.X), int32(ev.SWin.Config.Position.Y),
				int32(ev.SWin.Config.Size.X), int32(ev.SWin.Config.Size.Y),
				winapi.SWP_NOSIZE)
			mouseX = ev.Position.X
			mouseY = ev.Position.Y

		}

	case winapi.Press:
		log.Println("Mouse key press ", ev.Position)
		if ev.SWin.PointerBtns == winapi.ButtonPrimary {
			startMove = true
		}
	case winapi.Release:
		log.Println("Mouse key release ", ev.Position)
		startMove = false

	}
}
