package utils

import (
	"encoding/csv"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// PortInfo contains status and process information for a specific port
type PortInfo struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Process  string `json:"process"`
	PID      int    `json:"pid"`
	State    string `json:"state"`
}

// ScanListeningPorts fetches all active listening TCP ports and their owning process information.
func ScanListeningPorts() (map[int]PortInfo, error) {
	if runtime.GOOS == "windows" {
		return scanListeningPortsWindows()
	}
	return scanListeningPortsUnix()
}

// KillProcess forcefully terminates a process by its PID
func KillProcess(pid int) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
	} else {
		cmd = exec.Command("kill", "-9", strconv.Itoa(pid))
	}
	return cmd.Run()
}

// ScanPorts scans a range of ports (legacy function, kept for compatibility)
func ScanPorts(startRange, endRange int) {
	ports, err := ScanListeningPorts()
	if err != nil {
		// Fallback to basic dialing if commands fail
		for port := startRange; port <= endRange; port++ {
			conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
			if err != nil {
				continue
			}
			conn.Close()
			fmt.Printf("Port %d is open\n", port)
		}
		return
	}

	for port := startRange; port <= endRange; port++ {
		if info, ok := ports[port]; ok {
			fmt.Printf("Port %d is open (Process: %s, PID: %d)\n", port, info.Process, info.PID)
		}
	}
}

// ScanSinglePort scans a single port (legacy function, kept for compatibility)
func ScanSinglePort(port int) {
	ports, err := ScanListeningPorts()
	if err != nil {
		// Fallback to basic dialing
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			fmt.Printf("Port %d is closed\n", port)
			return
		}
		conn.Close()
		fmt.Printf("Port %d is open\n", port)
		return
	}

	if info, ok := ports[port]; ok {
		fmt.Printf("Port %d is open (Process: %s, PID: %d)\n", port, info.Process, info.PID)
	} else {
		fmt.Printf("Port %d is closed\n", port)
	}
}

// getProcessMapWindows queries tasklist to map PIDs to process names in bulk
func getProcessMapWindows() (map[int]string, error) {
	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	pidMap := make(map[int]string)
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		r := csv.NewReader(strings.NewReader(line))
		records, err := r.Read()
		if err != nil || len(records) < 2 {
			continue
		}
		name := records[0]
		pidStr := records[1]
		pid, err := strconv.Atoi(pidStr)
		if err == nil {
			pidMap[pid] = name
		}
	}
	return pidMap, nil
}

// scanListeningPortsWindows parses netstat and correlates with tasklist
func scanListeningPortsWindows() (map[int]PortInfo, error) {
	pids, err := getProcessMapWindows()
	if err != nil {
		pids = make(map[int]string)
	}

	cmd := exec.Command("netstat", "-ano", "-p", "tcp")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	portMap := make(map[int]PortInfo)
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, "LISTENING") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		localAddr := fields[1]
		pidStr := fields[4]

		lastColon := strings.LastIndex(localAddr, ":")
		if lastColon == -1 {
			continue
		}
		portStr := localAddr[lastColon+1:]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		procName := "Unknown"
		if name, ok := pids[pid]; ok {
			procName = name
		} else if pid == 0 {
			procName = "System Idle"
		} else if pid == 4 {
			procName = "System"
		}

		portMap[port] = PortInfo{
			Port:     port,
			Protocol: "TCP",
			Process:  procName,
			PID:      pid,
			State:    "LISTENING",
		}
	}
	return portMap, nil
}

// scanListeningPortsUnix parses lsof for Unix-like environments
func scanListeningPortsUnix() (map[int]PortInfo, error) {
	cmd := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-n", "-P")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	portMap := make(map[int]PortInfo)
	lines := strings.Split(string(out), "\n")
	if len(lines) <= 1 {
		return portMap, nil
	}

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		procName := fields[0]
		pidStr := fields[1]
		nameField := fields[8]

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		lastColon := strings.LastIndex(nameField, ":")
		if lastColon == -1 {
			continue
		}
		portStr := nameField[lastColon+1:]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		portMap[port] = PortInfo{
			Port:     port,
			Protocol: "TCP",
			Process:  procName,
			PID:      pid,
			State:    "LISTENING",
		}
	}
	return portMap, nil
}