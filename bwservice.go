package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

const (
	BACKGROUND_IMAGE = "http://bwservice.bwgame.com.cn/images/bk.jpg"
	CLOSE_IMAGE1     = "http://bwservice.bwgame.com.cn/images/close1.gif"
	CLOSE_IMAGE2     = "http://bwservice.bwgame.com.cn/images/close2.gif"
	CLOSE_BUTTON_X   = "293"
	CLOSE_BUTTON_Y   = "4"
)

var (
	NotifyStartTime time.Time
	NotifyEndTime   time.Time
	ClickURL        string
)

func Patrol(rw http.ResponseWriter, req *http.Request) {
	notify := false
	now := time.Now()
	req.ParseForm()
	if req.Form.Get("devmod") == "1" {
		notify = true
	} else {
		s := req.Form.Get("lasttime")
		lasttime, err := strconv.Atoi(s)
		if err != nil || lasttime == 0 {
			if now.After(NotifyStartTime) && now.Before(NotifyEndTime) {
				notify = true
			}
		} else {
			t := time.Unix(int64(lasttime), 0)
			if t.Before(NotifyStartTime) && now.After(NotifyStartTime) && now.Before(NotifyEndTime) {
				notify = true
			}
		}
	}
	result := ""
	result += "bkimage=" + BACKGROUND_IMAGE
	result += "&url=" + ClickURL
	result += "&lasttime=" + strconv.Itoa(int(now.Unix()))
	result += "&clsbtimage1=" + CLOSE_IMAGE1
	result += "&clsbtimage2=" + CLOSE_IMAGE2
	result += "&clsbtx=" + CLOSE_BUTTON_X
	result += "&clsbty=" + CLOSE_BUTTON_Y
	if notify {
		result += "&notify"
	}
	rw.Write([]byte(result))
	remote_addr := req.Header.Get("X-Forwarded-For")
	if remote_addr == "" {
		remote_addr = req.Header.Get("X-Real-IP")
	}
	if remote_addr == "" {
		remote_addr = req.RemoteAddr
	}
	log.Printf("[%s] %s?%s %v\n", remote_addr, req.URL.Path, req.URL.RawQuery, notify)
}

func main() {
	var host = flag.String("host", "", "Server listen host, default 0.0.0.0")
	var port = flag.Int("port", 80, "Server listen port, default 80")
	var url = flag.String("url", "", "Click open url")
	var nst = flag.String("nst", "", "Notify start time, format: YYYY-mm-dd HH:MM")
	var nvt = flag.String("nvt", "", "Notify over time, format as start time")
	var cpu = flag.Int("cpu", 0, "Go runtime max process")
	flag.Parse()
	if *cpu > 0 {
		runtime.GOMAXPROCS(*cpu)
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04", *nst, loc); err != nil {
		panic(err)
	} else {
		NotifyStartTime = t
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04", *nvt, loc); err != nil {
		panic(err)
	} else {
		NotifyEndTime = t
	}
	if *url == "" {
		panic(*url)
	} else {
		ClickURL = *url
	}
	var addr = net.JoinHostPort(*host, strconv.Itoa(*port))
	http.HandleFunc("/patrol", Patrol)
	log.Printf("[%s] (%s) - (%s) => %s\n", addr, NotifyStartTime.String(), NotifyEndTime.String(), ClickURL)
	log.Fatal(http.ListenAndServe(addr, nil))
}
