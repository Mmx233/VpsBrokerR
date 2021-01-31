package Modules

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type worker struct{}

var Worker worker

func (*worker) backReq(sign string) bool {
	//心跳超时，反向请求
	h, err := http.Get(Global.Online[sign].Backend)
	if err != nil {
		Callback.CMD.Error(err)
		return false
	}
	defer h.Body.Close()
	body, err := ioutil.ReadAll(h.Body)
	if err != nil {
		Callback.CMD.Error(err)
		return false
	}
	if strings.TrimSpace(string(body)) == "1" {
		return true
	} else {
		Callback.CMD.Error(errors.New("「" + Global.Online[sign].Name + "」(" +  Global.Online[sign].Ip + ") 反向响应不正确"))
		return false
	}
}

func (a *worker) Do(sign string, vps VpsInfo) {
	vps.Count = vps.Time
	Global.Online[sign] = &vps
	go func() {
		for {
			time.Sleep(time.Second)
			vps.Count--
			if vps.Count == 0 {
				if vps.Backend != "" {
					if a.backReq(sign) {
						vps.Count = vps.Time
						Callback.CMD.Info("「" + vps.Name + "」(" + vps.Ip + ") 心跳超时，反向正常")
						continue
					}
					Callback.Backend.Down(sign)
					//提升精度，持续反向
					for {
						if vps.Count == 0 {
							time.Sleep(time.Second)
							if a.backReq(sign) {
								vps.Count=vps.Time
								Callback.Backend.Up(sign)
								break
							}
						} else {
							break
						}
					}
				}else {//未配置backend
					Callback.Backend.Down(sign)
					return
				}
			}
		}
	}()
}
