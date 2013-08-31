package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	BACKGROUND_IMAGE = "http://conf.bwgame.org/static/res/sbk.jpg"
	CLOSE_IMAGE1     = "http://conf.bwgame.org/static/res/close1.gif"
	CLOSE_IMAGE2     = "http://conf.bwgame.org/static/res/close2.gif"
	CLICK_URL        = "http://shenzuo.bwgame.com.cn"
	CLOSE_BUTTON_X   = "293"
	CLOSE_BUTTON_Y   = "4"
)

var (
	NotifyStartTime time.Time
	NotifyEndTime   time.Time
)

func init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	NotifyStartTime = time.Date(2013, time.August, 31, 16, 1, 0, 0, loc)
	NotifyEndTime = time.Date(2013, time.September, 31, 16, 1, 0, 0, loc)
}

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
	result += "&url=" + CLICK_URL
	result += "&lasttime=" + strconv.Itoa(int(now.Unix()))
	result += "&clsbtimage1=" + CLOSE_IMAGE1
	result += "&clsbtimage2=" + CLOSE_IMAGE2
	result += "&clsbtx=" + CLOSE_BUTTON_X
	result += "&clsbty=" + CLOSE_BUTTON_Y
	if notify {
		result += "&notify"
	}
	rw.Write([]byte(result))
	log.Printf("[%s] %s?%s %v\n", req.RemoteAddr, req.URL.Path, req.URL.RawQuery, notify)
}

func main() {
	var host = flag.String("host", "", "Server listen host, default 0.0.0.0")
	var port = flag.Int("port", 80, "Server listen port, default 80")
	flag.Parse()
	var addr = net.JoinHostPort(*host, strconv.Itoa(*port))
	http.HandleFunc("/patrol", Patrol)
	log.Println(addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
