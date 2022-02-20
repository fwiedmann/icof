package gocrazy_permanent_data

import (
	"context"
	"encoding/json"
	"github.com/fwiedmann/icof"
	"io/ioutil"
	"sync"
)

// StateLocation where the state will be stored
var StateLocation = "/perm/icof/state.json"

// NewStateRepository init a new State repository
func NewStateRepository() *State {
	return &State{}
}

// State implements the icof.StateRepository
type State struct {
	mtx sync.RWMutex
}

// SavedState content for state file
type SavedState struct {
	State icof.ObserverState `json:"state"`
}

// Save stores state in state file under StateLocation
func (s *State) Save(ctx context.Context, state icof.ObserverState) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	content, err := json.Marshal(&SavedState{State: state})
	if err != nil {
		return err
	}

	return ioutil.WriteFile(StateLocation, content, 0700)
}

// GetLatest gets state from state file under StateLocation
func (s *State) GetLatest(ctx context.Context) (icof.ObserverState, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	content, err := ioutil.ReadFile(StateLocation)
	if err != nil {
		return false, err
	}

	var state SavedState
	if err := json.Unmarshal(content, &s); err != nil {
		return false, err
	}
	return state.State, nil
}
