package main

import (
	"fmt"
	"image"
	"sync"

	"github.com/gbatanov/golearn/wingui3/winapi"
)

const VERSION = "v0.0.4"

func main() {
	var wg sync.WaitGroup
	var config winapi.Config
	config.Decorated = true
	config.MaxSize = image.Pt(800, 600)
	config.MinSize = image.Pt(100, 100)
	config.Size = image.Pt(320, 120)
	config.Title = "GsbTest"
	config.Wg = &wg
	config.EventChan = make(chan winapi.Event, 32)

	go func() {
		for {
			ev, ok := <-config.EventChan
			if !ok {
				return
			}
			switch ev.Kind {
			case winapi.Move:
				fmt.Println("Move", ev.Position)
				winapi.SetWindowPos(ev.SWin.Hwnd, 0, int32(ev.Position.X), int32(ev.Position.Y), int32(ev.SWin.Config.Size.X), int32(ev.SWin.Config.Size.Y), winapi.SWP_FRAMECHANGED /*SWP_NOSIZE*/ /*SWP_SHOWWINDOW*/)
			}

		}
	}()

	err := winapi.CreateNativeMainWindow(config)
	if err == nil {
		close(config.EventChan)
		fmt.Println("Quit")
	} else {
		panic(err.Error())
	}

}
