package section

import (
	"github.com/tamnd/voicewiki/model/audio"
)

type Section struct {
	Content  string         `gorethink:"content" json:"content"`
	AudiosId []string       `gorethink:"audios" json:"-"`
	Audios   []*audio.Audio `gorethink:"-" json:"audios"`
}

func BuildFromText(text string) *Section {
	return &Section{
		Content: text,
	}
}
