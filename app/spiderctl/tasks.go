package main

import "errors"

var (
	TASKS = map[string]func(host string) error{
		"stats":                 taskGetStats,
		"reset-spider":          taskResetSpider,
		"reset-prey":            taskResetPrey,
		"firewall-spider":       taskFirewallSpider,
		"firewall-metric-nexus": taskFirewallMetricNexus,
		"start-spider":          taskStartSpider,
		"stop-spider":           taskStopSpider,
		"start-metric-nexus":    taskStartMetricNexus,
		"stop-metric-nexus":     taskStopMetricNexus,
		"upload-spider":         taskUploadSpider,
		"upload-metric-nexus":   taskUploadMetricNexus,
	}
)

func execTask(host, name string) error {
	if f, ok := TASKS[name]; ok {
		return f(host)
	}
	return errors.New("no matching tasks found")
}
