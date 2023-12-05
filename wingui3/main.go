package main

import (
	"fmt"
	"image"
	"log"

	"github.com/gbatanov/golearn/wingui3/winapi"
)

const VERSION = "v0.0.18"

var mouseX, mouseY int = 0, 0
var serverList []string = []string{"192.168.76.106", "192.168.76.80"}

// Конфиг основного окна
var config = winapi.Config{
	Position:   image.Pt(20, 20),
	MaxSize:    image.Pt(240, 240),
	MinSize:    image.Pt(240, 100),
	Size:       image.Pt(240, 100),
	Title:      "Server check",
	EventChan:  make(chan winapi.Event, 256),
	BorderSize: image.Pt(1, 1),
	Mode:       winapi.Windowed,
	BgColor:    0x00dedede,
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
	TextColor:  0x0011aa11,
	BgColor:    config.BgColor,
}

func main() {

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

	win, err := winapi.CreateNativeMainWindow(config)
	if err == nil {
		defer winapi.WinMap.Delete(win.Hwnd)

		// Label с текстом
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

func AddLabel(win *winapi.Window, lblConfig winapi.Config, id int) error {

	lblConfig.Position.Y = 10 + (lblConfig.Size.Y)*(id)
	chWin, err := winapi.CreateLabel(win, lblConfig, id)
	if err == nil {
		winapi.WinMap.Store(chWin.Hwnd, chWin)
		defer winapi.WinMap.Delete(chWin.Hwnd)
		win.Childrens[id] = chWin

		return nil
	}
	return err
}

// Обработка событий мыши
func MouseHandler(ev winapi.Event) {
	mouseX = ev.Position.X
	mouseY = ev.Position.Y

	switch ev.Kind {
	case winapi.Move:
		//		log.Println("Mouse move ", ev.Position)
	case winapi.Press:
		log.Println("Mouse key press ", ev.Position)
	case winapi.Release:
		log.Println("Mouse key release ", ev.Position)
	case winapi.Leave:
		log.Println("Mouse lost focus ")
	case winapi.Enter:
		log.Println("Mouse enter focus ")

	}
}
