package main

import (
	"fmt"
	"github.com/dancannon/gorethink"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ArticleRaw struct {
	Id      string `gorethink:"id"`
	Title   string `gorethink:"title"`
	Content string `gorethink:"content"`
}

func getShardId(slug string) string {
	r := 0
	for _, c := range slug {
		r = r*256 + int(c)
		r = r % 10
	}
	return strconv.FormatInt(int64(r), 10)
}

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

var count int64 = 0

type Title struct {
	Id   int64  `gorethink:"id"`
	Slug string `gorethink:"slug"`
}

type Index struct {
	Id   string  `gorethink:"id"`
	Docs []int64 `gorethink:"docs"`
}

var titles = []*Title{}
var idxes = make(map[string]*Index)

func index(session *gorethink.Session, title string) {
	count++
	slug := strings.ToLower(title)
	titles = append(titles, &Title{Id: count, Slug: slug})
	tokens := tokenizer(title)
	for n := 1; n <= 3; n++ {
		grams := ngram(tokens, n)
		for _, gram := range grams {
			idx := idxes[gram]
			if idx == nil {
				idx = &Index{
					Id:   gram,
					Docs: []int64{count},
				}
				idxes[gram] = idx
			} else {
				idx.Docs = append(idx.Docs, count)
			}
		}
	}
}

func pushIndex(session *gorethink.Session) {
	i := 0
	fmt.Println("Pushing title")
	for i < len(titles) {
		var chunk []*Title
		if i+100 < len(titles) {
			chunk = titles[i : i+100]
		} else {
			chunk = titles[i:len(titles)]
		}
		gorethink.Table("titles").Insert(chunk, gorethink.InsertOpts{Durability: "soft", Upsert: true}).RunWrite(session)
		fmt.Println("Title %d / %d\n", i, len(titles))
		i += 100
	}
	fmt.Println("Pushing index")
	idxChunk := []*Index{}
	i = 0
	for _, idx := range idxes {
		i++
		if len(idxChunk) < 100 {
			idxChunk = append(idxChunk, idx)
			continue
		}
		gorethink.Table("indexes").Insert(idxChunk, gorethink.InsertOpts{Durability: "soft", Upsert: true}).RunWrite(session)
		fmt.Println("Index %d / %d\n", i, len(idxes))
		idxChunk = []*Index{}
	}
}

func importFile(session *gorethink.Session, file string) {
	b, _ := ioutil.ReadFile(file)
	text := string(b)
	var title string
	titles := []string{}
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "<doc") {
			index := strings.Index(line, "title=")
			if index < 0 {
				continue
			}
			title = line[index+7 : len(line)-2]
			titles = append(titles, title)
		}
	}
	for _, title := range titles {
		index(session, title)
	}
}

func importDir(session *gorethink.Session, dir string) {
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		importFile(session, filepath.Join(dir, file.Name()))
		fmt.Printf("Imported: `%s`\n", filepath.Join(dir, file.Name()))
	}
}

func main() {
	session, _ := gorethink.Connect(gorethink.ConnectOpts{
		Address:  "127.1:28015",
		Database: "wikivoice",
	})
	dir := "extracted"
	subdirs, _ := ioutil.ReadDir(dir)
	for _, subdir := range subdirs {
		importDir(session, filepath.Join(dir, subdir.Name()))
	}
	pushIndex(session)
}
