package main

import (
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/global"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/server"
)

func main() {
	global.Init()
	server.Start()
}
