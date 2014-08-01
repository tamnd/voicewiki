package middleware

import (
	"net/http"
	"os"
	"path/filepath"
)

func Route(pattern string, handler http.HandlerFunc) {
	if pattern == "/" {
		homeHandler = handler
	} else {
		http.HandleFunc(pattern, handler)
	}
}

var homeHandler http.HandlerFunc

func homeDefaultHandler(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

func InitRouter() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			homeHandler(w, req)
			return
		}
		filename := filepath.Join(Config.App.WebRoot, req.URL.Path)
		info, err := os.Stat(filename)
		if err != nil || info.IsDir() {
			homeHandler(w, req)
			return
		}
		http.ServeFile(w, req, filename)
	})
}

func Run() {
	http.ListenAndServe(Config.App.Bind, nil)
}
