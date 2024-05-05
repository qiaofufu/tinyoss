package main

import (
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/global"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/meta"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/server"
)

func main() {
	global.Init()
	meta.InitMeta(global.Etcd)
	server.Start()
}
