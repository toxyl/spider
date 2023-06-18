package main

import (
	"net"

	"github.com/toxyl/glog"
)

var (
	log = glog.NewLoggerSimple("Spider")
)

func connInfo(conn net.Conn, action string) {
	log.Info("%s spider %s prey %s", glog.Auto(conn2name(conn)), colorizeAction(action), glog.Auto(conn2prey(conn)))
}

func spiderInfo(spider int, action string, suffix string) {
	log.Info("%s spider %s%s", glog.Auto(spider2name(spider)), colorizeAction(action), suffix)
}

func spiderNotOK(spider int, action string, suffix string) {
	log.NotOK("%s spider %s%s", glog.Auto(spider2name(spider)), colorizeAction(action), suffix)
}
