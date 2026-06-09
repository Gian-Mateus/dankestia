package sysinfo

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/AvengeMedia/Dankestia/core/pkg/syncmap"
)

type Manager struct {
	state       State
	stateMutex  sync.RWMutex
	subscribers syncmap.Map[string, chan State]
	stopChan    chan struct{}

	lastIdle  float64
	lastTotal float64
}

func NewManager() (*Manager, error) {
	m := &Manager{
		stopChan: make(chan struct{}),
	}
	m.state.Cpu.Name = "CPU"
	m.state.OsName = m.readOsName()
	m.state.KernelVersion = m.readKernelVersion()
	go m.loop()
	return m, nil
}

func (m *Manager) Subscribe(clientID string) chan State {
	ch := make(chan State, 10)
	m.subscribers.Store(clientID, ch)
	return ch
}

func (m *Manager) Unsubscribe(clientID string) {
	if ch, ok := m.subscribers.Load(clientID); ok {
		close(ch)
		m.subscribers.Delete(clientID)
	}
}

func (m *Manager) GetState() State {
	m.stateMutex.RLock()
	defer m.stateMutex.RUnlock()
	return m.state
}

func (m *Manager) loop() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// Initial fetch
	m.updateStats()

	for {
		select {
		case <-ticker.C:
			m.updateStats()
		case <-m.stopChan:
			return
		}
	}
}

func (m *Manager) updateStats() {
	cpuUsage := m.readCpuUsage()
	memTotal, memAvail := m.readMemory()
	uptime := m.readUptime()
	storageTotal, storageFree := m.readStorage()

	m.stateMutex.Lock()
	m.state.Cpu.Percentage = cpuUsage
	m.state.Memory.TotalMB = memTotal
	m.state.Memory.AvailableMB = memAvail
	m.state.Memory.UsedMB = memTotal - memAvail
	m.state.UptimeSeconds = uptime
	m.state.Storage.TotalMB = storageTotal
	m.state.Storage.FreeMB = storageFree
	m.state.Storage.UsedMB = storageTotal - storageFree
	newState := m.state
	m.stateMutex.Unlock()

	m.subscribers.Range(func(key string, ch chan State) bool {
		select {
		case ch <- newState:
		default:
		}
		return true
	})
}

func (m *Manager) readCpuUsage() float64 {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > 4 && fields[0] == "cpu" {
			var total float64
			var idle float64
			for i := 1; i < len(fields); i++ {
				val, _ := strconv.ParseFloat(fields[i], 64)
				total += val
				if i == 4 || i == 5 { // idle and iowait
					idle += val
				}
			}

			diffIdle := idle - m.lastIdle
			diffTotal := total - m.lastTotal
			var usage float64
			if diffTotal > 0 {
				usage = (1000 * (diffTotal - diffIdle) / diffTotal + 5) / 10
			}

			m.lastIdle = idle
			m.lastTotal = total
			return usage
		}
	}
	return 0
}

func (m *Manager) readMemory() (float64, float64) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0
	}
	defer file.Close()

	var total, avail float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				val, _ := strconv.ParseFloat(fields[1], 64)
				total = val / 1024 // KB to MB
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				val, _ := strconv.ParseFloat(fields[1], 64)
				avail = val / 1024 // KB to MB
			}
		}
	}
	return total, avail
}

func (m *Manager) readOsName() string {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "Linux"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
		}
	}
	return "Linux"
}

func (m *Manager) readKernelVersion() string {
	var uts syscall.Utsname
	if err := syscall.Uname(&uts); err == nil {
		var buf []byte
		for _, b := range uts.Release {
			if b == 0 {
				break
			}
			buf = append(buf, byte(b))
		}
		return string(buf)
	}
	return "Unknown"
}

func (m *Manager) readUptime() float64 {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	fields := strings.Fields(string(data))
	if len(fields) > 0 {
		uptime, _ := strconv.ParseFloat(fields[0], 64)
		return uptime
	}
	return 0
}

func (m *Manager) readStorage() (float64, float64) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		return 0, 0
	}
	total := float64(stat.Blocks) * float64(stat.Bsize) / (1024 * 1024)
	free := float64(stat.Bavail) * float64(stat.Bsize) / (1024 * 1024)
	return total, free
}
