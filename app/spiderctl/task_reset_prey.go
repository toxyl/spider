package main

import (
	"github.com/toxyl/spider/utils"
)

func taskResetPrey(host string) error {
	for _, s := range inventory.Spider.Spiders {
		_ = inventory.client.Update(utils.GetMetricName(s, "prey"), 0)
		_ = sshRm(host, utils.GetMetricFileName(s, "prey"))
	}
	return nil
}
