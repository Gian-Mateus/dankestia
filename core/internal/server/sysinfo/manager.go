package sysinfo

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"
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

	m.stateMutex.Lock()
	m.state.Cpu.Percentage = cpuUsage
	m.state.Memory.TotalMB = memTotal
	m.state.Memory.AvailableMB = memAvail
	m.state.Memory.UsedMB = memTotal - memAvail
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
