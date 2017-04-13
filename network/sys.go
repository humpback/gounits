package network

import (
	"fmt"
	"net"
	"strconv"
)

// GrabSystemRandomIdlePort is exported
// Get system random idle port, auto bind
func GrabSystemRandomIdlePort() (uint32, error) {

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}

	defer l.Close()
	_, portstr, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return 0, err
	}

	port, err := strconv.Atoi(portstr)
	if err != nil {
		return 0, err
	}
	return uint32(port), nil
}

// GrabSystemRangeIdlePort is exported
func GrabSystemRangeIdlePort(kind string, minPort uint32, maxPort uint32) (uint32, error) {

	if minPort < 0 || maxPort < 0 || maxPort < minPort {
		return 0, fmt.Errorf("ports range invalid.")
	}

	if minPort > 65535 || maxPort > 65535 {
		return 0, fmt.Errorf("ports out of range.")
	}

	if minPort == 0 && maxPort == 0 {
		return GrabSystemRandomIdlePort()
	}

	localPorts := []uint32{}
	connectionStat, err := ConnectionStats(kind)
	if err != nil {
		return 0, err
	}

	for _, stat := range connectionStat {
		if stat.Status == "LISTEN" {
			localPorts = append(localPorts, stat.Laddr.Port)
		}
	}

	for port := minPort; port <= maxPort; port++ {
		found := false
		for _, localport := range localPorts {
			if localport == port {
				found = true
				break
			}
		}
		if !found {
			return port, nil
		}
	}
	return 0, fmt.Errorf("call not be allocated valid.")
}
