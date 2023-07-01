package main

func taskStopSpider(host string) error {
	_ = sshServiceStop(host, "spider")
	return nil
}

func taskStopMetricNexus(host string) error {
	_ = sshServiceStop(host, "metric-nexus")
	return nil
}
