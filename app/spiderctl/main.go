package main

import (
	"os"

	"github.com/toxyl/glog"
	"github.com/toxyl/spider/log"
)

var (
	inventory *Inventory
)

func init() {
	glog.LoggerConfig.ShowIndicator = false
	glog.LoggerConfig.ShowSubsystem = false
	glog.LoggerConfig.ShowDateTime = false
	glog.LoggerConfig.ShowRuntimeMilliseconds = false
	glog.LoggerConfig.ShowRuntimeSeconds = false
}

func main() {
	if len(os.Args) < 3 {
		playbooks := []string{}
		for key := range PLAYBOOKS {
			playbooks = append(playbooks, key)
		}
		log.Info("Usage:     %s [inventory file] [playbook]", glog.File(os.Args[0]))
		log.Info("Example:   %s inventory.yaml stats", glog.File(os.Args[0]))
		log.Info("Playbooks: %s", glog.Auto(playbooks))
		os.Exit(0)
	}

	inv, err := readInventory(os.Args[1])
	if err != nil {
		log.Error("Error reading inventory file: %s", err)
		os.Exit(1)
	}
	inventory = inv

	if playbook, ok := PLAYBOOKS[os.Args[2]]; ok {
		playbook()
	} else {
		playbooks := []string{}
		for key := range PLAYBOOKS {
			playbooks = append(playbooks, key)
		}
		log.Warning(
			"Unknown playbook '%s', use one of these: %s",
			glog.Auto(os.Args[2]),
			glog.Auto(playbooks),
		)
	}
}
