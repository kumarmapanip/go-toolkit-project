package main

import (
	"log"

	"github.com/kumarmapanip/toolkit"
)

func main() {
	toolKit := toolkit.ToolKit{}

	err := toolKit.CreateDirectoryIfNotExists("./logs-dir/hello/op")
	if err != nil {
		log.Println(err)
	}
}