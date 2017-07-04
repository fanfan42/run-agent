package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

//Metrics represents a basic struct with systems stats
type Metrics struct {
	//CPU
	CPUStats []cpu.TimesStat `json:"cpustats"`
	//Memory
	MemStats  *mem.VirtualMemoryStat `json:"memstats"`
	SwapStats *mem.SwapMemoryStat    `json:"swapstats"`
	//Disk
	IOStats        map[string]disk.IOCountersStat `json:"iostats"`
	PartitionStats []disk.PartitionStat           `json:"partitionstats"`
	UsageStats     []*disk.UsageStat              `json:"usagestats"`
	//Load Average
	LoadAverage *load.AvgStat  `json:"loadaverage"`
	MiscStats   *load.MiscStat `json:"miscstats"`
}

func getLoadAverage() (*load.AvgStat, *load.MiscStat, error) {
	lastats, err := load.Avg()
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting load average information: %v", err)
	}
	msstats, err := load.Misc()
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting Misc information: %v", err)
	}
	return lastats, msstats, nil
}

func getDiskInfo() (map[string]disk.IOCountersStat, []disk.PartitionStat, []*disk.UsageStat, error) {
	pstats, err := disk.Partitions(false)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Error getting partition stats: %v", err)
	}
	iostats, err := disk.IOCounters()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Error getting io stats for %v", err)
	}
	//var iostats []map[string]disk.IOCountersStat
	var usstats []*disk.UsageStat
	//fmt.Println(pstats)
	for _, partition := range pstats {
		//fmt.Println(partition)

		//iostats = append(iostats, iostat)
		usstat, err := disk.Usage(partition.Device)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("Error getting usage stats for %s: %v", partition.Device, err)
		}
		usstats = append(usstats, usstat)
		//partitions = append(partitions, partition.Device)
		//fmt.Println(disk.IOCounters(partition.Device))
	}
	return iostats, pstats, usstats, nil
}

func getSWAPInfo() (*mem.SwapMemoryStat, error) {
	swapstats, err := mem.SwapMemory()
	if err != nil {
		return nil, fmt.Errorf("Error getting swap information: %v", err)
	}
	return swapstats, nil
}

func getRAMInfo() (*mem.VirtualMemoryStat, error) {
	memstats, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("Error getting memory information: %v", err)
	}
	return memstats, nil
}

func getCPUInfo() ([]cpu.TimesStat, error) {
	cpustats, err := cpu.Times(true)
	if err != nil {
		return nil, fmt.Errorf("Error getting cpu information: %v", err)
	}
	return cpustats, nil
}

func (agent *Agent) basicMetrics() {
	for {
		c, err := getCPUInfo()
		if err != nil {
			log.Println(err)
		}
		m, err := getRAMInfo()
		if err != nil {
			log.Println(err)
		}
		s, err := getSWAPInfo()
		if err != nil {
			log.Println(err)
		}
		io, ps, us, err := getDiskInfo()
		if err != nil {
			log.Println(err)
		}
		la, ms, err := getLoadAverage()
		metrics.CPUStats = c
		metrics.MemStats = m
		metrics.SwapStats = s
		metrics.IOStats = io
		metrics.PartitionStats = ps
		metrics.UsageStats = us
		metrics.LoadAverage = la
		metrics.MiscStats = ms
		time.Sleep(3 * time.Second)
	}
}
