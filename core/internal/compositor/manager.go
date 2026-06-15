package compositor

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/AvengeMedia/Dankestia/core/internal/server/models"
)

type Manager struct {
	provider        Provider
	stateUpdateChan chan CompositorState
	currentState    *CompositorState
	mu              sync.RWMutex
	subscribers     sync.Map
}

func NewManager() (*Manager, error) {
	m := &Manager{
		stateUpdateChan: make(chan CompositorState, 10),
	}

	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") != "" {
		m.provider = NewHyprlandProvider()
	} else if os.Getenv("NIRI_SOCKET") != "" || os.Getenv("WAYLAND_DISPLAY") != "" {
		m.provider = NewNiriProvider()
	} else {
		return nil, fmt.Errorf("no supported compositor detected")
	}

	if err := m.provider.Start(m.stateUpdateChan); err != nil {
		return nil, err
	}

	go m.listenUpdates()

	return m, nil
}

func (m *Manager) Stop() {
	if m.provider != nil {
		m.provider.Stop()
	}
	close(m.stateUpdateChan)
}

func (m *Manager) listenUpdates() {
	for state := range m.stateUpdateChan {
		m.mu.Lock()
		m.currentState = &state
		m.mu.Unlock()
		
		m.subscribers.Range(func(key, value interface{}) bool {
			if ch, ok := value.(chan CompositorState); ok {
				select {
				case ch <- state:
				default:
				}
			}
			return true
		})
	}
}

func (m *Manager) Subscribe(clientID string) chan CompositorState {
	ch := make(chan CompositorState, 64)
	m.subscribers.Store(clientID, ch)
	return ch
}

func (m *Manager) Unsubscribe(clientID string) {
	if ch, ok := m.subscribers.LoadAndDelete(clientID); ok {
		close(ch.(chan CompositorState))
	}
}

func (m *Manager) GetState() CompositorState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.currentState == nil {
		return CompositorState{}
	}
	return *m.currentState
}

func (m *Manager) Dispatch(command string) error {
	if m.provider != nil {
		return m.provider.Dispatch(command)
	}
	return fmt.Errorf("provider not initialized")
}

func HandleRequest(conn net.Conn, req models.Request, manager *Manager) {
	switch req.Method {
	case "compositor.getState":
		state := manager.GetState()
		models.Respond(conn, req.ID, state)
		
	case "compositor.dispatch":
		if cmd, ok := models.Get[string](req, "command"); ok {
			err := manager.Dispatch(cmd)
			if err != nil {
				models.RespondError(conn, req.ID, err.Error())
			} else {
				models.Respond(conn, req.ID, models.SuccessResult{Success: true})
			}
		} else {
			models.RespondError(conn, req.ID, "missing 'command' parameter")
		}
		
	default:
		models.RespondError(conn, req.ID, fmt.Sprintf("unknown method: %s", req.Method))
	}
}
