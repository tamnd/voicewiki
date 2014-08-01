package article

import (
	"github.com/dancannon/gorethink"
	"github.com/keimoon/gore"
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
	id = strings.Replace(strings.ToLower(id), "_", " ", -1)
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
	article, err := Get(query)
	if err != nil {
		return nil, err
	}
	return []*Article{article}, nil
}

func List() ([]*Article, error) {
	pool := model.Redis["list"]
	conn, err := pool.Acquire()
	if conn == nil {
		return nil, nil
	}
	defer pool.Release(conn)
	result, err := gore.NewCommand("ZREVRANGE", "articles", 0, 10).Run(conn)
	if err != nil {
		return nil, err
	}
	idList := []string{}
	err = result.Slice(&idList)
	articles := []*Article{}
	for _, id := range idList {
		article, err := Get(id)
		if err == nil {
			articles = append(articles, article)
		}
	}
	return articles, nil
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
