package Modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type callback struct {
	CMD     cmd
	Request request
	Backend backend
}
type cmd struct{}     //控制台
type request struct{} //请求
type backend struct{} //上报
var Callback callback

func (*cmd) do(msg string, color int) {
	fmt.Printf(time.Now().Format("[2006-01-02 15:04:05] ")+"\033[1;%v;40m%s\033[0m\n", color, msg)
}

func (a *cmd) Error(err error) {
	a.do(err.Error(), 31)
}

func (a *cmd) Info(msg string) {
	a.do(msg, 34)
}

type info struct { //响应体
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func (*request) Error(c *gin.Context, err error) {
	c.AsciiJSON(500, info{
		"error",
		err.Error(),
	})
}

func (*request) Info(c *gin.Context, msg string) {
	c.AsciiJSON(200, info{
		"ok",
		msg,
	})
}

func (*backend) do(params map[string]interface{}) error {
	URL, err := url.Parse(Config.SUrl)
	if err != nil {
		return err
	}
	p := url.Values{}
	for k, v := range params {
		p.Set(k, fmt.Sprintf("%v", v))
	}
	URL.RawQuery = p.Encode()
	_, err = http.Get(URL.String())
	if err != nil {
		return err
	}
	return nil
}

func (a *backend) doAction(msg string, name string, time int64, Type string) {
	a.do(map[string]interface{}{
		"msg":  msg,
		"name": name,
		"time": time,
		"type": Type,
	})
}

func (a *backend) Down(sign string) {
	Global.Offline[sign] = time.Now().Unix()
	msg := "「" + Global.Online[sign].Name + "」(" + Global.Online[sign].Ip + ") 宕机"
	Callback.CMD.Error(errors.New(msg))
	a.doAction(msg, Global.Online[sign].Name, 0, "down")
}

func (a *backend) New(sign string) {
	msg := "新主机「" + Global.Online[sign].Name + "」(" + Global.Online[sign].Ip + ") 上线"
	Callback.CMD.Info(msg)
	a.doAction(msg, Global.Online[sign].Name, 0, "new")
}

func (a *backend) Up(sign string) {
	if Config.Logg == "true" { //记录日志
		path := "log/" + time.Now().Format("2006-01") + "/" + sign
		if !File.IsExist(path) {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				Callback.CMD.Error(err)
			}
		}
		t := strconv.FormatInt(time.Now().Unix(), 10)
		data := map[string]string{
			"name":  Global.Online[sign].Name,
			"ip":    Global.Online[sign].Ip,
			"break": strconv.FormatInt(Global.Offline[sign], 10),
			"back":  t,
		}
		bytes, err := json.Marshal(&data)
		if err != nil {
			Callback.CMD.Error(err)
		} else {
			File.Log(path+"/"+t, bytes)
		}
	}
	temp := time.Now().Unix() - Global.Offline[sign]
	delete(Global.Offline, sign)
	m := ""
	if temp > (24 * 60 * 60) {
		temp2 := temp % (24 * 60 * 60)
		m += strconv.FormatInt((temp-temp2)/(24*60*60), 10) + "天"
		temp = temp2
	}
	if temp > (60 * 60) {
		temp2 := temp % (60 * 60)
		m += strconv.FormatInt((temp-temp2)/(60*60), 10) + "时"
		temp = temp2
	}
	if temp > 60 {
		temp2 := temp % 60
		m += strconv.FormatInt((temp-temp2)/60, 10) + "分"
		temp = temp2
	}
	if temp > 0 {
		m += strconv.FormatInt(temp, 10) + "秒"
	}
	msg := "「" + Global.Online[sign].Name + "」(" + Global.Online[sign].Ip + ") 恢复，历时" + m
	Callback.CMD.Info(msg)
	a.doAction(msg, Global.Online[sign].Name, time.Now().Unix() - Global.Offline[sign], "up")
}
