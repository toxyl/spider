package main

import "github.com/toxyl/spider/stats"

func playbookStats() {
	inventory.WithSpiderHosts("", "stats")
	stats.PrintStatsTable(prey, kills, wasted, active)
}
