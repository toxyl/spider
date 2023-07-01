package main

import (
	"os"

	"github.com/toxyl/spider/log"
)

func playbookDeploy() {
	prepareTemplates(inventory)
	log.Normal("Building MetricNexus...")
	err := buildMetricNexus(inventory)
	if err != nil {
		log.Error("Error building MetricNexus binary: %s", err)
		os.Exit(1)
	}
	inventory.WithMetricNexusHost("Stopping MetricNexus...", "stop-metric-nexus")
	inventory.WithMetricNexusHost("Uploading MetricNexus...", "upload-metric-nexus")
	inventory.WithMetricNexusHost("Setting up MetricNexus firewall...", "firewall-metric-nexus")
	inventory.WithMetricNexusHost("Starting MetricNexus...", "start-metric-nexus")

	log.Normal("Building Spider...")
	err = buildSpider(inventory)
	if err != nil {
		log.Error("Error building Spider binary: %s", err)
		os.Exit(1)
	}
	inventory.WithSpiderHosts("Stopping Spiders...", "stop-spider")
	inventory.WithSpiderHosts("Uploading Spiders...", "upload-spider")
	inventory.WithSpiderHosts("Resetting prey...", "reset-prey")
	inventory.WithSpiderHosts("Setting up Spider firewalls...", "firewall-spider")
	inventory.WithSpiderHosts("Starting Spiders...", "start-spider")
	printTaskHistory()
}
