package compositor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"sync"
)

type NiriProvider struct {
	stateUpdateChan chan<- CompositorState
	done            chan struct{}
	wg              sync.WaitGroup
	mu              sync.Mutex
	currentState    CompositorState
}

func NewNiriProvider() *NiriProvider {
	return &NiriProvider{
		done: make(chan struct{}),
	}
}

func (p *NiriProvider) GetName() string {
	return "niri"
}

func (p *NiriProvider) Start(stateUpdateChan chan<- CompositorState) error {
	p.stateUpdateChan = stateUpdateChan
	p.refreshAll()

	p.wg.Add(1)
	go p.listenEvents()

	return nil
}

func (p *NiriProvider) Stop() {
	close(p.done)
	p.wg.Wait()
}

func (p *NiriProvider) Dispatch(command string) error {
	// Simple mapping for demonstration
	// In reality, hyprland commands need to be translated to niri commands
	// e.g. "workspace 1" -> "niri msg action focus-workspace 1"
	return nil
}

func (p *NiriProvider) listenEvents() {
	defer p.wg.Done()

	for {
		select {
		case <-p.done:
			return
		default:
			cmd := exec.Command("niri", "msg", "-j", "event-stream")
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				return
			}
			
			if err := cmd.Start(); err != nil {
				return
			}
			
			go func() {
				<-p.done
				cmd.Process.Kill()
			}()

			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				// We don't necessarily need to parse the precise event, just trigger a refresh
				// Niri events are full JSON objects
				p.refreshAll()
			}
			cmd.Wait()
		}
	}
}

func (p *NiriProvider) refreshAll() {
	p.mu.Lock()
	defer p.mu.Unlock()

	var state CompositorState

	// Outputs (Monitors)
	out, err := exec.Command("niri", "msg", "-j", "outputs").Output()
	if err == nil {
		var nOutputs map[string]struct {
			Name        string  `json:"name"`
			Make        string  `json:"make"`
			Model       string  `json:"model"`
			X           int     `json:"x"`
			Y           int     `json:"y"`
			Width       int     `json:"width"`
			Height      int     `json:"height"`
			Logical     struct {
				X      int `json:"x"`
				Y      int `json:"y"`
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"logical"`
			Scale       float64 `json:"scale"`
		}
		json.Unmarshal(out, &nOutputs)

		i := 0
		for _, no := range nOutputs {
			state.Monitors = append(state.Monitors, Monitor{
				ID:          i,
				Name:        no.Name,
				Description: fmt.Sprintf("%s %s", no.Make, no.Model),
				Width:       no.Logical.Width,
				Height:      no.Logical.Height,
				X:           no.Logical.X,
				Y:           no.Logical.Y,
				Scale:       no.Scale,
				RefreshRate: 60.0, // Niri doesn't always provide this in logical
			})
			i++
		}
	}

	// Workspaces
	out, err = exec.Command("niri", "msg", "-j", "workspaces").Output()
	if err == nil {
		var nWorkspaces []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Output    string `json:"output"`
			IsActive  bool   `json:"is_active"`
			IsFocused bool   `json:"is_focused"`
		}
		json.Unmarshal(out, &nWorkspaces)

		for _, nw := range nWorkspaces {
			var mId *int
			for _, m := range state.Monitors {
				if m.Name == nw.Output {
					mId = &m.ID
					break
				}
			}
			
			// Try to parse ID from name if it's numeric, or use 1-indexed
			id := nw.ID
			if idStr, err := strconv.Atoi(nw.Name); err == nil {
				id = idStr
			}

			state.Workspaces = append(state.Workspaces, Workspace{
				ID:        id,
				Name:      nw.Name,
				MonitorID: mId,
				Windows:   0, // Filled later
			})
			
			if nw.IsFocused {
				state.ActiveWorkspaceID = id
				if mId != nil {
					for i, m := range state.Monitors {
						if m.ID == *mId {
							state.Monitors[i].ActiveWorkspace = &id
							state.Monitors[i].Focused = true
						}
					}
				}
			}
		}
	}

	// Windows
	out, err = exec.Command("niri", "msg", "-j", "windows").Output()
	if err == nil {
		var nWindows []struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			AppID       string `json:"app_id"`
			WorkspaceID int    `json:"workspace_id"`
			IsFocused   bool   `json:"is_focused"`
			X           int    `json:"x"`
			Y           int    `json:"y"`
			Width       int    `json:"width"`
			Height      int    `json:"height"`
		}
		json.Unmarshal(out, &nWindows)

		for _, nw := range nWindows {
			// Find mapped workspace
			wsId := nw.WorkspaceID
			
			state.Windows = append(state.Windows, Window{
				Address:     strconv.Itoa(nw.ID),
				Mapped:      true,
				Hidden:      false,
				X:           nw.X,
				Y:           nw.Y,
				Width:       nw.Width,
				Height:      nw.Height,
				WorkspaceID: wsId,
				Title:       nw.Title,
				AppID:       nw.AppID,
				Focused:     nw.IsFocused,
			})
			
			for i, w := range state.Workspaces {
				if w.ID == wsId {
					state.Workspaces[i].Windows++
				}
			}
		}
	}

	p.currentState = state
	if p.stateUpdateChan != nil {
		p.stateUpdateChan <- state
	}
}
