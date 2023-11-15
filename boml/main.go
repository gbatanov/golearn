package main

import (
	"log"

	"github.com/gbatanov/golearn/boml/boml"
)

const VERSION = "v0.0.5"

func main() {
	var config boml.BomlConfig = boml.BomlConfig{}
	filename := "./settings.conf"
	err := config.Load(filename)
	if err != nil {
		log.Println("incorrect file with configuration")
		return
	} else {
		log.Printf("%v", config)
	}

}
