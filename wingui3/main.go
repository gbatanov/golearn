package main

import (
	"fmt"
	"image"
	"sync"

	"github.com/gbatanov/golearn/wingui3/winapi"
)

const VERSION = "v0.0.1"

func main() {
	var wg sync.WaitGroup
	var config winapi.Config
	config.Decorated = true
	config.MaxSize = image.Pt(100, 100)
	config.MinSize = image.Pt(100, 100)
	config.Size = image.Pt(100, 100)
	config.Title = "GsbTest"
	config.Wg = &wg

	win, err := winapi.CreateNativeWindow(config)
	if err == nil {
		fmt.Println("New window")
		fmt.Println(win)
	} else {
		panic(err.Error())
	}
}
