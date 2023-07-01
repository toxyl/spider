package main

import (
	"os"
	"strconv"
)

func playbookReset() {
	if len(os.Args) > 3 {
		for _, arg := range os.Args[3:] {
			if spider, err := strconv.Atoi(arg); err == nil {
				inventory.Spider.Spiders = append(inventory.Spider.Spiders, spider)
			}
		}
	}
	inventory.WithSpiderHosts("Stopping Spiders...", "stop-spider")
	inventory.WithSpiderHosts("Resetting Spiders...", "reset-spider")
	inventory.WithSpiderHosts("Starting Spiders...", "start-spider")
	printTaskHistory()
}
