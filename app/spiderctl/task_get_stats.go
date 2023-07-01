package main

import (
	"sync"

	"github.com/toxyl/spider/utils"
)

var (
	lock   = &sync.Mutex{}
	prey   = map[int]int{}
	kills  = map[int]int{}
	wasted = map[int]float64{}
	active = map[int]float64{}
)

func taskGetStats(host string) error {
	lock.Lock()
	defer lock.Unlock()
	for _, s := range inventory.Spider.Spiders {
		if _, ok := prey[s]; !ok {
			if p, err := inventory.client.Read(utils.GetMetricName(s, "prey")); err == nil {
				prey[s] = int(p)
			}
		}

		if _, ok := kills[s]; !ok {
			if k, err := inventory.client.Read(utils.GetMetricName(s, "kills")); err == nil {
				kills[s] = int(k)
			}
		}

		if _, ok := wasted[s]; !ok {
			if w, err := inventory.client.Read(utils.GetMetricName(s, "wasted")); err == nil {
				wasted[s] = w
			}
		}

		if _, ok := active[s]; !ok {
			if w, err := inventory.client.Read(utils.GetMetricName(s, "active")); err == nil {
				active[s] = w
			}
		}
	}
	return nil
}
