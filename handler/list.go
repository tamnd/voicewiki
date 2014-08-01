package handler

import (
        "github.com/tamnd/voicewiki/middleware"
	"net/http"
)

func List(w http.ResponseWriter, r *http.Request) {
	data := make(middleware.Data)
	middleware.RenderView(w, "list.html", data)
}
