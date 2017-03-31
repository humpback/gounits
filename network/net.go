package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"syscall"
)

var (
	ErrIpAddrEmpty      = errors.New("local ipaddr is empty.")
	ErrDriveNameInvalid = errors.New("drivename invalid.")
)

// Addr is implemented compatibility to psutil
type Addr struct {
	IP   string `json:"ip"`
	Port uint32 `json:"port"`
}

type AddrInfo struct {
	Name string `json:"name"`
	Mac  string `json:"mac"`
	IP   string `json:"ip"`
}

type ConnectionStat struct {
	Fd     uint32  `json:"fd"`
	Family uint32  `json:"family"`
	Type   uint32  `json:"type"`
	Laddr  Addr    `json:"localaddr"`
	Raddr  Addr    `json:"remoteaddr"`
	Status string  `json:"status"`
	Uids   []int32 `json:"uids"`
	Pid    int32   `json:"pid"`
}

// NetInterfaceAddr is designed for represent interface addresses
type InterfaceAddr struct {
	Addr string `json:"addr"`
}

type InterfaceStat struct {
	MTU          int             `json:"mtu"`          // maximum transmission unit
	Name         string          `json:"name"`         // e.g., "en0", "lo0", "eth0.100"
	HardwareAddr string          `json:"hardwareaddr"` // IEEE MAC-48, EUI-48 and EUI-64 form
	Flags        []string        `json:"flags"`        // e.g., FlagUp, FlagLoopback, FlagMulticast
	Addrs        []InterfaceAddr `json:"addrs"`
}

var constMap = map[string]int{
	"TCP":  syscall.SOCK_STREAM,
	"UDP":  syscall.SOCK_DGRAM,
	"IPv4": syscall.AF_INET,
	"IPv6": syscall.AF_INET6,
}

func (n ConnectionStat) String() string {
	s, _ := json.Marshal(n)
	return string(s)
}

func (a Addr) String() string {
	s, _ := json.Marshal(a)
	return string(s)
}

func (n InterfaceStat) String() string {
	s, _ := json.Marshal(n)
	return string(s)
}

func (n InterfaceAddr) String() string {
	s, _ := json.Marshal(n)
	return string(s)
}

func ConvertToIPStr(ip uint32) string {

	var bytes [4]byte
	bytes[0] = byte(ip & 0xFF)
	bytes[1] = byte((ip >> 8) & 0xFF)
	bytes[2] = byte((ip >> 16) & 0xFF)
	bytes[3] = byte((ip >> 24) & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}

func ConvertToIPInt(ip string) uint32 {

	netip := net.ParseIP(ip)
	bits := strings.Split(netip.String(), ".")
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

// Interfaces is exported
func Interfaces() ([]InterfaceStat, error) {
	is, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	ret := make([]InterfaceStat, 0, len(is))
	for _, ifi := range is {

		var flags []string
		if ifi.Flags&net.FlagUp != 0 {
			flags = append(flags, "up")
		}
		if ifi.Flags&net.FlagBroadcast != 0 {
			flags = append(flags, "broadcast")
		}
		if ifi.Flags&net.FlagLoopback != 0 {
			flags = append(flags, "loopback")
		}
		if ifi.Flags&net.FlagPointToPoint != 0 {
			flags = append(flags, "pointtopoint")
		}
		if ifi.Flags&net.FlagMulticast != 0 {
			flags = append(flags, "multicast")
		}

		r := InterfaceStat{
			Name:         ifi.Name,
			MTU:          ifi.MTU,
			HardwareAddr: ifi.HardwareAddr.String(),
			Flags:        flags,
		}
		addrs, err := ifi.Addrs()
		if err == nil {
			r.Addrs = make([]InterfaceAddr, 0, len(addrs))
			for _, addr := range addrs {
				r.Addrs = append(r.Addrs, InterfaceAddr{
					Addr: addr.String(),
				})
			}

		}
		ret = append(ret, r)
	}

	return ret, nil
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

func parseNetLine(line string) (ConnectionStat, error) {
	f := strings.Fields(line)
	if len(f) < 9 {
		return ConnectionStat{}, fmt.Errorf("wrong line,%s", line)
	}

	pid, err := strconv.Atoi(f[1])
	if err != nil {
		return ConnectionStat{}, err
	}
	fd, err := strconv.Atoi(strings.Trim(f[3], "u"))
	if err != nil {
		return ConnectionStat{}, fmt.Errorf("unknown fd, %s", f[3])
	}
	netFamily, ok := constMap[f[4]]
	if !ok {
		return ConnectionStat{}, fmt.Errorf("unknown family, %s", f[4])
	}
	netType, ok := constMap[f[7]]
	if !ok {
		return ConnectionStat{}, fmt.Errorf("unknown type, %s", f[7])
	}

	laddr, raddr, err := parseNetAddr(f[8])
	if err != nil {
		return ConnectionStat{}, fmt.Errorf("failed to parse netaddr, %s", f[8])
	}

	n := ConnectionStat{
		Fd:     uint32(fd),
		Family: uint32(netFamily),
		Type:   uint32(netType),
		Laddr:  laddr,
		Raddr:  raddr,
		Pid:    int32(pid),
	}
	if len(f) == 10 {
		n.Status = strings.Trim(f[9], "()")
	}

	return n, nil
}

func parseNetAddr(line string) (laddr Addr, raddr Addr, err error) {
	parse := func(l string) (Addr, error) {
		host, port, err := net.SplitHostPort(l)
		if err != nil {
			return Addr{}, fmt.Errorf("wrong addr, %s", l)
		}
		lport, err := strconv.Atoi(port)
		if err != nil {
			return Addr{}, err
		}
		return Addr{IP: host, Port: uint32(lport)}, nil
	}

	addrs := strings.Split(line, "->")
	if len(addrs) == 0 {
		return laddr, raddr, fmt.Errorf("wrong netaddr, %s", line)
	}
	laddr, err = parse(addrs[0])
	if len(addrs) == 2 { // remote addr exists
		raddr, err = parse(addrs[1])
		if err != nil {
			return laddr, raddr, err
		}
	}

	return laddr, raddr, err
}
