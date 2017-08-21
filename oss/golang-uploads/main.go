package main

import (
	"bufio"
	"encoding/base64"
	"html/template"
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
)

var views = template.Must(template.ParseGlob("views/*.html"))

func main() {
	log.SetHandler(logfmt.New(os.Stdout))
	addr := ":" + os.Getenv("PORT")
	http.HandleFunc("/submit", submit)
	http.HandleFunc("/", index)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("error listening: %s", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	views.ExecuteTemplate(w, "index.html", nil)
}

func submit(w http.ResponseWriter, r *http.Request) {
	file, hdr, err := r.FormFile("image")
	if err != nil {
		log.WithError(err).Error("parsing form")
		http.Error(w, "Error parsing form.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	buf := make([]byte, hdr.Size)

	// read file content into buffer
	fReader := bufio.NewReader(file)
	fReader.Read(buf)

	w.Header().Set("Content-Type", "text/html")
	views.ExecuteTemplate(w, "index.html", struct {
		Name         string
		Size         int64
		Type         string
		ImgBase64Str string
	}{
		Name:         hdr.Filename,
		Size:         hdr.Size,
		Type:         hdr.Header.Get("Content-Type"),
		ImgBase64Str: base64.StdEncoding.EncodeToString(buf),
	})
}
