package main

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type rtinfoCollector struct {
	rtinfo          *rtinfo
	mux             sync.RWMutex
	UptimeDesc      *prometheus.Desc
	CPUUsageDesc    *prometheus.Desc
	CPUTempDesc     *prometheus.Desc
	LoadAvgDesc     *prometheus.Desc
	MemoryTotalDesc *prometheus.Desc
	MemoryUsedDesc  *prometheus.Desc
	SwapTotalDesc   *prometheus.Desc
	SwapFreeDesc    *prometheus.Desc
	NetworkRxDesc   *prometheus.Desc
	NetworkTxDesc   *prometheus.Desc
	DiskWriteDesc   *prometheus.Desc
	DiskReadDesc    *prometheus.Desc
}

func newRtinfoCollector() *rtinfoCollector {
	return &rtinfoCollector{
		UptimeDesc: prometheus.NewDesc(
			"system_uptime_second",
			"Uptime of the system",
			[]string{"hostname", "remoteIP"},
			nil,
		),
		CPUUsageDesc: prometheus.NewDesc(
			"cpu_usage_percent",
			"Percentage of the cpu used",
			[]string{"hostname", "remoteIP", "cpu"},
			nil,
		),
		CPUTempDesc: prometheus.NewDesc(
			"cpu_temperature_celcius",
			"Average temperature of the cpu",
			[]string{"hostname", "remoteIP"},
			nil,
		),
		LoadAvgDesc: prometheus.NewDesc(
			"load_average_unit",
			"Average load of each cpu",
			[]string{"hostname", "remoteIP", "cpu"},
			nil,
		),
		SwapFreeDesc: prometheus.NewDesc(
			"swap_free_byte",
			"Free swap of the system",
			[]string{"hostname", "remoteIP"},
			nil,
		),
		SwapTotalDesc: prometheus.NewDesc(
			"swap_total_byte",
			"Total amount of swap",
			[]string{"hostname", "remoteIP"},
			nil,
		),
		MemoryTotalDesc: prometheus.NewDesc(
			"memory_total_byte",
			"Total memory available of the system",
			[]string{"hostname", "remoteIP"},
			nil,
		),
		MemoryUsedDesc: prometheus.NewDesc(
			"memory_used_byte",
			"Memory used of the system",
			[]string{"hostname", "remoteIP"},
			nil,
		),
		NetworkRxDesc: prometheus.NewDesc(
			"network_received_byte",
			"Data received on network interface ",
			[]string{"hostname", "remoteIP", "nic", "ip"},
			nil,
		),
		NetworkTxDesc: prometheus.NewDesc(
			"network_sent_byte",
			"Data sent on network interface ",
			[]string{"hostname", "remoteIP", "nic", "ip"},
			nil,
		),
		DiskReadDesc: prometheus.NewDesc(
			"disk_read_byte",
			"Data read on disk",
			[]string{"hostname", "remoteIP", "disk"},
			nil,
		),
		DiskWriteDesc: prometheus.NewDesc(
			"disk_write_byte",
			"Data written on disk",
			[]string{"hostname", "remoteIP", "disk"},
			nil,
		),
	}
}

// RAMTotal  int `json:"ram_total"`
// RAMUsed   int `json:"ram_used"`
// SwapTotal int `json:"swap_total"`
// SwapFree  int `json:"swap_free"`

func (c *rtinfoCollector) SetInfo(rtinfo *rtinfo) {
	c.mux.Lock()
	c.rtinfo = rtinfo
	c.mux.Unlock()
}

func (c *rtinfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.UptimeDesc
	ch <- c.CPUUsageDesc
	ch <- c.CPUTempDesc
	ch <- c.LoadAvgDesc
	ch <- c.MemoryTotalDesc
	ch <- c.MemoryUsedDesc
	ch <- c.SwapTotalDesc
	ch <- c.SwapFreeDesc
	ch <- c.NetworkRxDesc
	ch <- c.NetworkTxDesc
}

func (c *rtinfoCollector) Collect(ch chan<- prometheus.Metric) {
	c.mux.Lock()
	defer c.mux.Unlock()

	for _, rtinfo := range c.rtinfo.Rtinfo {

		ch <- prometheus.MustNewConstMetric(
			c.UptimeDesc,
			prometheus.CounterValue,
			float64(rtinfo.Uptime),
			rtinfo.Hostname, rtinfo.Remoteip,
		)

		for i, cpu := range rtinfo.CPUUsage {
			ch <- prometheus.MustNewConstMetric(
				c.CPUUsageDesc,
				prometheus.GaugeValue,
				float64(cpu),
				rtinfo.Hostname,
				rtinfo.Remoteip,
				fmt.Sprintf("%d", i),
			)
		}

		for i, load := range rtinfo.Loadavg {
			ch <- prometheus.MustNewConstMetric(
				c.LoadAvgDesc,
				prometheus.GaugeValue,
				float64(load),
				rtinfo.Hostname,
				rtinfo.Remoteip,
				fmt.Sprintf("%d", i),
			)
		}

		ch <- prometheus.MustNewConstMetric(
			c.CPUTempDesc,
			prometheus.GaugeValue,
			float64(rtinfo.Sensors.CPU.Average),
			rtinfo.Hostname,
			rtinfo.Remoteip,
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemoryTotalDesc,
			prometheus.GaugeValue,
			float64(rtinfo.Memory.RAMTotal),
			rtinfo.Hostname,
			rtinfo.Remoteip,
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemoryUsedDesc,
			prometheus.GaugeValue,
			float64(rtinfo.Memory.RAMUsed),
			rtinfo.Hostname,
			rtinfo.Remoteip,
		)

		ch <- prometheus.MustNewConstMetric(
			c.SwapTotalDesc,
			prometheus.GaugeValue,
			float64(rtinfo.Memory.SwapTotal),
			rtinfo.Hostname,
			rtinfo.Remoteip,
		)

		ch <- prometheus.MustNewConstMetric(
			c.SwapFreeDesc,
			prometheus.GaugeValue,
			float64(rtinfo.Memory.SwapFree),
			rtinfo.Hostname,
			rtinfo.Remoteip,
		)

		for _, nic := range rtinfo.Network {
			ch <- prometheus.MustNewConstMetric(
				c.NetworkTxDesc,
				prometheus.CounterValue,
				float64(nic.TxData),
				rtinfo.Hostname, rtinfo.Remoteip, nic.Name, nic.IP,
			)

			ch <- prometheus.MustNewConstMetric(
				c.NetworkRxDesc,
				prometheus.CounterValue,
				float64(nic.RxData),
				rtinfo.Hostname, rtinfo.Remoteip, nic.Name, nic.IP,
			)
		}

		for _, disk := range rtinfo.Disks {
			ch <- prometheus.MustNewConstMetric(
				c.DiskReadDesc,
				prometheus.CounterValue,
				float64(disk.BytesRead),
				rtinfo.Hostname, rtinfo.Remoteip, disk.Name,
			)

			ch <- prometheus.MustNewConstMetric(
				c.DiskWriteDesc,
				prometheus.CounterValue,
				float64(disk.BytesWritten),
				rtinfo.Hostname, rtinfo.Remoteip, disk.Name,
			)
		}
	}
}
