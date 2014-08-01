package model

import (
	"github.com/dancannon/gorethink"
	"github.com/keimoon/gore"
	"github.com/tamnd/voicewiki/middleware"
)

var Rethink *gorethink.Session
var Redis = make(map[string]*gore.Pool)

func Init() {
	session, err := gorethink.Connect(gorethink.ConnectOpts{
		Address:   middleware.Config.Rethink.Address,
		Database:  middleware.Config.Rethink.Database,
		MaxIdle:   middleware.Config.Rethink.MaxIdle,
		MaxActive: middleware.Config.Rethink.MaxActive,
	})
	if err != nil {
		panic(err)
	}
	Rethink = session
	for name, instance := range middleware.Config.Redis {
		pool := &gore.Pool{}
		err = pool.Dial(instance.Address)
		if err != nil {
			panic(err)
		}
		Redis[name] = pool
	}
}
