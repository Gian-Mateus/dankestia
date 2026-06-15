package compositor

import (
	"fmt"
	"time"
)

type MangoWCProvider struct {
	updateChan chan<- CompositorState
	stopChan   chan struct{}
}

func NewMangoWCProvider() *MangoWCProvider {
	return &MangoWCProvider{
		stopChan: make(chan struct{}),
	}
}

func (p *MangoWCProvider) Start(updateChan chan<- CompositorState) error {
	p.updateChan = updateChan
	go p.mockPoll()
	return nil
}

func (p *MangoWCProvider) Stop() {
	close(p.stopChan)
}

func (p *MangoWCProvider) GetName() string {
	return "mangowc"
}

func (p *MangoWCProvider) Dispatch(command string) error {
	return fmt.Errorf("mangowc dispatch not fully implemented")
}

func (p *MangoWCProvider) mockPoll() {
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
