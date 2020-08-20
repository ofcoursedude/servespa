package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/handlers"
)

func main() {
	addr := flag.String("listen", ":4200", "address to listen on, can be either with IP or without")
	path := flag.String("path", "dist", "location of the files to serve")
	index := flag.String("index", "index.html", "name of the main (index) file")
	log := flag.String("log", "true", "whether to print log to stdout")
	flag.Parse()
	if os.Args[0] == "help" {
		flag.PrintDefaults()
		os.Exit(0)
	}
	addressItems := strings.Split(*addr, ":")
	if len(addressItems) == 1 {

	}
	var h http.Handler
	spaServer := NewSpaServer(*path, *index)
	if *log == "true" {
		h = handlers.CombinedLoggingHandler(os.Stdout, spaServer)
	} else {
		h = spaServer
	}
	server := &http.Server{
		Addr:    *addr,
		Handler: h,
	}

	server.ListenAndServe()
}

type SpaServer struct {
	Path          string
	IndexFileName string
	fullIndexFile string
}

func NewSpaServer(path string, indexFileName string) *SpaServer {
	return &SpaServer{Path: path, IndexFileName: indexFileName, fullIndexFile: filepath.Join(path, indexFileName)}
}

func (s SpaServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	file := filepath.Join(s.Path, request.URL.Path)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// accommodate for SPA routing, send index.html or whatever
		http.ServeFile(writer, request, s.fullIndexFile)
	} else if err != nil {
		// something's wrong
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	// send file
	http.FileServer(http.Dir(s.Path)).ServeHTTP(writer, request)
}
