package article

import (
	"github.com/dancannon/gorethink"
	"github.com/tamnd/voicewiki/model"
	"github.com/tamnd/voicewiki/model/section"
	"strings"
)

type ArticleRaw struct {
	Id      string `gorethink:"id"`
	Title   string `gorethink:"title"`
	Content string `gorethink:"content"`
}

type Article struct {
	Id         string             `gorethink:"id" json:"id"`
	Title      string             `gorethink:"title" json:"title"`
	SectionsId []string           `gorethink:"sections" json:"-"`
	Sections   []*section.Section `gorethink:"-" json:"sections"`
}

func Get(id string) (*Article, error) {
	shard := model.GetShardID(id)
	rows, err := gorethink.Table("articles" + shard).Get(id).Run(model.Rethink)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.IsNil() {
		article := &Article{}
		err = rows.One(article)
		if err != nil {
			return nil, err
		}
		err = article.merge()
		return article, err
	}
	return getRaw(id)
}

func getRaw(id string) (*Article, error) {
	shard := model.GetShardID(id)
	rows, err := gorethink.Table("article_raws" + shard).Get(id).Run(model.Rethink)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.IsNil() {
		return nil, model.ErrNotFound
	}
	articleRaw := &ArticleRaw{}
	err = rows.One(articleRaw)
	if err != nil {
		return nil, err
	}
	return articleRaw.build(), nil
}

func Search(query string) ([]*Article, error) {
	return nil, nil
}

func (article *Article) merge() error {
	return nil
}

func (raw *ArticleRaw) build() *Article {
	article := &Article{
		Id:    raw.Id,
		Title: raw.Title,
	}
	for _, line := range strings.Split(raw.Content, "\n") {
		article.Sections = append(article.Sections, section.BuildFromText(line))
	}
	return article
}
