package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/toxyl/glog"
	metrics "github.com/toxyl/metric-nexus"
	"github.com/toxyl/spider/log"
	"gopkg.in/yaml.v3"
)

var (
	colHost    = glog.NewTableColumnLeftCustom("Host", ' ', nil)
	colAction  = glog.NewTableColumnLeftCustom("Action", ' ', nil)
	colSuccess = glog.NewTableColumnLeftCustom("Success", ' ', nil)
	colActive  = glog.NewTableColumnLeftCustom("Time", ' ', func(a ...interface{}) string {
		str := make([]string, len(a))
		for i, v := range a {
			switch t := v.(type) {
			case float64:
				str[i] = glog.DurationShort(t, glog.DURATION_SCALE_AVERAGE)
			case string:
				if n, err := strconv.Atoi(fmt.Sprint(t)); err == nil {
					str[i] = glog.DurationShort(n, glog.DURATION_SCALE_AVERAGE)
				}
			case time.Duration:
				str[i] = glog.DurationShort(t.Seconds(), glog.DURATION_SCALE_AVERAGE)
			}
		}
		return strings.Join(str, ", ")
	})
	taskHistory = struct {
		data map[string][]struct {
			name     string
			success  bool
			duration float64
		}
	}{
		data: map[string][]struct {
			name     string
			success  bool
			duration float64
		}{},
	}
	failedHosts = map[string]struct{}{}
)

func isHostConfigured(host string) bool {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	configFile := filepath.Join(usr.HomeDir, ".ssh", "config")
	file, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "Host "+host {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return false
}

type Inventory struct {
	client      *metrics.Client
	Credentials struct {
		User  string `yaml:"user"`
		Group string `yaml:"group"`
		Key   string `yaml:"key"`
	} `yaml:"credentials"`
	Spider struct {
		Source       string   `yaml:"source"`
		Destination  string   `yaml:"destination"`
		Spiders      []int    `yaml:"spiders"`
		Taunts       []string `yaml:"taunts"`
		Whitelist    []string `yaml:"whitelist"`
		AttackLength int      `yaml:"attack_length"`
	} `yaml:"spider"`
	MetricNexus struct {
		Source      string `yaml:"source"`
		Destination string `yaml:"destination"`
		IP          string `yaml:"ip"`
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		Key         string `yaml:"key"`
		CertFile    string `yaml:"ssl_cert_file"`
		KeyFile     string `yaml:"ssl_key_file"`
	} `yaml:"metric_nexus"`
	Hosts []string `yaml:"hosts"`
}

func (i *Inventory) WithMetricNexusHost(logMsg, task string) error {
	var (
		success = true
		t       = time.Now()
	)

	if logMsg != "" {
		log.Normal(logMsg)
	}

	if err := execTask(i.MetricNexus.Host, task); err != nil {
		log.Failed("%s: '%s' failed: %s", glog.Auto(i.MetricNexus.IP), glog.Auto(task), err)
		success = false
	}

	taskHistory.data[i.MetricNexus.IP] = append(taskHistory.data[i.MetricNexus.IP], struct {
		name     string
		success  bool
		duration float64
	}{
		name:     task,
		success:  success,
		duration: time.Since(t).Seconds(),
	})

	return nil
}

func (i *Inventory) WithSpiderHosts(logMsg, task string) error {
	var (
		wg  = sync.WaitGroup{}
		res = struct {
			lock       *sync.Mutex
			success    map[string]time.Duration
			failure    map[string]time.Duration
			setFailure func(host string, t time.Time)
			setSuccess func(host string, t time.Time)
		}{
			lock:    &sync.Mutex{},
			success: map[string]time.Duration{},
			failure: map[string]time.Duration{},
		}
	)

	res.setFailure = func(host string, t time.Time) {
		res.lock.Lock()
		defer res.lock.Unlock()
		res.failure[host] = time.Since(t)
	}

	res.setSuccess = func(host string, t time.Time) {
		res.lock.Lock()
		defer res.lock.Unlock()
		res.success[host] = time.Since(t)
	}

	if logMsg != "" {
		log.Normal(logMsg)
	}

	for _, host := range i.Hosts {
		if _, ok := failedHosts[host]; ok {
			continue // let's not process tasks on hosts with failures, shall we?
		}
		wg.Add(1)
		go func(host, dst string, t time.Time) {
			defer wg.Done()

			if !isHostConfigured(host) {
				log.Error("Host '%s' is not configured in the ~/.ssh/config file", glog.Auto(host))
				res.setFailure(host, t)
				return
			}

			if err := execTask(host, task); err != nil {
				log.Failed("%s: '%s' failed: %s", glog.Auto(host), glog.Auto(task), err)
				res.setFailure(host, t)
				return
			}

			res.setSuccess(host, t)
		}(host, fmt.Sprintf("%s@%s", i.Credentials.User, host), time.Now())
	}
	wg.Wait()

	for _, host := range i.Hosts {
		if _, ok := taskHistory.data[host]; !ok {
			taskHistory.data[host] = []struct {
				name     string
				success  bool
				duration float64
			}{}
		}

		if t, ok := res.success[host]; ok {
			taskHistory.data[host] = append(taskHistory.data[host], struct {
				name     string
				success  bool
				duration float64
			}{
				name:     task,
				success:  true,
				duration: t.Seconds(),
			})
			continue
		}

		if t, ok := res.failure[host]; ok {
			failedHosts[host] = struct{}{}
			taskHistory.data[host] = append(taskHistory.data[host], struct {
				name     string
				success  bool
				duration float64
			}{
				name:     task,
				success:  false,
				duration: t.Seconds(),
			})
		}
	}

	return nil
}

func printTaskHistory() {
	hosts := []string{}
	for host := range taskHistory.data {
		hosts = append(hosts, host)
	}
	sort.Strings(hosts)
	for _, host := range hosts {
		for _, action := range taskHistory.data[host] {
			colHost.Push(host)
			colAction.Push(action.name)
			colSuccess.Push(action.success)
			colActive.Push(action.duration)
		}
	}

	log.Table(colHost, colAction, colSuccess, colActive)
}

func readInventory(file string) (*Inventory, error) {
	yamlData, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var inventory Inventory
	err = yaml.Unmarshal(yamlData, &inventory)
	if err != nil {
		return nil, err
	}
	inventory.client = metrics.NewClient(inventory.MetricNexus.IP, inventory.MetricNexus.Port, inventory.MetricNexus.Key, true)

	return &inventory, nil
}
