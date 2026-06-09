package sysinfo

type CpuState struct {
	Name        string  `json:"name"`
	Percentage  float64 `json:"percentage"`
	Temperature float64 `json:"temperature"`
}

type MemoryState struct {
	TotalMB     float64 `json:"totalMB"`
	UsedMB      float64 `json:"usedMB"`
	AvailableMB float64 `json:"availableMB"`
}

type State struct {
	Cpu    CpuState    `json:"cpu"`
	Memory MemoryState `json:"memory"`
}
