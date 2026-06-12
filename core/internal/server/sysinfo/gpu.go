package sysinfo

import (
	"bytes"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type gpuType int

const (
	gpuTypeAuto gpuType = iota
	gpuTypeNone
	gpuTypeNvidia
	gpuTypeGeneric
)

func detectGpuType() gpuType {
	// Nvidia detection
	if err := exec.Command("sh", "-c", "command -v nvidia-smi >/dev/null 2>&1 && nvidia-smi -L >/dev/null 2>&1").Run(); err == nil {
		return gpuTypeNvidia
	}
	// Generic (AMD/Intel) detection
	if err := exec.Command("sh", "-c", "ls /sys/class/drm/card*/device/gpu_busy_percent 2>/dev/null | grep -q .").Run(); err == nil {
		return gpuTypeGeneric
	}
	return gpuTypeNone
}

func detectGpuName() string {
	// Try Nvidia
	out, err := exec.Command("sh", "-c", "nvidia-smi --query-gpu=name --format=csv,noheader 2>/dev/null").Output()
	if err == nil && len(bytes.TrimSpace(out)) > 0 {
		return cleanGpuName(string(out))
	}
	// Try glxinfo
	out, err = exec.Command("sh", "-c", "glxinfo -B 2>/dev/null | grep 'Device:' | cut -d':' -f2 | cut -d'(' -f1").Output()
	if err == nil && len(bytes.TrimSpace(out)) > 0 {
		return cleanGpuName(string(out))
	}
	// Try lspci
	out, err = exec.Command("sh", "-c", "lspci 2>/dev/null | grep -i 'vga\\|3d controller\\|display' | head -1").Output()
	if err == nil && len(bytes.TrimSpace(out)) > 0 {
		return cleanGpuName(string(out))
	}
	return "Unknown GPU"
}

func cleanGpuName(s string) string {
	s = strings.TrimSpace(s)
	// Try matching brackets e.g. "VGA compatible controller: Intel Corporation [UHD Graphics 620]"
	reBracket := regexp.MustCompile(`\[([^\]]+)\][^\[]*$`)
	match := reBracket.FindStringSubmatch(s)
	if len(match) > 1 {
		s = match[1]
	} else {
		// Match colon e.g. "VGA compatible controller: Intel Graphics"
		reColon := regexp.MustCompile(`:\s*(.+)`)
		match = reColon.FindStringSubmatch(s)
		if len(match) > 1 {
			s = match[1]
		}
	}
	// Remove noise
	s = strings.ReplaceAll(s, "(R)", "")
	s = strings.ReplaceAll(s, "(TM)", "")
	s = strings.ReplaceAll(s, "Graphics", "")
	reSpaces := regexp.MustCompile(`\s+`)
	s = reSpaces.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func readNvidiaStats() (float64, float64) {
	out, err := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu,temperature.gpu", "--format=csv,noheader,nounits").Output()
	if err != nil {
		return 0, 0
	}
	parts := strings.Split(strings.TrimSpace(string(out)), ",")
	if len(parts) < 2 {
		return 0, 0
	}
	usage, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	temp, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	return usage / 100.0, temp
}

func readGenericUsage() float64 {
	matches, err := filepath.Glob("/sys/class/drm/card*/device/gpu_busy_percent")
	if err != nil || len(matches) == 0 {
		return 0
	}
	var sum float64
	var count int
	for _, match := range matches {
		data, err := os.ReadFile(match)
		if err == nil {
			v, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
			if err == nil {
				sum += v
				count++
			}
		}
	}
	if count > 0 {
		return sum / float64(count) / 100.0
	}
	return 0
}

func readGenericTemp() float64 {
	// Attempt to read from hwmon for amdgpu
	matches, err := filepath.Glob("/sys/class/hwmon/hwmon*/name")
	if err != nil {
		return 0
	}
	for _, match := range matches {
		data, err := os.ReadFile(match)
		if err == nil && strings.TrimSpace(string(data)) == "amdgpu" {
			dir := filepath.Dir(match)
			tempFiles, _ := filepath.Glob(filepath.Join(dir, "temp*_input"))
			if len(tempFiles) > 0 {
				tempData, err := os.ReadFile(tempFiles[0])
				if err == nil {
					v, err := strconv.ParseFloat(strings.TrimSpace(string(tempData)), 64)
					if err == nil {
						return v / 1000.0
					}
				}
			}
		}
	}
	return 0
}

func (m *Manager) updateGpuStats() {
	if m.gpuType == gpuTypeNone || m.gpuType == gpuTypeAuto {
		return
	}
	var usage float64
	var temp float64

	if m.gpuType == gpuTypeNvidia {
		usage, temp = readNvidiaStats()
	} else if m.gpuType == gpuTypeGeneric {
		usage = readGenericUsage()
		temp = readGenericTemp()
	}

	m.stateMutex.Lock()
	if math.Abs(usage-m.state.Gpu.Percentage) > 0.0001 {
		m.state.Gpu.Percentage = usage
	}
	if math.Abs(temp-m.state.Gpu.Temperature) > 0.05 {
		m.state.Gpu.Temperature = temp
	}
	m.stateMutex.Unlock()
}
