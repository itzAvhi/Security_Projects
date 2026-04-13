package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

type failedLoginAttemptsData struct {
	Username string
	IP       string
}

func main() {
	var loginDatas []failedLoginAttemptsData
	cmd := exec.Command("journalctl", "_SYSTEMD_UNIT=sshd.service")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command 'journalctl'")
		return
	}

	ipcount := make(map[string]int)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Failed password") {
			fields := strings.Fields(line)
			var username, ip string
			for i, word := range fields {
				if word == "for" && i+1 < len(fields) {
					if fields[i+1] == "invalid" {
						if i+3 < len(fields) {
							username = fields[i+3]
						}
					} else {
						username = fields[i+1]
					}
				}
				if word == "from" && i+1 < len(fields) {
					ip = fields[i+1]
				}
			}
			if username != "" && ip != "" {
				ipcount[ip]++
				newUser := failedLoginAttemptsData{Username: username, IP: ip}
				loginDatas = append(loginDatas, newUser)
			}
		}
	}

	fmt.Println("Failed login attempts")
	fmt.Printf("Total failed attempts: %d\n", len(loginDatas))
	fmt.Println("\nBreakdown by IP:")
	for ip, count := range ipcount {
		fmt.Printf("  %s: %d attempts\n", ip, count)
	}
	fmt.Println("\nDetailed attempts:")
	for _, data := range loginDatas {
		fmt.Printf("  User: %s, IP: %s\n", data.Username, data.IP)
	}
}
