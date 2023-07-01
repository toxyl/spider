package main

import (
	"path/filepath"
)

func taskUploadSpider(host string) error {
	_ = sshMkdir(host, filepath.Join(inventory.Spider.Destination, "data"))

	if err := scpExecutableToHost(host, "spider", inventory.Spider.Destination); err != nil {
		return err
	}

	if err := scpServiceToHost(host, "spider"); err != nil {
		return err
	}

	if err := sshServiceEnable(host, "spider"); err != nil {
		return err
	}

	return nil
}

func taskUploadMetricNexus(host string) error {
	_ = sshMkdir(host, inventory.MetricNexus.Destination)

	if err := scpExecutableToHost(host, "metric-nexus", inventory.MetricNexus.Destination); err != nil {
		return err
	}

	if err := scpServiceToHost(host, "metric-nexus"); err != nil {
		return err
	}

	if err := sshServiceEnable(host, "metric-nexus"); err != nil {
		return err
	}

	return nil
}
