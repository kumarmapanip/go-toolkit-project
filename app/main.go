package main

import (
	"log"

	"github.com/kumarmapanip/toolkit"
)

func main() {
	tools := &toolkit.ToolKit{} 

	password:=  tools.RandomPasswordGenerator()
	log.Println("My password: ", password)
}