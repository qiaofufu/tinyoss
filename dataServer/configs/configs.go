package configs

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

func NewConfigs() *Configs {
	return &Configs{}
}

func (c *Configs) Load() {
	c.Server.Ip = "localhost"
	c.Server.Port = 8000
	c.Server.BaseDir = "/tmp/tinyoss"
	c.Server.LocateTimeout = 2
	c.Etcd.Endpoints = []string{"http://162.14.115.114:2379"}
	c.Etcd.DialTimeout = 5
}
