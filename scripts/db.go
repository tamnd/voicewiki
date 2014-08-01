package main

import (
	"github.com/dancannon/gorethink"
	"strconv"
)

func main() {
	session, _ := gorethink.Connect(gorethink.ConnectOpts{
		Address:  "127.1:28015",
		Database: "wikivoice",
	})
	tables := []string{"article_raws", "articles", "article_versions", "sections", "audios"}
	for _, table := range tables {
		for i := 0; i < 10; i++ {
			gorethink.Db("wikivoice").TableCreate(table + strconv.FormatInt(int64(i), 10)).RunWrite(session)
		}
	}
}
