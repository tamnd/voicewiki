package section

import (
	"github.com/dancannon/gorethink"
	"github.com/tamnd/voicewiki/model"
	"github.com/tamnd/voicewiki/model/audio"
	"mime/multipart"
	"strconv"
)

type Section struct {
	Id       string         `gorethink:"-" json:"id"`
	Content  string         `gorethink:"content" json:"content"`
	AudiosId []string       `gorethink:"audios" json:"-"`
	Audios   []*audio.Audio `gorethink:"-" json:"audios"`
}

func Get(articleId string, id string) (*Section, error) {
	shard := model.GetShardID(articleId)
	rows, err := gorethink.Table("sections" + shard).Get(id).Run(model.Rethink)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.IsNil() {
		return nil, model.ErrNotFound
	}
	s := &Section{}
	err = rows.One(s)
	if err != nil {
		return nil, err
	}
	s.Id = id
	return s, nil
}

func (s *Section) Merge(shard string) error {
	for _, audioId := range s.AudiosId {
		audio, err := audio.Get(shard, audioId)
		if err != nil {
			return err
		}
		s.Audios = append(s.Audios, audio)
	}
	return nil
}

func BuildFromText(pos int, text string) *Section {
	return &Section{
		Id:      strconv.FormatInt(int64(pos), 10),
		Content: text,
	}
}

func (s *Section) CreateFromRaw(shard string) error {
	result, err := gorethink.Table("sections" + shard).Insert(s).RunWrite(model.Rethink)
	if err != nil {
		return err
	}
	s.Id = result.GeneratedKeys[0]
	return nil
}

func (s *Section) AddAudio(articleId string, uploadFile multipart.File) error {
	shard := model.GetShardID(articleId)
	audio, err := audio.Create(articleId, s.Id, uploadFile)
	if err != nil {
		return err
	}
	s.AudiosId = append([]string{audio.Id}, s.AudiosId...)
	_, err = gorethink.Table("sections" + shard).Get(s.Id).Update(s).RunWrite(model.Rethink)
	return err
}
