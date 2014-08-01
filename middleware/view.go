package middleware

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
)

type Data map[string]interface{}

var (
	ViewPath   = "view"
	LayoutFile = "layout/default.html"
)

var templateCache = make(map[string]string)

func loadTemplateContent(filename string) (string, error) {
	var content string
	if Config.App.Env != "dev" {
		content = templateCache[filename]
	}
	if len(content) == 0 {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return "", err
		}
		content = string(b)
	}
	if Config.App.Env != "dev" {
		templateCache[filename] = content
	}
	return content, nil
}

func RenderView(w io.Writer, filename string, data Data) error {
	content, err := loadTemplateContent(LayoutFile)
	if err != nil {
		return err
	}
	t := template.New("template")
	t.Funcs(template.FuncMap{
		"block": blockFunc,
	})
	_, err = t.Parse(content)
	if err != nil {
		return err
	}
	data["Content"] = filename
	return t.Execute(w, data)
}

func blockFunc(filename string, data interface{}) (string, error) {
	content, err := loadTemplateContent(filepath.Join(ViewPath, filename))
	if err != nil {
		return "", err
	}
	t := template.New("block")
	_, err = t.Parse(content)
	if err != nil {
		return "", err
	}
	buffer := &bytes.Buffer{}
	err = t.Execute(buffer, data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
