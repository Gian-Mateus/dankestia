package compositor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type HyprlandProvider struct {
	stateUpdateChan chan<- CompositorState
	done            chan struct{}
	wg              sync.WaitGroup
	mu              sync.Mutex
	currentState    CompositorState
}

func NewHyprlandProvider() *HyprlandProvider {
	return &HyprlandProvider{
		done: make(chan struct{}),
	}
}

func (p *HyprlandProvider) GetName() string {
	return "hyprland"
}

func (p *HyprlandProvider) Start(stateUpdateChan chan<- CompositorState) error {
	p.stateUpdateChan = stateUpdateChan
	
	// Initial fetch
	p.refreshAll()

	p.wg.Add(1)
	go p.listenEvents()

	return nil
}

func (p *HyprlandProvider) Stop() {
	close(p.done)
	p.wg.Wait()
}

func (p *HyprlandProvider) Dispatch(command string) error {
	cmd := exec.Command("hyprctl", "dispatch", command)
	return cmd.Run()
}

func (p *HyprlandProvider) listenEvents() {
	defer p.wg.Done()

	sig := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if sig == "" {
		fmt.Println("HYPRLAND_INSTANCE_SIGNATURE not set")
		return
	}

	xdg := os.Getenv("XDG_RUNTIME_DIR")
	if xdg == "" {
		xdg = fmt.Sprintf("/run/user/%d", os.Getuid())
	}

	sockPath := filepath.Join(xdg, "hypr", sig, ".socket2.sock")

	var conn net.Conn
	var err error

	for {
		select {
		case <-p.done:
			if conn != nil {
				conn.Close()
			}
			return
		default:
			if conn == nil {
				conn, err = net.Dial("unix", sockPath)
				if err != nil {
					time.Sleep(2 * time.Second)
					continue
				}
			}

			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				line := scanner.Text()
				parts := strings.SplitN(line, ">>", 2)
				if len(parts) == 2 {
					event := parts[0]
					// args := parts[1]
					
					// We just refresh everything on relevant events for simplicity and correctness
					// More optimized implementations would just update the specific part of the state
					switch event {
					case "workspace", "focusedmon", "activewindow", "windowtitle", "openwindow", "closewindow", "movewindow", "monitoradded", "monitorremoved":
						p.refreshAll()
					}
				}
			}

			// If we exit the loop, the connection died
			conn.Close()
			conn = nil
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *HyprlandProvider) refreshAll() {
	p.mu.Lock()
	defer p.mu.Unlock()

	var state CompositorState

	// Monitors
	out, err := exec.Command("hyprctl", "-j", "monitors").Output()
	if err == nil {
		var hMonitors []struct {
			ID              int     `json:"id"`
			Name            string  `json:"name"`
			Description     string  `json:"description"`
			Width           int     `json:"width"`
			Height          int     `json:"height"`
			RefreshRate     float64 `json:"refreshRate"`
			X               int     `json:"x"`
			Y               int     `json:"y"`
			Scale           float64 `json:"scale"`
			Focused         bool    `json:"focused"`
			ActiveWorkspace struct {
				ID int `json:"id"`
			} `json:"activeWorkspace"`
		}
		json.Unmarshal(out, &hMonitors)

		for _, hm := range hMonitors {
			activeWs := hm.ActiveWorkspace.ID
			state.Monitors = append(state.Monitors, Monitor{
				ID:              hm.ID,
				Name:            hm.Name,
				Description:     hm.Description,
				Width:           hm.Width,
				Height:          hm.Height,
				RefreshRate:     hm.RefreshRate,
				X:               hm.X,
				Y:               hm.Y,
				Scale:           hm.Scale,
				ActiveWorkspace: &activeWs,
				Focused:         hm.Focused,
			})
		}
	}

	// Workspaces
	out, err = exec.Command("hyprctl", "-j", "workspaces").Output()
	if err == nil {
		var hWorkspaces []struct {
			ID           int    `json:"id"`
			Name         string `json:"name"`
			MonitorID    int    `json:"monitorID"`
			Windows      int    `json:"windows"`
			HasFullscreen bool  `json:"hasfullscreen"`
		}
		json.Unmarshal(out, &hWorkspaces)

		for _, hw := range hWorkspaces {
			mId := hw.MonitorID
			state.Workspaces = append(state.Workspaces, Workspace{
				ID:            hw.ID,
				Name:          hw.Name,
				MonitorID:     &mId,
				Windows:       hw.Windows,
				HasFullscreen: hw.HasFullscreen,
				IsSpecial:     strings.HasPrefix(hw.Name, "special:"),
			})
		}
	}

	// Windows (Clients)
	out, err = exec.Command("hyprctl", "-j", "clients").Output()
	if err == nil {
		var hClients []struct {
			Address   string `json:"address"`
			Mapped    bool   `json:"mapped"`
			Hidden    bool   `json:"hidden"`
			At        []int  `json:"at"`
			Size      []int  `json:"size"`
			Workspace struct {
				ID int `json:"id"`
			} `json:"workspace"`
			Floating   bool   `json:"floating"`
			Monitor    int    `json:"monitor"`
			Class      string `json:"class"`
			Title      string `json:"title"`
			Pinned     bool   `json:"pinned"`
			Fullscreen int    `json:"fullscreen"`
		}
		json.Unmarshal(out, &hClients)

		for _, hc := range hClients {
			x, y, w, h := 0, 0, 0, 0
			if len(hc.At) == 2 {
				x, y = hc.At[0], hc.At[1]
			}
			if len(hc.Size) == 2 {
				w, h = hc.Size[0], hc.Size[1]
			}
			state.Windows = append(state.Windows, Window{
				Address:     hc.Address,
				Mapped:      hc.Mapped,
				Hidden:      hc.Hidden,
				X:           x,
				Y:           y,
				Width:       w,
				Height:      h,
				WorkspaceID: hc.Workspace.ID,
				MonitorID:   hc.Monitor,
				Floating:    hc.Floating,
				Fullscreen:  hc.Fullscreen != 0,
				Pinned:      hc.Pinned,
				Title:       hc.Title,
				AppID:       hc.Class,
			})
		}
	}
	
	// Active window to mark focused
	out, err = exec.Command("hyprctl", "-j", "activewindow").Output()
	if err == nil {
		var activeWin struct {
			Address string `json:"address"`
		}
		json.Unmarshal(out, &activeWin)
		for i, w := range state.Windows {
			if w.Address == activeWin.Address {
				state.Windows[i].Focused = true
				break
			}
		}
	}
	
	// Active workspace
	out, err = exec.Command("hyprctl", "-j", "activeworkspace").Output()
	if err == nil {
		var activeWs struct {
			ID int `json:"id"`
		}
		json.Unmarshal(out, &activeWs)
		state.ActiveWorkspaceID = activeWs.ID
	}

	p.currentState = state
	if p.stateUpdateChan != nil {
		p.stateUpdateChan <- state
	}
}
