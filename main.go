package main

import (
	"Mmx/Modules"
	"Mmx/Router"
)

func main() {
	Modules.Config.Init()
	Modules.Global.Init()
	Router.InitRouter()
}
