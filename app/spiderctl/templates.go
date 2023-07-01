package main

import (
	_ "embed"
	"html/template"
	"os"
	"path/filepath"
)

//go:embed "templates/metric_nexus.config.gotmpl"
var tmplMetricNexusConfig string

//go:embed "templates/metric_nexus.service.gotmpl"
var tmplMetricNexusService string

//go:embed "templates/spider.config.gotmpl"
var tmplSpiderConfig string

//go:embed "templates/spider.service.gotmpl"
var tmplSpiderService string

func parseTemplate(templateName, dstFile string, data any) error {
	tmpl := template.New(templateName)

	content := ""
	switch templateName {
	case "metric-nexus-config":
		content = tmplMetricNexusConfig
	case "metric-nexus-service":
		content = tmplMetricNexusService
	case "spider-config":
		content = tmplSpiderConfig
	case "spider-service":
		content = tmplSpiderService
	}

	_, err := tmpl.Parse(content)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		return err
	}
	return nil
}

func prepareTemplate(templateName, dstPath string, data any) error {
	ftmp := dstPath
	if ftmp == "" {
		ftmp = getTempFilePath(templateName)
	}
	if err := parseTemplate(templateName, ftmp, data); err != nil {
		return err
	}
	return nil
}

func prepareTemplates(inventory *Inventory) error {
	if err := prepareTemplate(
		"spider-config",
		filepath.Join(inventory.Spider.Source, "config.go"),
		struct {
			Spiders         []int
			Taunts          []string
			Whitelist       []string
			AttackLength    int
			MetricNexusHost string
			MetricNexusPort int
			MetricNexusKey  string
		}{
			Spiders:         inventory.Spider.Spiders,
			Taunts:          inventory.Spider.Taunts,
			Whitelist:       inventory.Spider.Whitelist,
			AttackLength:    inventory.Spider.AttackLength,
			MetricNexusHost: inventory.MetricNexus.IP,
			MetricNexusPort: inventory.MetricNexus.Port,
			MetricNexusKey:  inventory.MetricNexus.Key,
		},
	); err != nil {
		return err
	}

	if err := prepareTemplate(
		"spider-service",
		"",
		struct {
			User       string
			Group      string
			Executable string
		}{
			User:       inventory.Credentials.User,
			Group:      inventory.Credentials.Group,
			Executable: filepath.Join(inventory.Spider.Destination, "spider"),
		},
	); err != nil {
		return err
	}

	if err := prepareTemplate(
		"metric-nexus-config",
		filepath.Join(inventory.MetricNexus.Source, "config.go"),
		struct {
			Host      string
			Port      int
			CertFile  string
			KeyFile   string
			StateFile string
			APIKeys   []string
		}{
			Host:      inventory.MetricNexus.IP,
			Port:      inventory.MetricNexus.Port,
			CertFile:  inventory.MetricNexus.CertFile,
			KeyFile:   inventory.MetricNexus.KeyFile,
			StateFile: filepath.Join(inventory.MetricNexus.Destination, "state.yaml"),
			APIKeys:   []string{inventory.MetricNexus.Key},
		},
	); err != nil {
		return err
	}

	if err := prepareTemplate(
		"metric-nexus-service",
		"",
		struct {
			User       string
			Group      string
			Executable string
		}{
			User:       inventory.Credentials.User,
			Group:      inventory.Credentials.Group,
			Executable: filepath.Join(inventory.MetricNexus.Destination, "metric-nexus"),
		},
	); err != nil {
		return err
	}

	return nil
}
