package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"time"

	//"github.com/astaxie/beego"
	"liveplay/configure"

	"liveplay/protocol/rtmp"
)

var (
	rtmpAddr       = flag.String("rtmp-addr", ":1935", "RTMP server listen address")
	configfilename = flag.String("cfgfile", "livego.cfg", "live configure filename")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	flag.Parse()
}

func startRtmp(stream *rtmp.RtmpStream) {
	rtmpListen, err := net.Listen("tcp", *rtmpAddr)
	fmt.Println("==========", *rtmpAddr)
	if err != nil {
		log.Fatal(err)
	}

	var rtmpServer *rtmp.Server

	rtmpServer = rtmp.NewRtmpServer(stream, nil)

	defer func() {
		if r := recover(); r != nil {
			log.Println("RTMP server panic: ", r)
		}
	}()
	log.Println("RTMP Listen On", *rtmpAddr)
	rtmpServer.Serve(rtmpListen)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("livego panic: ", r)
			time.Sleep(1 * time.Second)
		}
	}()
	err := configure.LoadConfig(*configfilename)
	if err != nil {
		return
	}
	stream := rtmp.NewRtmpStream()
	startRtmp(stream)
}
