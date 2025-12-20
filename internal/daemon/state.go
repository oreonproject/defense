package daemon

import (
	"sync"
)

// State represents the current protection status of the daemon.
// The tray icon and UI will reflect this state.
type State int

const (
	StateStarting  State = iota // daemon is initializing
	StateProtected              // everything is good
	StateWarning                // something needs attention (e.g. rules outdated)
	StateAlert                  // something is wrong (e.g. threat detected)
	StateScanning               // scan in progress
	StatePaused                 // protection temporarily disabled
)

func (s State) String() string {
	switch s {
	case StateStarting:
		return "starting"
	case StateProtected:
		return "protected"
	case StateWarning:
		return "warning"
	case StateAlert:
		return "alert"
	case StateScanning:
		return "scanning"
	case StatePaused:
		return "paused"
	default:
		return "unknown"
	}
}

// StateListener is called whenever the state changes.
// Implement this to react to state changes (e.g. update tray icon).
type StateListener func(old, new State)

// StateManager handles state transitions and notifies listeners.
// Thread-safe - can be called from multiple goroutines.
type StateManager struct {
	mu        sync.RWMutex
	state     State
	listeners []StateListener
}

func NewStateManager() *StateManager {
	return &StateManager{
		state: StateStarting,
	}
}

// State returns the current state.
func (sm *StateManager) State() State {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.state
}

// SetState changes the state and notifies all listeners.
// Listeners are called synchronously - keep them fast.
func (sm *StateManager) SetState(s State) {
	sm.mu.Lock()
	old := sm.state
	sm.state = s
	listeners := sm.listeners // copy slice header for safe iteration
	sm.mu.Unlock()

	if old != s {
		for _, fn := range listeners {
			fn(old, s)
		}
	}
}

// OnStateChange registers a listener that's called when state changes.
// Returns a function to unregister the listener.
//
// Example (for tray icon):
//
//	sm.OnStateChange(func(old, new State) {
//	    updateTrayIcon(new)
//	})
//
// Example (for firewall integration):
//
//	sm.OnStateChange(func(old, new State) {
//	    if new == StatePaused {
//	        // maybe log that protection is paused
//	    }
//	})
func (sm *StateManager) OnStateChange(fn StateListener) func() {
	sm.mu.Lock()
	sm.listeners = append(sm.listeners, fn)
	idx := len(sm.listeners) - 1
	sm.mu.Unlock()

	return func() {
		sm.mu.Lock()
		defer sm.mu.Unlock()
		// remove listener by setting to nil (avoids slice realloc)
		if idx < len(sm.listeners) {
			sm.listeners[idx] = nil
		}
	}
}
