package main

import (
	"fmt"
	"image"
	"log"

	"github.com/gbatanov/golearn/wingui3/winapi"
)

const VERSION = "v0.0.13"

var mouseX, mouseY int = 0, 0
var startMove bool = false

func main() {

	// Конфиг основного окна
	var config winapi.Config
	//	config.Decorated = false
	config.Position = image.Pt(20, 20)
	config.MaxSize = image.Pt(240, 240)
	config.MinSize = image.Pt(240, 100)
	config.Size = image.Pt(240, 100)
	config.Title = "Server check"
	config.EventChan = make(chan winapi.Event, 256)
	config.BorderSize = image.Pt(1, 1)
	config.Mode = winapi.Windowed
	config.BgColor = 0x00dedede

	childConfig := &winapi.Config{
		Title:      "192.168.76.106",
		EventChan:  config.EventChan,
		Size:       image.Pt(int(config.Size.X-10), int(80)),
		MinSize:    config.MinSize,
		MaxSize:    config.MaxSize,
		Position:   image.Pt(int(20), int(15)),
		Mode:       winapi.Windowed,
		BorderSize: image.Pt(0, 0),
		TextColor:  0x0011aa11,
		BgColor:    config.BgColor,
	}

	// Обработчик событий
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

	win, err := winapi.CreateNativeMainWindow(&config)
	if err == nil {
		defer winapi.WinMap.Delete(win.Hwnd)

		// Label с текстом
		chWin, err := winapi.CreateLabel(win, childConfig)
		if err == nil {
			winapi.WinMap.Store(chWin.Hwnd, chWin)
			defer winapi.WinMap.Delete(chWin.Hwnd)
		}
		msg := new(winapi.Msg)
		for winapi.GetMessage(msg, 0, 0, 0) > 0 {
			winapi.TranslateMessage(msg)
			winapi.DispatchMessage(msg)
		}

		close(config.EventChan)
		fmt.Println("Quit")
	} else {
		panic(err.Error())
	}

}

// Обработка событий мыши
func MouseHandler(ev winapi.Event) {
	switch ev.Kind {
	case winapi.Move:
		//		log.Println("Mouse move ", ev.Position)
		/*
			if startMove {

				if ev.SWin.Id == 0 {

					ev.SWin.Config.Position.X += (ev.Position.X)
					ev.SWin.Config.Position.Y += (ev.Position.Y)

					winapi.SetWindowPos( // MoveWindow работает хуже
						ev.SWin.Hwnd,
						0, int32(ev.SWin.Config.Position.X), int32(ev.SWin.Config.Position.Y),
						int32(ev.SWin.Config.Size.X), int32(ev.SWin.Config.Size.Y),
						winapi.SWP_NOSIZE|winapi.SWP_FRAMECHANGED)

					mouseX = ev.Position.X
					mouseY = ev.Position.Y

				}

			}
		*/
	case winapi.Press:
		log.Println("Mouse key press ", ev.Position)
		/*
			if ev.SWin.Id == 0 && ev.SWin.PointerBtns == winapi.ButtonPrimary {
				startMove = true
			}
		*/
	case winapi.Release:
		log.Println("Mouse key release ", ev.Position)
		/*
			if ev.SWin.Id == 0 {
				startMove = false
			}
		*/
	case winapi.Leave:
		log.Println("Mouse lost focus ")
		//		startMove = false
		mouseX = ev.Position.X
		mouseY = ev.Position.Y
	case winapi.Enter:
		//		startMove = false

		log.Println("Mouse enter focus ")
		mouseX = ev.Position.X
		mouseY = ev.Position.Y

	}
}
