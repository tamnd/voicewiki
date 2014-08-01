package handler

import (
	"github.com/tamnd/voicewiki/middleware"
	"github.com/tamnd/voicewiki/model"
	"github.com/tamnd/voicewiki/model/article"
	"github.com/zenazn/goji/web"
	"net/http"
)

func Article(c web.C, w http.ResponseWriter, r *http.Request) {
	data := make(middleware.Data)
	id := c.URLParams["id"]
	article, err := article.Get(id)
	if err != nil {
		if err == model.ErrNotFound {
			middleware.NotFound(w)
		} else {
			middleware.Fatal(w)
		}
		return
	}
	data["Title"] = article.Title
	data["Article"] = article
	middleware.RenderView(w, "article.html", data)
}
