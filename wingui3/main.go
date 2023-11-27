package main

import (
	"fmt"
	"image"
	"sync"

	"github.com/gbatanov/golearn/wingui3/winapi"
)

const VERSION = "v0.0.2"

func main() {
	var wg sync.WaitGroup
	var config winapi.Config
	config.Decorated = true
	config.MaxSize = image.Pt(800, 600)
	config.MinSize = image.Pt(100, 100)
	config.Size = image.Pt(320, 120)
	config.Title = "GsbTest"
	config.Wg = &wg

	_, err := winapi.CreateNativeWindow(config)
	if err == nil {
		fmt.Println("New window")
		fmt.Println("Quit")
	} else {
		panic(err.Error())
	}
}
