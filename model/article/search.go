package article

import (
	"github.com/dancannon/gorethink"
	"github.com/tamnd/voicewiki/middleware"
	"math"
	"regexp"
	"sort"
	"strings"
)

var replaceRe = regexp.MustCompile("\\(|\\)|\"|'|_|\\-")
var splitRe = regexp.MustCompile("\\s+")

func tokenizer(title string) []string {
	title = replaceRe.ReplaceAllString(strings.ToLower(title), "")
	return splitRe.Split(title, -1)
}

func ngram(tokens []string, n int) []string {
	i := 0
	grams := []string{}
	limit := len(tokens) - (n - 1)
	for i < limit {
		grams = append(grams, strings.Join(tokens[i:i+n], " "))
		i++
	}
	return grams
}

type Document struct {
	Title string
	Score float64
}

type DocumentList []*Document

func (a DocumentList) Len() int {
	return len(a)
}

func (a DocumentList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a DocumentList) Less(i, j int) bool {
	return a[i].Score > a[j].Score
}

type Title struct {
	Id   int64  `gorethink:"id"`
	Slug string `gorethink:"slug"`
}

type Index struct {
	Id   string  `gorethink:"id"`
	Docs []int64 `gorethink:"docs"`
}

func search(session *gorethink.Session, query string) DocumentList {
	numDoc := middleware.Config.App.NumDoc
	tokens := tokenizer(query)
	docs := make(map[int64]*Document)
	for n := 1; n <= 3; n++ {
		grams := ngram(tokens, n)
		for _, gram := range grams {
			rows, err := gorethink.Table("indexes").Get(gram).Run(session)
			if err != nil {
				continue
			}
			if rows.IsNil() {
				rows.Close()
				continue
			}
			idx := &Index{}
			rows.One(idx)
			for _, docId := range idx.Docs {
				doc := docs[docId]
				if doc == nil {
					rows, err = gorethink.Table("titles").Get(docId).Run(session)
					if err != nil {
						continue
					}
					if rows.IsNil() {
						rows.Close()
						continue
					}
					title := &Title{}
					rows.One(title)
					doc = &Document{Score: 0, Title: title.Slug}
				}
				tf := float64(n) / float64(len(doc.Title))
				idf := math.Log(float64(numDoc) / float64(len(idx.Docs)))
				doc.Score += tf * idf * float64(n)
				docs[docId] = doc
			}
		}
	}
	docList := DocumentList{}
	for _, doc := range docs {
		docList = append(docList, doc)
	}
	sort.Sort(docList)
	return docList
}
