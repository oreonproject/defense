// oreon/defense · watchthelight <wtl>

package daemon

import (
	"sync"
	"testing"
	"time"
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

	done := make(chan struct{})
	var gotOld, gotNew State

	sm.OnStateChange(func(old, new State) {
		gotOld = old
		gotNew = new
		close(done)
	})

	sm.SetState(StateProtected)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("listener was not called within timeout")
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
	time.Sleep(10 * time.Millisecond) // wait for first state change to process

	called := make(chan struct{}, 1)
	sm.OnStateChange(func(old, new State) {
		called <- struct{}{}
	})

	sm.SetState(StateProtected) // same state

	select {
	case <-called:
		t.Error("listener should not be called when state doesn't change")
	case <-time.After(50 * time.Millisecond):
		// expected - listener was not called
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
