package middleware

import (
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func Route(pattern string, handler interface{}) {
	goji.Handle(pattern, handler)
}

func NotFound(w http.ResponseWriter) {	
	http.Error(w, "Not Found", http.StatusNotFound)
}

func Fatal(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

type ServerMux struct {
}

func (s *ServerMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	filename := filepath.Join(Config.App.WebRoot, req.URL.Path)
	info, err := os.Stat(filename)
	if err == nil && !info.IsDir() {
		http.ServeFile(w, req, filename)
		return
	}
	goji.DefaultMux.ServeHTTP(w, req)
}

var DefaultServerMux = &ServerMux{}

func Run() {
	http.Handle("/", DefaultServerMux)
	listener, err := net.Listen("tcp", Config.App.Bind)
	if err != nil {
		panic(err)
	}
	err = graceful.Serve(listener, DefaultServerMux)
	if err != nil {
		panic(err)
	}
	graceful.Wait()
}
