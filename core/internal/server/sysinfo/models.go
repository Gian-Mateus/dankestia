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

type StorageState struct {
	TotalMB float64 `json:"totalMB"`
	UsedMB  float64 `json:"usedMB"`
	FreeMB  float64 `json:"freeMB"`
}

type State struct {
	OsName        string       `json:"osName"`
	KernelVersion string       `json:"kernelVersion"`
	UptimeSeconds float64      `json:"uptimeSeconds"`
	Cpu           CpuState     `json:"cpu"`
	Memory        MemoryState  `json:"memory"`
	Storage       StorageState `json:"storage"`
}
