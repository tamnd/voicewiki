package handler

import (
	"github.com/tamnd/voicewiki/middleware"
	"github.com/tamnd/voicewiki/model/article"
	"net/http"
)

func List(w http.ResponseWriter, r *http.Request) {
	data := make(middleware.Data)
	articles, err := article.List()
	if err != nil {
		middleware.Fatal(w)
		return
	}
	data["Title"] = "Recent articles"
	data["Articles"] = articles
	middleware.RenderView(w, "list.html", data)
}
