package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func tprint(a string) {
	fmt.Println(time.Now().Format("[2006-01-02 15:04:05] ") + a)
}

func errer(a string, err error) {
	fmt.Println(err)
	tprint(a)
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func backe(b string) bool {
	//心跳超时，反向请求
	_, err := http.Get(popes[b].d)
	if err == nil {
		return false
	} else {
		return true
	}
}

func worker() {
	//计数者
	b := sign
	for {
		time.Sleep(time.Duration(1) * time.Second)
		popes[b].a--
		if popes[b].a < 0 {
			if backe(b) {
				popes[b].a = popes[b].aa
				continue
			}
			msg := "「" + popes[b].b + "」(" + popes[b].c + ") 宕机"
			tprint(msg)
			go urler(msg, popes[b].b, 0, "down")
			delete(popes, b)
			dpper[b] = &dpp{time.Now().Unix()}
			break
		}
	}
}

func timecount(a string) string {
	temp := time.Now().Unix() - dpper[sign].a
	a = ""
	if temp > (24 * 60 * 60) {
		temp2 := temp % (24 * 60 * 60)
		a += strconv.FormatInt((temp-temp2)/(24*60*60), 10) + "天"
		temp = temp2
	}
	if temp > (60 * 60) {
		temp2 := temp % (60 * 60)
		a += strconv.FormatInt((temp-temp2)/(60*60), 10) + "时"
		temp = temp2
	}
	if temp > 60 {
		temp2 := temp % 60
		a += strconv.FormatInt((temp-temp2)/60, 10) + "分"
		temp = temp2
	}
	if temp > 0 {
		a += strconv.FormatInt(temp, 10) + "秒"
	}
	return a
}

func urler(msg string, name string, time int64, mtype string) {
	_, err := http.Get(surl + "?msg=" + url.QueryEscape(msg) + "&name=" + url.QueryEscape(name) + "&time=" + strconv.FormatInt(time, 10) + "&type=" + mtype)
	if err != nil {
		errer("上报URL请求失败", err)
	}
}

type pp struct {
	a  int
	b  string
	c  string
	aa int
	d  string
}
type dpp struct {
	a int64
}

var popes map[string]*pp = make(map[string]*pp)
var dpper map[string]*dpp = make(map[string]*dpp)

//接收get参数
var name string
var sign string
var surl string

var ip string

func main() {
	//接收参数
	var port string
	var path string
	flag.StringVar(&port, "p", "233", "监听端口")
	flag.StringVar(&path, "path", "/", "心跳路径")
	flag.StringVar(&surl, "url", "", "上报URL")
	flag.Parse()

	//检查参数
	if surl == "" {
		fmt.Println("上报URL为必填项")
		os.Exit(3)
	}

	//监听
	s := &http.Server{
		Addr: ":" + port,
	}
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		get := r.URL.Query()
		if get["time"] == nil || get["name"] == nil || get["sign"] == nil || get["backend"] == nil {
			w.Write([]byte(`{"status":"fail","message":"缺少参数"}`))
		} else {
			ip, _, _ = net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
			name = get["name"][0]
			if get["sign"][0] == "name" {
				sign = name
			} else if get["sign"][0] == "ip" {
				sign = ip
			} else {
				w.Write([]byte(`{"status":"fail","message":"sign参数有误"}`))
				return
			}
			ttime, err := strconv.Atoi(get["time"][0])
			if err != nil {
				errer("请求中的time参数有误", err)
				w.Write([]byte(`{"status":"fail","message":"time参数有误"}`))
				return
			}
			if popes[sign] == nil {
				popes[sign] = &pp{ttime, name, ip, ttime, get["backend"][0]}
				go worker()
				if dpper[sign] != nil {
					msg := "「" + name + "」(" + ip + ") 恢复，历时" + timecount(sign)
					tprint(msg)
					go urler(msg, name, time.Now().Unix()-dpper[sign].a, "up")
					delete(dpper, sign)
				} else {
					msg := "新主机「" + name + "」(" + ip + ") 上线"
					tprint(msg)
					go urler(msg, name, 0, "new")
				}
				w.Write([]byte(`{"status":"ok","data":"init"}`))
			} else {
				popes[sign].a = ttime
				popes[sign].aa = ttime
				popes[sign].b = name
				popes[sign].c = ip
				w.Write([]byte(`{"status":"ok","data":"continue"}`))
			}
		}
	})
	s.ListenAndServe()
}
