package main

import (
	"fmt"
	"github.com/tamnd/voicewiki/handler"
	"github.com/tamnd/voicewiki/middleware"
	"os"
)

func main() {
	var configFile string = "app.config"
	if len(os.Args) >= 2 {
		configFile = os.Args[1]
	}
	err := middleware.LoadConfig(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	middleware.InitRouter()
	middleware.Route("/", handler.Home)
	middleware.Route("/list", handler.List)

	middleware.Run()
}
