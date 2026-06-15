package compositor

// Monitor represents a physical display
type Monitor struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description,omitempty"`
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	RefreshRate     float64 `json:"refreshRate"`
	X               int     `json:"x"`
	Y               int     `json:"y"`
	Scale           float64 `json:"scale"`
	ActiveWorkspace *int    `json:"activeWorkspaceId,omitempty"`
	Focused         bool    `json:"focused"`
}

// Workspace represents a virtual desktop
type Workspace struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	MonitorID *int   `json:"monitorId,omitempty"`
	Windows   int    `json:"windows"`
	HasFullscreen bool `json:"hasFullscreen"`
	IsSpecial bool   `json:"isSpecial"`
}

// Window represents an open application window
type Window struct {
	Address     string `json:"address"` // Unique identifier (hex address for Hyprland, integer ID for Niri)
	Mapped      bool   `json:"mapped"`
	Hidden      bool   `json:"hidden"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	WorkspaceID int    `json:"workspaceId"`
	MonitorID   int    `json:"monitorId,omitempty"`
	Floating    bool   `json:"floating"`
	Fullscreen  bool   `json:"fullscreen"`
	Pinned      bool   `json:"pinned"`
	Title       string `json:"title"`
	AppID       string `json:"appId"` // "class" in Hyprland, "app_id" in Niri
	Focused     bool   `json:"focused"`
}

// CompositorState represents the full state broadcast to the frontend
type CompositorState struct {
	Monitors   []Monitor   `json:"monitors"`
	Workspaces []Workspace `json:"workspaces"`
	Windows    []Window    `json:"windows"`
	ActiveWorkspaceID int  `json:"activeWorkspaceId"`
}

// Provider represents a backend source for compositor data (Hyprland, Niri, Sway)
type Provider interface {
	Start(stateUpdateChan chan<- CompositorState) error
	Stop()
	Dispatch(command string) error
	GetName() string
}
