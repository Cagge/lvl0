package main

import (
	"github.com/labstack/gommon/log"
)

func main() {
	err := control()
	if err != nil {
		log.Error(err)
	}
}
