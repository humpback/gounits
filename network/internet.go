/*
* (C) 2001-2017 humpback Inc.
*
* gounits source code
* version: 1.0.0
* author: bobliu0909@gmail.com
* datetime: 2015-10-14
*
 */

package network

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

var (
	ErrIpAddrEmpty      = errors.New("local ipaddr is empty.")
	ErrDriveNameInvalid = errors.New("drivename invalid.")
)

func NetAddrItoa(ip uint32) net.IP {

	var bytes [4]byte
	bytes[0] = byte(ip & 0xFF)
	bytes[1] = byte((ip >> 8) & 0xFF)
	bytes[2] = byte((ip >> 16) & 0xFF)
	bytes[3] = byte((ip >> 24) & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func NetAddrAtoi(ip net.IP) uint32 {

	bits := strings.Split(ip.String(), ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var dwip uint32
	dwip += uint32(b0) << 24
	dwip += uint32(b1) << 16
	dwip += uint32(b2) << 8
	dwip += uint32(b3)
	return dwip
}

type AddrInfo struct {
	Name string
	Mac  string
	IP   string
}

func GetLocalNetDriveInfo(name string) (*AddrInfo, error) {

	addrs, err := GetLocalNetAddrs()
	if err != nil {
		return nil, err
	}

	if len(addrs) == 0 {
		if name == "" {
			ip := GetDefaultIP()
			if ip == "" {
				return nil, ErrIpAddrEmpty
			}
			return &AddrInfo{
				Name: "",
				Mac:  "",
				IP:   ip,
			}, nil
		} else {
			return nil, ErrIpAddrEmpty
		}
	}

	var addr *AddrInfo = nil
	if name == "" {
		addr = addrs[0] //未传入设备信息，默认获取第一张网卡
	} else {
		for _, value := range addrs {
			if value.Name == name {
				addr = value //根据设备名称获取网卡信息
				break
			}
		}
		if addr == nil {
			return nil, ErrDriveNameInvalid
		}
	}
	return addr, nil
}

func GetLocalNetAddrs() ([]*AddrInfo, error) {

	ips := make([]*AddrInfo, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		return ips, err
	}

	for _, iface := range ifaces {
		name := iface.Name
		mac := iface.HardwareAddr
		if iface.Flags&net.FlagUp == 0 {
			continue // 忽略禁用的网卡
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // 忽略loopback回路接口
		}

		// 忽略 docker与网桥
		if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "w-") {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return ips, err
		}

		for _, addr := range addrs {

			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // 不是ipv4地址，放弃
			}

			ipStr := ip.String()
			if IsIntranet(ipStr) {
				addrinfo := &AddrInfo{
					Name: name,
					Mac:  mac.String(),
					IP:   ipStr,
				}
				ips = append(ips, addrinfo)
			}
		}
	}
	return ips, nil
}

func GetDefaultIP() string {

	conn, err := net.Dial("udp", "192.168.0.1:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

func IsIntranet(ipStr string) bool {

	if strings.HasPrefix(ipStr, "10.") || strings.HasPrefix(ipStr, "192.168.") {
		return true
	}

	if strings.HasPrefix(ipStr, "172.") {
		// 172.16.0.0-172.31.255.255
		arr := strings.Split(ipStr, ".")
		if len(arr) != 4 {
			return false
		}

		second, err := strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			return false
		}

		if second >= 16 && second <= 31 {
			return true
		}
	}
	return false
}
