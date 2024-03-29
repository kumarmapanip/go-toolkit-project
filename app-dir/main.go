package main

import (
	"fmt"
	"log"

	"github.com/kumarmapanip/toolkit"
)

func main() {
	toolKit := toolkit.ToolKit{}

	err := toolKit.CreateDirectoryIfNotExists("./logs-dir/hello/op")
	if err != nil {
		log.Println(err)
	}

	// slugify
	toSlug := "now is the time 123"
	slug, err := toolKit.Slugify(toSlug)
	fmt.Println(slug)
}