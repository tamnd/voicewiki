package audio

import (
	"github.com/dancannon/gorethink"
	"github.com/tamnd/voicewiki/middleware"
	"github.com/tamnd/voicewiki/model"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Audio struct {
	Id     string `gorethink:"-" json:"id"`
	Server string `gorethink:"server" json:"-"`
	Path   string `gorethink:"path" json:"-"`
	URL    string `gorethink:"-" json:"url"`
}

func Get(shard string, id string) (*Audio, error) {
        rows, err := gorethink.Table("audios" + shard).Get(id).Run(model.Rethink)
	if err != nil {
                return nil, err
	}
	defer rows.Close()
	if rows.IsNil() {
		return nil, model.ErrNotFound
	}
	audio := &Audio{}
	err = rows.One(audio)
	if err != nil {
		return nil, err
	}
	audio.Id = id
	audio.URL = audio.Server + "/" + audio.Path
        return audio, nil
}

func Create(articleId string, sectionId string, uploadFile multipart.File) (*Audio, error) {
	folder := filepath.Join(articleId, sectionId)
	fullFolder := filepath.Join(middleware.Config.Upload.Root, folder)
	path := filepath.Join(folder, strconv.FormatInt(time.Now().UnixNano(), 10) + ".mp3")
	fullpath := filepath.Join(middleware.Config.Upload.Root, path)
	_, err := os.Stat(fullFolder)
	if err != nil {
		err = os.MkdirAll(fullFolder, 0755)
		if err != nil {
			return nil, err
		}
	}
	newFile, err := os.Create(fullpath)
	if err != nil {
		return nil, err
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, uploadFile)
	if err != nil {
		return nil, err
	}
	audio := &Audio{
		Server: middleware.Config.Upload.Server,
		Path:   path,
	}
	shard := model.GetShardID(articleId)
	result, err := gorethink.Table("audios" + shard).Insert(audio).RunWrite(model.Rethink)
	if err == nil {
		audio.Id = result.GeneratedKeys[0]
		audio.URL = audio.Server + "/" + audio.Path
	}
	return audio, nil
}
