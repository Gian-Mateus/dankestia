package compositor

import (
	"fmt"
	"time"
)

type SwayProvider struct {
	updateChan chan<- CompositorState
	stopChan   chan struct{}
}

func NewSwayProvider() *SwayProvider {
	return &SwayProvider{
		stopChan: make(chan struct{}),
	}
}

func (p *SwayProvider) Start(updateChan chan<- CompositorState) error {
	p.updateChan = updateChan
	go p.mockPoll()
	return nil
}

func (p *SwayProvider) Stop() {
	close(p.stopChan)
}

func (p *SwayProvider) GetName() string {
	return "sway"
}

func (p *SwayProvider) Dispatch(command string) error {
	// Future: run `swaymsg`
	return fmt.Errorf("sway dispatch not fully implemented")
}

func (p *SwayProvider) mockPoll() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Stub: Would query `swaymsg -t get_workspaces`
			state := CompositorState{
				Workspaces: []Workspace{{ID: 1, Name: "1"}},
				Monitors:   []Monitor{{ID: 1, Name: "eDP-1", Focused: true}},
				Windows:    []Window{},
				ActiveWorkspaceID: 1,
			}
			select {
			case p.updateChan <- state:
			case <-p.stopChan:
				return
			}
		case <-p.stopChan:
			return
		}
	}
}
