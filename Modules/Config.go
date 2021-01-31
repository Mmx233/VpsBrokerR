package Modules

import "flag"

type config struct {
	Port string
	Path string
	SUrl string
	Logg string
}

var Config config

func (a *config) Init() {
	flag.StringVar(&a.Port, "p", "233", "监听端口")
	flag.StringVar(&a.Path, "path", "/", "心跳路径")
	flag.StringVar(&a.SUrl, "url", "", "上报URL")
	flag.StringVar(&a.Logg, "log", "false", "日志开关")
	flag.Parse()
}
