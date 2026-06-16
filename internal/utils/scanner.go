package utils 

import (
	"fmt"
	"net"
)

func ScanPorts(startRange, endRange int) {
	for port := startRange; port <= endRange; port++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			continue
		}
		conn.Close()
		fmt.Printf("Port %d is open\n", port)
	}
}

func ScanSinglePort(port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Printf("Port %d is closed\n", port)
		return
	}
	conn.Close()
	fmt.Printf("Port %d is open\n", port)
}