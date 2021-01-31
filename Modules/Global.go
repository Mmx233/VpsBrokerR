package Modules

type VpsInfo struct {
	Time    uint
	Count   uint
	Name    string
	Ip      string
	Backend string
}

type global struct {
	Online  map[string]*VpsInfo
	Offline map[string]int64
}

var Global global

func (a *global) Init() {
	a.Online = make(map[string]*VpsInfo, 0)
	a.Offline = make(map[string]int64, 0)
}
