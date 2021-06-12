package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/danielhoward314/microservices-test/product-images/files"
	"github.com/gorilla/mux"
)

type Images struct {
	l *log.Logger
	p string
	s files.Storage
}

func NewImages(l *log.Logger, p string, s files.Storage) *Images {
	return &Images{l, p, s}
}

func (i *Images) Upload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fn := vars["filename"]
	i.l.Printf("Uploading file with filename=%v", fn)
	i.saveFile(fn, w, r.Body)
}

func (i *Images) GetImage(w http.ResponseWriter, r *http.Request) {
	imagesDir, err := filepath.Abs(i.p)
	if err != nil {
		i.l.Print("unable to get absolute path of imagestore")
		http.Error(w, "unable to get absolute path of imagestore", http.StatusInternalServerError)
	}
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]
	fullPath := filepath.Join(imagesDir, id, fn)
	// check whether a file exists at the given path
	_, err = os.Stat(fullPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, fullPath)
	// http.FileServer(http.Dir(fullPath)).ServeHTTP(w, r)
}

func (i *Images) saveFile(fn string, w http.ResponseWriter, r io.ReadCloser) {
	err := i.s.Save(fn, r)
	if err != nil {
		i.l.Print("error in saveFile")
		http.Error(w, "unable to save file", http.StatusInternalServerError)
	}
}
