package main

import (
	"net"

	"github.com/toxyl/glog"
	"github.com/toxyl/spider/log"
	"github.com/toxyl/spider/services"
)

func colorizeAction(action string) string {
	if action == "kills" {
		action = glog.WrapRed(action)
	} else if action == "attacks" {
		action = glog.WrapOrange(action)
	} else if action == "pokes" {
		action = glog.WrapYellow(action)
	} else {
		action = glog.Auto(action)
	}
	return action
}

func SpiderPreyInfo(conn net.Conn, action string) {
	log.Info("%s spider %s prey %s", glog.Auto(services.Conn2name(conn)), colorizeAction(action), glog.Auto(services.Conn2prey(conn)))
}

func SpiderInfo(spider int, action string, suffix string) {
	log.Info("%s spider %s%s", glog.Auto(services.Spider2name(spider)), colorizeAction(action), suffix)
}

func SpiderFailed(spider int, action string, suffix string) {
	log.Failed("%s spider %s%s", glog.Auto(services.Spider2name(spider)), colorizeAction(action), suffix)
}
