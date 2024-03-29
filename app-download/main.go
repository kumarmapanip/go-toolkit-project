package main

import (
	"fmt"
	//"math/rand"
	"net/http"
	// "time"

	"github.com/kumarmapanip/toolkit"
)

func main() {
	// start router
	mux := routers()

	err := http.ListenAndServe(":8089", mux)
	if err != nil {
		panic(err)
	}
}

func routers() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/download", downloadFile)

	fmt.Println("started web server at 8089")
	http.ListenAndServe(":8089", mux)

	return mux
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	toolKit := toolkit.ToolKit{}

	// source := rand.NewSource(time.Now().UnixNano()) 
	// randGen := rand.New(source)

	// toolKit.DownloadStaticFile(w, r, "./files", "gFlknSf1Ag.png", fmt.Sprintf("download-%d.png", randGen.Intn(100)))

	toolKit.DownloadStaticFile(w, r, "./files", "gFlknSf1Ag.png", "test-download.png")

}
