package main

import (
	"github.com/toxyl/glog"
	"github.com/toxyl/spider/log"
)

func taskFirewallSpider(host string) error {
	err := firewallOpen(host, inventory.Spider.Spiders)
	if err != nil {
		log.Error("%s: Error opening ports: %s", glog.Auto(host), err)
		return err
	}

	return nil
}

func taskFirewallMetricNexus(host string) error {
	err := firewallOpen(host, []int{inventory.MetricNexus.Port})
	if err != nil {
		log.Error("%s: Error opening ports: %s", glog.Auto(host), err)
		return err
	}

	return nil
}
