package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"

	"github.com/kumarmapanip/toolkit"
)

func main() {
	mux := routes()

	log.Println("Starting server on port 8085")

	http.ListenAndServe(":8085", mux)
}

func routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(handleHtml))
	mux.HandleFunc("/upload", uploadMultipleFiles)
	mux.HandleFunc("/upload-one", uploadOneFile)

	return mux
}

func handleHtml(w http.ResponseWriter, r *http.Request)  {
	temp, err := template.ParseFiles("index.html")
	if err != nil {
		log.Println(err)
		return
	}

	err = temp.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
}
func uploadMultipleFiles(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST" {
		http.Error(w, "Method Not allowed", http.StatusMethodNotAllowed)
	}

	toolkit := toolkit.ToolKit{
		MaxFileSize: int64(math.Pow(2, 30)),
		AllowedFileTypes: []string{"image/jpg", "image/png", "image/gif"},
	}

	uploadedFiles, err := toolkit.UploadFiles(r, "./uploads")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} 

	out := ""
	for _, file := range uploadedFiles {
		out += fmt.Sprintf("Uploaded %s file to uploads folder and renamed to %s\n", file.OriginalFileName, file.NewFileName)
	}

	w.Write([]byte(out))
}

func uploadOneFile(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST" {
		http.Error(w, "Method Not allowed", http.StatusMethodNotAllowed)
	}

	toolkit := toolkit.ToolKit{
		MaxFileSize: int64(math.Pow(2, 30)),
		AllowedFileTypes: []string{"image/jpg", "image/png", "image/gif"},
	}

	uploadedFile, err := toolkit.UploadOneFile(r, "./uploads")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} 

	out := fmt.Sprintf("Uploaded %s file to uploads folder and renamed to %s\n", uploadedFile.OriginalFileName, uploadedFile.NewFileName)

	w.Write([]byte(out))
}