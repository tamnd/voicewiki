package api

import (
	"github.com/tamnd/voicewiki/middleware"
	"github.com/tamnd/voicewiki/model"
	"github.com/tamnd/voicewiki/model/article"
	"github.com/zenazn/goji/web"
	"net/http"
)

func GetArticle(c web.C, w http.ResponseWriter, r *http.Request) {
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
	middleware.RenderJSON(w, article)
}

func SearchArticle(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()["q"]
	if len(queries) == 0 {
		middleware.NotFound(w)
		return
	}
	query := queries[0]
	articles, _ := article.Search(query)
	middleware.RenderJSON(w, articles)
}

func ListArticle(w http.ResponseWriter, r *http.Request) {
	articles, _ := article.List()
	middleware.RenderJSON(w, articles)
}
