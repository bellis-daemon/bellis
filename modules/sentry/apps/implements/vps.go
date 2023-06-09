package implements

import (
	"context"
	"encoding/json"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"net/http"
	"time"
)

type VPS struct {
	Options vpsOptions
	client  *http.Client
}

func (this *VPS) Fetch(ctx context.Context) (any, error) {
	resp, err := this.client.Get(this.Options.URL)
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

func (this *VPS) Init(setOptions func(options any) error) error {
	this.client = http.DefaultClient
	this.client.Timeout = 5 * time.Second
	return setOptions(&this.Options)
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

type vpsMetrics struct {
	Load   *load.AvgStat          `json:"load"`
	CPU    float64                `json:"cpus"`
	Memory *mem.VirtualMemoryStat `json:"memory"`
	Disks  *disk.UsageStat        `json:"disks"`
	Host   *host.InfoStat         `json:"host"`
}
