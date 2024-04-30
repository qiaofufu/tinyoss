package main

import (
	"github.com/qiaofufu/tinyoss_kernal/dataServer/internal/global"
	"github.com/qiaofufu/tinyoss_kernal/dataServer/internal/server"
)

func main() {
	global.Init()
	server.Start()
}
