package Modules

import (
	"io/ioutil"
	"os"
)

type file struct{}

var File file

func (*file) IsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func (*file) Log(name string, content []byte) {
	name += ".txt"
	err := ioutil.WriteFile(name, content, 0666)
	if err != nil {
		Callback.CMD.Error(err)
	} else {
		Callback.CMD.Info("「已写入日志」")
	}
}
