package vps

import (
	"context"
	"encoding/json"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/option"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type VPS struct {
	options vpsOptions
	client  *http.Client
}

func (this *VPS) Fetch(ctx context.Context) (status.Status, error) {
	resp, err := this.client.Get(this.options.URL)
	if err != nil {
		return &vpsStatus{}, err
	}
	defer resp.Body.Close()
	var m vpsMetrics
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return &vpsStatus{}, err
	}
	return &vpsStatus{
		Load1:       m.Load.Load1,
		Load5:       m.Load.Load5,
		Load15:      m.Load.Load15,
		CPUUsage:    m.CPU,
		MemoryUsage: m.Memory.UsedPercent,
		DiskUsage:   m.Disks.UsedPercent,
		Hostname:    m.Host.Hostname,
		OS:          m.Host.OS,
		Uptime:      m.Host.Uptime,
		Platform:    m.Host.Platform,
		Process:     m.Host.Procs,
	}, nil
}

type vpsOptions struct {
	URL string `json:"URL"`
}

type vpsStatus struct {
	Load1       float64 `json:"load1"`
	Load5       float64 `json:"load5"`
	Load15      float64 `json:"load15"`
	CPUUsage    float64 `json:"CPUUsage"`
	MemoryUsage float64 `json:"memoryUsage"`
	DiskUsage   float64 `json:"diskUsage"`
	Hostname    string  `json:"hostname"`
	OS          string  `json:"OS"`
	Uptime      uint64  `json:"uptime"`
	Platform    string  `json:"platform"`
	Process     uint64  `json:"process"`
}

func (this *vpsStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}

type vpsMetrics struct {
	Load   *load.AvgStat          `json:"load"`
	CPU    float64                `json:"cpus"`
	Memory *mem.VirtualMemoryStat `json:"memory"`
	Disks  *disk.UsageStat        `json:"disks"`
	Host   *host.InfoStat         `json:"host"`
}

func init() {
	implements.Register("vps", func(options bson.M) implements.Implement {
		return &VPS{
			options: option.ToOption[vpsOptions](options),
			client:  http.DefaultClient,
		}
	})
}
