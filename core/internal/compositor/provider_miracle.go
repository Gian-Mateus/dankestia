package compositor

import (
	"fmt"
	"time"
)

type MiracleProvider struct {
	updateChan chan<- CompositorState
	stopChan   chan struct{}
}

func NewMiracleProvider() *MiracleProvider {
	return &MiracleProvider{
		stopChan: make(chan struct{}),
	}
}

func (p *MiracleProvider) Start(updateChan chan<- CompositorState) error {
	p.updateChan = updateChan
	go p.mockPoll()
	return nil
}

func (p *MiracleProvider) Stop() {
	close(p.stopChan)
}

func (p *MiracleProvider) GetName() string {
	return "miracle"
}

func (p *MiracleProvider) Dispatch(command string) error {
	return fmt.Errorf("miracle dispatch not fully implemented")
}

func (p *MiracleProvider) mockPoll() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
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
