package global

import (
	"github.com/qiaofufu/tinyoss_kernal/dataServer/configs"
	"log"
	"os"
)

var (
	Cfg *configs.Configs
)

func initConfig() {
	Cfg = configs.NewConfigs()
	Cfg.Load()

	if _, err := os.Stat(Cfg.Server.BaseDir); os.IsNotExist(err) {
		e := os.MkdirAll(Cfg.Server.BaseDir, os.ModePerm)
		if e != nil {
			log.Println(e)
			panic(e)
		}
		e = os.MkdirAll(Cfg.Server.BaseDir+"/objects", os.ModePerm)
		if e != nil {
			log.Println(e)
			panic(e)
		}
	}
}
