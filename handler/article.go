package handler

import (
        "github.com/tamnd/voicewiki/middleware"
	"net/http"
)

func Article(w http.ResponseWriter, r *http.Request) {
	data := make(middleware.Data)
	middleware.RenderView(w, "article.html", data)
}
