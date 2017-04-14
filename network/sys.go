package network

import "github.com/humpback/gounits/utils"

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
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

	ports := []uint32{}
	for port := minPort; port <= maxPort; port++ {
		if ret := utils.Contains(port, localPorts); !ret {
			ports = append(ports, port)
		}
	}

	if len(ports) == 0 {
		return 0, fmt.Errorf("call not be allocated valid.")
	}
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	return ports[r.Intn(len(ports)-1)], nil
}
