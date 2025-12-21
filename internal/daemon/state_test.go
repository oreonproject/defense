// oreon/defense · watchthelight <wtl>

package daemon

import (
	"sync"
	"testing"
)

func TestStateString(t *testing.T) {
	tests := []struct {
		state State
		want  string
	}{
		{StateStarting, "starting"},
		{StateProtected, "protected"},
		{StateWarning, "warning"},
		{StateAlert, "alert"},
		{StateScanning, "scanning"},
		{StatePaused, "paused"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.want {
			t.Errorf("State(%d).String() = %q, want %q", tt.state, got, tt.want)
		}
	}
}

func TestStateManager(t *testing.T) {
	sm := NewStateManager()

	if sm.State() != StateStarting {
		t.Errorf("initial state = %v, want StateStarting", sm.State())
	}

	sm.SetState(StateProtected)
	if sm.State() != StateProtected {
		t.Errorf("state = %v, want StateProtected", sm.State())
	}
}

func TestStateListener(t *testing.T) {
	sm := NewStateManager()

	var called bool
	var gotOld, gotNew State

	sm.OnStateChange(func(old, new State) {
		called = true
		gotOld = old
		gotNew = new
	})

	sm.SetState(StateProtected)

	if !called {
		t.Error("listener was not called")
	}
	if gotOld != StateStarting {
		t.Errorf("old = %v, want StateStarting", gotOld)
	}
	if gotNew != StateProtected {
		t.Errorf("new = %v, want StateProtected", gotNew)
	}
}

func TestStateListenerNotCalledOnSameState(t *testing.T) {
	sm := NewStateManager()
	sm.SetState(StateProtected)

	var called bool
	sm.OnStateChange(func(old, new State) {
		called = true
	})

	sm.SetState(StateProtected) // same state

	if called {
		t.Error("listener should not be called when state doesn't change")
	}
}

func TestStateManagerConcurrent(t *testing.T) {
	sm := NewStateManager()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				sm.SetState(StateProtected)
			} else {
				sm.SetState(StateScanning)
			}
			_ = sm.State()
		}(i)
	}
	wg.Wait()
}
