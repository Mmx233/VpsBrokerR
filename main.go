package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

func loger() {
	// 级别为info warning none
}

func worker() {
	//计数者
	b := sign
	for {
		time.Sleep(time.Duration(1) * time.Second)
		popes[sign].a--
		if popes[sign].a < 0 {
			tprint(popes[sign].b + " 超时")
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

type pp struct {
	a int
	b string
	c string
}
type dpp struct {
	a int64
}

var popes map[string]*pp = make(map[string]*pp)
var dpper map[string]*dpp = make(map[string]*dpp)

//接收get参数
var name string
var sign string

func main() {
	//接收参数
	var port string
	var path string
	var log string
	var url string
	flag.StringVar(&port, "p", "233", "监听端口")
	flag.StringVar(&path, "path", "/", "心跳路径")
	flag.StringVar(&log, "log", "info", "日志基本")
	flag.StringVar(&url, "url", "", "上报URL")
	flag.Parse()

	//检查参数
	if url == "" {
		fmt.Println("上报URL为必填项")
		os.Exit(3)
	}

	//监听
	s := &http.Server{
		Addr: ":" + port,
	}
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		get := r.URL.Query()
		if get["time"] == nil || get["name"] == nil || get["sign"] == nil {
			w.Write([]byte(`{"status":"fail","message":"缺少参数"}`))
		} else {
			name = get["name"][0]
			if get["sign"][0] == "name" {
				sign = name
			} else if get["sign"][0] == "ip" {
				sign = r.RemoteAddr
			} else {
				w.Write([]byte(`{"status":"fail","message":"sign参数有误"}`))
				return
			}
			time, err := strconv.Atoi(get["time"][0])
			if err != nil {
				errer("请求中的time参数有误", err)
				w.Write([]byte(`{"status":"fail","message":"time参数有误"}`))
				return
			}
			if popes[sign] == nil {
				popes[sign] = &pp{time, name, r.RemoteAddr}
				go worker()
				if dpper[sign] != nil {
					tprint("「" + name + "」恢复，历时" + timecount(sign))
					delete(dpper, sign)
				} else {
					tprint("新主机「" + name + "」上线")
				}
				w.Write([]byte(`{"status":"ok","data":"init"}`))
			} else {
				popes[sign].a = time
				popes[sign].b = name
				popes[sign].c = r.RemoteAddr
				w.Write([]byte(`{"status":"ok","data":"continue"}`))
			}
		}
	})
	s.ListenAndServe()
}
