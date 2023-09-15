package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const PowerPath = "/sys/class/power_supply/"

type Battery struct {
	now      float64
	full     float64
	charging bool
}

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	b, err := getBattery()
	if err != nil {
		log.Fatalln(err)
	}

	var s strings.Builder
	s.WriteString(strconv.Itoa(int(b.now / b.full * 100)))
	if b.charging {
		s.WriteString("+")
	} else {
		s.WriteString("-")
	}
	fmt.Println(s.String())
}

func getBattery() (Battery, error) {
	var b Battery
	path, err := firstBattery()
	if err != nil {
		return b, err
	}

	now, full, status, err := readStats(path)
	if err != nil {
		return b, err
	}
	return parseStats(now, full, status)
}

// firstBattery finds the first battery and returns its path or an error.
func firstBattery() (string, error) {
	entries, err := os.ReadDir(PowerPath)
	if err != nil {
		return "", err
	}

	var path string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "BAT") {
			path = filepath.Join(PowerPath, entry.Name())
			break
		}
	}
	if path == "" {
		return "", fmt.Errorf("no battery found")
	}
	return path, nil
}

// readStats from a battery using either its charge or energy properies.
// Now, full, and status or an error are returned.
func readStats(path string) (string, string, string, error) {
	var now, full, status string
	entries, err := os.ReadDir(path)
	if err != nil {
		return "", "", "", err
	}
	for _, entry := range entries {
		switch entry.Name() {
		case "energy_now", "charge_now":
			data, err := os.ReadFile(filepath.Join(path, entry.Name()))
			if err != nil {
				return "", "", "", err
			}
			now = string(data)
		case "energy_full", "charge_full":
			data, err := os.ReadFile(filepath.Join(path, entry.Name()))
			if err != nil {
				return "", "", "", err
			}
			full = string(data)
		case "status":
			data, err := os.ReadFile(filepath.Join(path, entry.Name()))
			if err != nil {
				return "", "", "", err
			}
			status = string(data)
		}
	}
	return now, full, status, nil
}

// parseStats reads the string properties and parses them into a Battery
// object.
func parseStats(now, full, status string) (Battery, error) {
	var b Battery
	n, err := strconv.ParseFloat(strings.TrimSpace(now), 64)
	if err != nil {
		return b, err
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(full), 64)
	if err != nil {
		return b, err
	}

	b.now = n
	b.full = f

	if status == "Charging" {
		b.charging = true
	}

	return b, nil
}
