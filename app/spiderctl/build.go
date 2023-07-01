package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func buildSpider(inventory *Inventory) error {
	sourceDir := inventory.Spider.Source
	cmd := exec.Command("go", "build", "-C", sourceDir, "-trimpath", "-o", filepath.Join(os.TempDir(), "spider.tmp"))
	cmd.Env = append(cmd.Environ(), "CGO_ENABLED=0")
	return cmd.Run()
}

func buildMetricNexus(inventory *Inventory) error {
	sourceDir := inventory.MetricNexus.Source
	cmd := exec.Command("go", "build", "-C", sourceDir, "-trimpath", "-o", filepath.Join(os.TempDir(), "metric-nexus.tmp"))
	cmd.Env = append(cmd.Environ(), "CGO_ENABLED=0")
	return cmd.Run()
}
