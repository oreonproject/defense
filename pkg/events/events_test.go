// oreon/defense · watchthelight <wtl>

package events

import (
	"context"
	"testing"
	"time"
)

func TestBuilder(t *testing.T) {
	evt := Start(EventTypeScan, "scanner").
		Set("path", "/tmp/test").
		Set("count", 42).
		End()

	if evt.Type != EventTypeScan {
		t.Errorf("Type = %v, want %v", evt.Type, EventTypeScan)
	}
	if evt.Component != "scanner" {
		t.Errorf("Component = %v, want scanner", evt.Component)
	}
	if evt.Fields["path"] != "/tmp/test" {
		t.Errorf("Fields[path] = %v, want /tmp/test", evt.Fields["path"])
	}
	if evt.Fields["count"] != 42 {
		t.Errorf("Fields[count] = %v, want 42", evt.Fields["count"])
	}
	if !evt.Success {
		t.Error("Success should be true when no error set")
	}
	if evt.OperationID == "" {
		t.Error("OperationID should be generated")
	}
	if evt.DurationMs < 0 {
		t.Error("DurationMs should be non-negative")
	}
}

func TestBuilderSetError(t *testing.T) {
	evt := Start(EventTypeScan, "scanner").
		SetError(context.DeadlineExceeded).
		End()

	if evt.Success {
		t.Error("Success should be false when error set")
	}
	if evt.Error != "context deadline exceeded" {
		t.Errorf("Error = %v, want context deadline exceeded", evt.Error)
	}
}

func TestBuilderSetErrorNil(t *testing.T) {
	evt := Start(EventTypeScan, "scanner").
		SetError(nil).
		End()

	if !evt.Success {
		t.Error("Success should be true when nil error set")
	}
}

func TestBuilderWithOperationID(t *testing.T) {
	evt := Start(EventTypeScan, "scanner").
		WithOperationID("custom-id").
		End()

	if evt.OperationID != "custom-id" {
		t.Errorf("OperationID = %v, want custom-id", evt.OperationID)
	}
}

func TestContextOperationID(t *testing.T) {
	ctx := context.Background()

	// No ID initially
	if id := OperationIDFromContext(ctx); id != "" {
		t.Errorf("OperationIDFromContext = %v, want empty", id)
	}

	// Add ID
	ctx = WithOperationID(ctx, "test-op-123")
	if id := OperationIDFromContext(ctx); id != "test-op-123" {
		t.Errorf("OperationIDFromContext = %v, want test-op-123", id)
	}
}

func TestNewOperationContext(t *testing.T) {
	ctx := context.Background()
	ctx, id := NewOperationContext(ctx)

	if id == "" {
		t.Error("NewOperationContext should generate an ID")
	}
	if len(id) != 8 {
		t.Errorf("ID length = %d, want 8", len(id))
	}
	if extracted := OperationIDFromContext(ctx); extracted != id {
		t.Errorf("OperationIDFromContext = %v, want %v", extracted, id)
	}
}

func TestEmitterSampling(t *testing.T) {
	// Test that errors are always emitted
	e := NewEmitter(WithSampleRate(0)) // 0% sample rate

	errorEvt := Event{Success: false, Error: "test error"}
	if !e.shouldEmit(errorEvt) {
		t.Error("Errors should always be emitted")
	}

	// Test that slow operations are always emitted
	slowEvt := Event{Success: true, Duration: 2 * time.Second}
	if !e.shouldEmit(slowEvt) {
		t.Error("Slow operations should always be emitted")
	}

	// Test that fast successes respect sample rate
	fastEvt := Event{Success: true, Duration: 10 * time.Millisecond}
	if e.shouldEmit(fastEvt) {
		t.Error("Fast successes should be sampled at 0% rate")
	}

	// Test 100% sample rate
	e100 := NewEmitter(WithSampleRate(1.0))
	if !e100.shouldEmit(fastEvt) {
		t.Error("Fast successes should be emitted at 100% rate")
	}
}

func TestTypedBuilders(t *testing.T) {
	t.Run("ScanBuilder", func(t *testing.T) {
		evt := StartScan("quick", "job-123").
			FilesScanned(100).
			ThreatsFound(2).
			Path("/home").
			End()

		if evt.Type != EventTypeScan {
			t.Errorf("Type = %v, want %v", evt.Type, EventTypeScan)
		}
		if evt.Fields[FieldScanType] != "quick" {
			t.Errorf("scan_type = %v, want quick", evt.Fields[FieldScanType])
		}
		if evt.Fields[FieldJobID] != "job-123" {
			t.Errorf("job_id = %v, want job-123", evt.Fields[FieldJobID])
		}
		if evt.Fields[FieldFilesScanned] != 100 {
			t.Errorf("files_scanned = %v, want 100", evt.Fields[FieldFilesScanned])
		}
		if evt.Fields[FieldThreatsFound] != 2 {
			t.Errorf("threats_found = %v, want 2", evt.Fields[FieldThreatsFound])
		}
	})

	t.Run("IPCRequestBuilder", func(t *testing.T) {
		evt := StartIPCRequest("status", "req-456").
			ClientVersion(1).
			ResponseSize(256).
			End()

		if evt.Type != EventTypeIPCRequest {
			t.Errorf("Type = %v, want %v", evt.Type, EventTypeIPCRequest)
		}
		if evt.Fields[FieldCommand] != "status" {
			t.Errorf("command = %v, want status", evt.Fields[FieldCommand])
		}
	})

	t.Run("StateChangeBuilder", func(t *testing.T) {
		evt := StartStateChange("protected", "warning").
			Reason("rules outdated").
			End()

		if evt.Type != EventTypeStateChange {
			t.Errorf("Type = %v, want %v", evt.Type, EventTypeStateChange)
		}
		if evt.Fields[FieldFromState] != "protected" {
			t.Errorf("from_state = %v, want protected", evt.Fields[FieldFromState])
		}
		if evt.Fields[FieldToState] != "warning" {
			t.Errorf("to_state = %v, want warning", evt.Fields[FieldToState])
		}
	})

	t.Run("ThreatBuilder", func(t *testing.T) {
		evt := StartThreat("/tmp/virus.exe", "Eicar-Test").
			Action("quarantined").
			FileSize(1024).
			End()

		if evt.Type != EventTypeThreat {
			t.Errorf("Type = %v, want %v", evt.Type, EventTypeThreat)
		}
		if evt.Fields[FieldThreatName] != "Eicar-Test" {
			t.Errorf("threat_name = %v, want Eicar-Test", evt.Fields[FieldThreatName])
		}
	})

	t.Run("HealthCheckBuilder", func(t *testing.T) {
		evt := StartHealthCheck().
			ClamAVAvailable(true).
			FirewallEnabled(false).
			End()

		if evt.Type != EventTypeHealthCheck {
			t.Errorf("Type = %v, want %v", evt.Type, EventTypeHealthCheck)
		}
		if evt.Fields[FieldClamAvailable] != true {
			t.Errorf("clamav_available = %v, want true", evt.Fields[FieldClamAvailable])
		}
	})
}
