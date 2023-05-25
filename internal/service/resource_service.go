package service

import (
	"Time_k8s_operator/pkg/logger"
	"math"

	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/shirou/gopsutil/net"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

var ()

type ResourceService struct {
	logger *logrus.Logger
	// dao    *dao.UserDao
}

func NewResourceService() *ResourceService {
	return &ResourceService{
		logger: logger.Logger(),
		//		dao:    dao.NewUserDao(),
	}
}

type ReceiveBytes uint64
type TransmitBytes uint64

type DownloadFlow string
type UploadFlow string

var wg sync.WaitGroup
var (
	cpu_core    int
	cpu_utility float64

	memory_total         string
	memory_free_percent  float64
	memory_used_percent  float64
	memory_other_percent float64

	disk_total         string
	disk_free_percent  float64
	disk_used_percent  float64
	disk_other_percent float64

	net_download DownloadFlow
	net_upload   UploadFlow

	err error
)

func getCPU() {
	cpu_core, _ = cpu.Counts(true)
	c, _ := cpu.Percent(time.Second, false)
	cpu_utility, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", c[0]), 64)
	fmt.Println(cpu_core, cpu_utility)
	wg.Done()
}

func getMemory() {
	memInfo, _ := mem.VirtualMemory()
	total := memInfo.Total
	memory_total = strconv.FormatInt(int64(math.Ceil(float64(total)/1024/1024/1024)), 10) + "GB"

	memory_free_percent = float64(float64(memInfo.Available)/float64(memInfo.Total)) * 100
	memory_free_percent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", memory_free_percent), 64)

	memory_used_percent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", memInfo.UsedPercent), 64)

	memory_other_percent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", 100-memory_free_percent-memory_used_percent), 64)

	fmt.Println(memory_free_percent, memory_used_percent, memory_other_percent, memory_total)
	wg.Done()
}

func getDisk() {
	partitioncstat, _ := disk.Partitions(false)
	for _, partition := range partitioncstat {
		device := partition.Device
		if strings.HasPrefix(device, "/dev/") {
			if strings.HasPrefix(device, "/dev/sda") {
				continue
			}
			diskInfo, _ := disk.Usage(partition.Mountpoint)
			disk_total = strconv.FormatInt(int64(math.Ceil(float64(diskInfo.Total)/1024/1024/1024)), 10) + "GB"
			fmt.Println(disk_total, "===")
			disk_free_percent = float64(float64(diskInfo.Free)/float64(diskInfo.Total)) * 100

			disk_free_percent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", disk_free_percent), 64)

			disk_used_percent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", diskInfo.UsedPercent), 64)

			disk_other_percent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", 100.0-disk_free_percent-disk_used_percent), 64)
			if disk_other_percent <= 0 {
				disk_other_percent = 0
			}
		}
	}
	fmt.Println(disk_free_percent, disk_used_percent, disk_other_percent, disk_total)
	wg.Done()
}

func UploadDownloadFlow(dev string) (DownloadFlow, UploadFlow, error) {
	down, up, err := TotalFlowByDevice(dev)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	time.Sleep(time.Second * 1)
	down2, up2, err := TotalFlowByDevice(dev)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	downStr := strconv.FormatInt(int64((down2-down)/1024), 10) + "Kbps"
	upStr := strconv.FormatInt(int64((up2-up)/1024), 10) + "Kbps"
	return DownloadFlow(downStr), UploadFlow(upStr), nil
}

func TotalFlowByDevice(dev string) (ReceiveBytes, TransmitBytes, error) {
	devInfo, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}
	var receive int = -1
	var transmit int = -1
	var receiveBytes uint64
	var transmitBytes uint64
	lines := strings.Split(string(devInfo), "\n")
	for _, line := range lines {
		if strings.Contains(line, dev) {
			i := 0
			fields := strings.Split(line, ":")
			for _, field := range fields {
				if strings.Contains(field, dev) {
					i = 1
				} else {
					values := strings.Fields(field)
					for _, value := range values {
						if receive == i {
							bytes, _ := strconv.ParseInt(value, 10, 64)
							receiveBytes = uint64(bytes)
						} else if transmit == i {
							bytes, _ := strconv.ParseInt(value, 10, 64)
							transmitBytes = uint64(bytes)
						}
						i++
					}
				}
			}
		} else if strings.Contains(line, "face") {
			index := 0
			tag := false
			fields := strings.Split(line, "|")
			for _, field := range fields {
				if strings.Contains(field, "face") {
					index = 1
				} else if strings.Contains(field, "bytes") {
					values := strings.Fields(field)
					for _, value := range values {
						if strings.Contains(value, "bytes") {
							if !tag {
								tag = true
								receive = index
							} else {
								transmit = index
							}
						}
						index++
					}
				}
			}
		}
	}
	return ReceiveBytes(receiveBytes), TransmitBytes(transmitBytes), nil
}
func getNetIO() {
	ioCountStat, _ := net.IOCounters(true)
	for _, val := range ioCountStat {
		name := val.Name
		if strings.HasPrefix(name, "ens") {
			net_download, net_upload, err = UploadDownloadFlow(name)
			if err != nil {
				logrus.Println(err)
			}
			fmt.Println(net_download, net_upload)
		}
	}
	wg.Done()
}

func (rs *ResourceService) GetResource() (data map[string]map[string]interface{}) {
	wg.Add(4)
	go getCPU()
	go getMemory()
	go getDisk()
	go getNetIO()
	wg.Wait()
	data = map[string]map[string]interface{}{
		"cpu": map[string]interface{}{
			"cpu_core":    cpu_core,
			"cpu_utility": cpu_utility,
		},
		"memory": map[string]interface{}{
			"memory_total":         memory_total,
			"memory_free_percent":  memory_free_percent,
			"memory_used_percent":  memory_used_percent,
			"memory_other_percent": memory_other_percent,
		},
		"disk": map[string]interface{}{
			"disk_total":         disk_total,
			"disk_free_percent":  disk_free_percent,
			"disk_used_percent":  disk_used_percent,
			"disk_other_percent": disk_other_percent,
		},
		"net": map[string]interface{}{
			"net_download": net_download,
			"net_upload":   net_upload,
		},
	}
	fmt.Println(data)
	return
}

func (rs *ResourceService) GetCpu() (data map[string]interface{}) {
	wg.Add(1)
	go getCPU()
	wg.Wait()
	data = map[string]interface{}{
		"cpu_core":    cpu_core,
		"cpu_utility": cpu_utility,
	}
	return
}

func (rs *ResourceService) GetMemory() (data map[string]interface{}) {
	wg.Add(1)
	go getMemory()
	wg.Wait()
	data = map[string]interface{}{
		"memory_total":         memory_total,
		"memory_free_percent":  memory_free_percent,
		"memory_used_percent":  memory_used_percent,
		"memory_other_percent": memory_other_percent,
	}
	return
}

func (rs *ResourceService) GetDisk() (data map[string]interface{}) {
	wg.Add(1)
	go getDisk()
	wg.Wait()
	data = map[string]interface{}{
		"disk_total":         disk_total,
		"disk_free_percent":  disk_free_percent,
		"disk_used_percent":  disk_used_percent,
		"disk_other_percent": disk_other_percent,
	}
	return
}

func (rs *ResourceService) GetNetwork() (data map[string]interface{}) {
	wg.Add(1)
	go getNetIO()
	wg.Wait()
	data = map[string]interface{}{
		"net_download": net_download,
		"net_upload":   net_upload,
	}
	return
}
