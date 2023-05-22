package main

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type ReceiveBytes uint64
type TransmitBytes uint64

type DownloadFlow string
type UploadFlow string

func UploadDownloadFlow(dev string) (DownloadFlow, UploadFlow, error) {
	down, up, err := TotalFlowByDevice(dev)
	if err != nil {
		logrus.Println(err)
		return "", "", err
	}
	logrus.Debug(down, up)
	time.Sleep(time.Second * 1)
	down2, up2, err := TotalFlowByDevice(dev)
	if err != nil {
		logrus.Println(err)
		return "", "", err
	}
	logrus.Debug(down2, up2)
	downStr := strconv.FormatInt(int64((down2-down)/1024), 10) + "Kbps"
	upStr := strconv.FormatInt(int64((up2-up)/1024), 10) + "Kbps"
	return DownloadFlow(downStr), UploadFlow(upStr), nil
}

func TotalFlowByDevice(dev string) (ReceiveBytes, TransmitBytes, error) {
	devInfo, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		logrus.Println(err)
		return 0, 0, err
	}
	var receive int = -1
	var transmit int = -1
	var receiveBytes uint64
	var transmitBytes uint64
	lines := strings.Split(string(devInfo), "\n")
	for _, line := range lines {
		logrus.Debug(line)
		if strings.Contains(line, dev) {
			i := 0
			fields := strings.Split(line, ":")
			for _, field := range fields {
				if strings.Contains(field, dev) {
					i = 1
				} else {
					values := strings.Fields(field)
					for _, value := range values {
						//logrus.Debug(value)
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
						//logrus.Debug(value)
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
	logrus.Debug("receive_bytes :", receiveBytes)
	logrus.Debug("transmit_bytes :", transmitBytes)
	return ReceiveBytes(receiveBytes), TransmitBytes(transmitBytes), nil
}
