package main

import (
	"os"
	"strings"

	"github.com/toxyl/glog"
	"github.com/toxyl/spider/log"
)

func playbookExec() {
	if len(os.Args) <= 3 {
		log.Warning("Not enough arguments!")
		log.Normal("")
		log.Normal("Usage: %s %s %s [cmd]", os.Args[0], os.Args[1], os.Args[2])
		os.Exit(0)
	}
	cmd := strings.Join(os.Args[3:], " ")

	for _, h := range inventory.Hosts {
		log.Info("%s:", glog.Auto(h))
		sshExecInteractive(h, cmd)
		log.Normal("")
	}
}
