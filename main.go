package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	//create the handler
	http.HandleFunc("/upload", FileUploadHandler)

	//start server
	log.Fatal(http.ListenAndServe("localhost:8085", nil))
}

//FileUploadHandler takes the files from FilePond and stores them in a temporary folder
func FileUploadHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")

	if req.Method == "OPTIONS" {
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		return
	}

	req.ParseMultipartForm(10 << 20) //10MB max file size
	for _, files := range req.MultipartForm.File {
		for _, file := range files {
			fmt.Printf("Uploading File: %+v\n", file.Filename)
			fmt.Printf("File Size: %+v\n", file.Size)

			//validate the extension
			tokens := strings.Split(file.Filename, ".")
			extension := strings.ToLower(tokens[len(tokens)-1])
			if extension != "png" && extension != "jpg" {
				checkError(errors.New("wrong extension"), http.StatusBadRequest, w)
				return

			}

			//create our destination file
			filename := "upload-*." + extension
			tempFile, err := ioutil.TempFile("temp-images", filename)
			defer tempFile.Close()
			if checkError(err, http.StatusBadRequest, w) {
				return
			}

			//open the uploaded file
			f, err := file.Open()
			if checkError(err, http.StatusInternalServerError, w) {
				return
			}
			//Read the uploaded file
			fileBytes, err := ioutil.ReadAll(f)
			if checkError(err, http.StatusInternalServerError, w) {
				return
			}
			//Write the read bytes to the file created on the server
			tempFile.Write(fileBytes)

			//respond to the client.
			w.WriteHeader(http.StatusOK)

		}

	}

	return

}

func checkError(err error, status int, w http.ResponseWriter) bool {
	if err != nil {
		w.WriteHeader(status)
		w.Write([]byte(err.Error()))
		return true
	}
	return false

}
