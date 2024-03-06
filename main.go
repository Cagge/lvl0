package main

import (
	"github.com/labstack/gommon/log"
	//_ "github.com/lib/pq"
)

func main() {
	err := control()
	if err != nil {
		log.Error(err)
	}
}
