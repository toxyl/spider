package main

func playbookStop() {
	inventory.WithSpiderHosts("Stopping Spiders...", "stop-spider")
	printTaskHistory()
}
