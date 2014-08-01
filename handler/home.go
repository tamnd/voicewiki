package handler

import (
	"github.com/tamnd/voicewiki/middleware"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	data := make(middleware.Data)
	middleware.RenderView(w, "home.html", data)
}
