package main

import (
	"github.com/toxyl/spider/utils"
)

func taskResetSpider(host string) error {
	for _, s := range inventory.Spider.Spiders {
		_ = inventory.client.Update(utils.GetMetricName(s, "prey"), 0)
		_ = inventory.client.Update(utils.GetMetricName(s, "kills"), 0)
		_ = inventory.client.Update(utils.GetMetricName(s, "wasted"), 0)
		_ = inventory.client.Update(utils.GetMetricName(s, "active"), 0)
		_ = sshRm(host, utils.GetMetricFileName(s, "prey"))
	}
	_ = inventory.client.Update(utils.GetMetricName(0, "uptime"), 0)
	return nil
}
