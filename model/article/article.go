package article

import (
	"github.com/dancannon/gorethink"
	"github.com/keimoon/gore"
	"github.com/tamnd/voicewiki/model"
	"github.com/tamnd/voicewiki/model/section"
	"mime/multipart"
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
	Existed    bool               `gorethink:"-" json:"-"`
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
		article.Existed = true
		return article, nil
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
	docs := search(model.Rethink, query)[:1]
	articles := []*Article{}
	for _, doc := range docs {
		article, err := Get(doc.Title)
		if err == nil {
			articles = append(articles, article)
		}
	}
	return articles, nil
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

func (article *Article) Merge() error {
	shard := model.GetShardID(article.Id)
	if article.Existed {
		for _, sectionId := range article.SectionsId {
			sect, err := section.Get(article.Id, sectionId)
			if err != nil {
				return err
			}
			err = sect.Merge(shard)
			if err != nil {
				return err
			}
			article.Sections = append(article.Sections, sect)
		}
	}
	return nil
}

func (article *Article) AddAudio(sectionId string, uploadFile multipart.File) error {
	if article.Existed {
		for _, sectionId := range article.SectionsId {
			sect, err := section.Get(article.Id, sectionId)
			if err != nil {
				return err
			}
			article.Sections = append(article.Sections, sect)
		}
	}
	var sect *section.Section = nil
	for i := range article.Sections {
		if article.Sections[i].Id == sectionId {
			sect = article.Sections[i]
		}
	}
	if sect == nil {
		return model.ErrNotFound
	}
	var err error
	if !article.Existed {
		err = article.CreateFromRaw()
		if err != nil {
			return err
		}
	} else {
		sect, err = section.Get(article.Id, sect.Id)
		if err != nil {
			return err
		}
	}
	return sect.AddAudio(article.Id, uploadFile)
}

func (article *Article) CreateFromRaw() error {
	shard := model.GetShardID(article.Id)
	for _, sect := range article.Sections {
		err := sect.CreateFromRaw(shard)
		if err != nil {
			return err
		}
		article.SectionsId = append(article.SectionsId, sect.Id)
	}
	_, err := gorethink.Table("articles" + shard).Insert(article).RunWrite(model.Rethink)
	if err == nil {
		article.Existed = true
	}
	return err
}

func (raw *ArticleRaw) build() *Article {
	article := &Article{
		Id:    raw.Id,
		Title: raw.Title,
	}
	for i, line := range strings.Split(raw.Content, "\n") {
		article.Sections = append(article.Sections, section.BuildFromText(i, line))
	}
	article.Existed = false
	return article
}
