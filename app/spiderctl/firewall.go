package main

import (
	"github.com/toxyl/glog"
	"github.com/toxyl/spider/log"
)

func firewallOpen(host string, ports []int) error {
	for _, port := range ports {
		if sshCommandExists(host, "ufw") {
			err := sshOpenPortUFW(host, port)
			if err != nil {
				return err
			}
		} else if sshCommandExists(host, "iptables") {
			err := sshOpenPortIPTables(host, port)
			if err != nil {
				return err
			}
		} else {
			log.Error("Error: Neither ufw nor iptables found on %s. Unable to open port %s.", glog.Auto(host), glog.Auto(port))
		}
	}

	return nil
}
