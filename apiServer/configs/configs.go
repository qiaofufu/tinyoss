package configs

import "flag"

type Configs struct {
	Server
	Etcd
}

type Server struct {
	Ip            string
	Port          int32
	BaseDir       string
	LocateTimeout int64
}

type Etcd struct {
	Endpoints   []string
	DialTimeout int
}

var (
	port int
)

func NewConfigs() *Configs {
	return &Configs{}
}

func (c *Configs) Load() {
	flag.IntVar(&port, "p", 8000, "server port")
	flag.Parse()
	c.Server.Ip = "localhost"
	c.Server.Port = int32(port)
	c.Server.BaseDir = "/tmp/tinyoss"
	c.Server.LocateTimeout = 2
	c.Etcd.Endpoints = []string{"http://162.14.115.114:2379"}
	c.Etcd.DialTimeout = 5
}
