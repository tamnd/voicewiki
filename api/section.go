package api

import (
	"github.com/tamnd/voicewiki/middleware"
	"github.com/tamnd/voicewiki/model"
	"github.com/tamnd/voicewiki/model/article"
	"github.com/zenazn/goji/web"
	"net/http"
	"fmt"
)

func ReadSection(c web.C, w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024 * 1024)
	if err != nil {
		middleware.Fatal(w)
		return
	}
	articleId := r.FormValue("article_id")
	sectionId := r.FormValue("section_id")
	article, err := article.Get(articleId)
	if err != nil {
		if err == model.ErrNotFound {
			middleware.NotFound(w)
		} else {
			middleware.Fatal(w)
		}
		return
	}
	uploadFile, _, err := r.FormFile("audio")
	if err != nil {
		middleware.BadRequest(w)
		return
	}
	err = article.AddAudio(sectionId, uploadFile)
	if err != nil {
		if err == model.ErrNotFound {
			middleware.NotFound(w)
		} else {
			fmt.Println(err)
			middleware.Fatal(w)
		}
		return
	}
	middleware.RenderJSON(w, &struct {
		Code int `json:"code"`
	}{
		Code: 0,
	})
}
