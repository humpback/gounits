package network

import (
	"net"
	"strconv"
)

// Get system idle port
func GrabSystemIdlePort() (uint16, error) {

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
	return uint16(port), nil
}
