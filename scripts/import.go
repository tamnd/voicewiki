package main

import (
	"fmt"
	"github.com/dancannon/gorethink"
	"io/ioutil"
	"path/filepath"
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

func importFile(session *gorethink.Session, file string) {
	b, _ := ioutil.ReadFile(file)
	text := string(b)
	var title string
	var nextDoc bool
	var sections []string
	articles := make(map[string][]*ArticleRaw)
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "<doc") {
			index := strings.Index(line, "title=")
			if index < 0 {
				nextDoc = true
				continue
			}
			nextDoc = false
			title = line[index+7 : len(line)-2]
		} else if line != "</doc>" {
			if nextDoc {
				continue
			}
			line = strings.TrimSpace(line)
			if len(line) <= 0 {
				continue
			}
			sections = append(sections, line)
		} else {
			if nextDoc {
				continue
			}
			article := &ArticleRaw{
				Id:      strings.ToLower(title),
				Title:   title,
				Content: strings.Join(sections[1:], "\n"),
			}
			articles["article_raws"+getShardId(article.Id)] = append(articles["article_raws"+getShardId(article.Id)], article)
			sections = []string{}
		}
	}
	for table, shard := range articles {
		gorethink.Table(table).Insert(shard, gorethink.InsertOpts{Durability: "soft", Upsert: true}).RunWrite(session)
	}
}

func importDir(dir string) {
	session, _ := gorethink.Connect(gorethink.ConnectOpts{
		Address:  "127.1:28015",
		Database: "wikivoice",
	})
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		importFile(session, filepath.Join(dir, file.Name()))
		fmt.Printf("Imported: `%s`\n", filepath.Join(dir, file.Name()))
	}
}

func main() {
	dir := "extracted"
	subdirs, _ := ioutil.ReadDir(dir)
	for _, subdir := range subdirs {
		importDir(filepath.Join(dir, subdir.Name()))
	}
}
