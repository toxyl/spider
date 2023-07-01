package main

func taskStartSpider(host string) error {
	return sshServiceStart(host, "spider")
}

func taskStartMetricNexus(host string) error {
	return sshServiceStart(host, "metric-nexus")
}
