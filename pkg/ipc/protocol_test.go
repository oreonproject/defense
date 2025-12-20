package ipc

import (
	"encoding/json"
	"testing"
)

func TestRequestJSON(t *testing.T) {
	req := Request{
		ID:      "123",
		Command: CmdScan,
		Params:  json.RawMessage(`{"type":"quick"}`),
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got Request
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.ID != req.ID || got.Command != req.Command {
		t.Errorf("round-trip failed: got %+v, want %+v", got, req)
	}
}

func TestResponseJSON(t *testing.T) {
	resp := Response{
		ID:      "123",
		Success: true,
		Data: StatusResponse{
			State:           "protected",
			FirewallEnabled: true,
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	// check it contains expected fields
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}

	if m["success"] != true {
		t.Errorf("success = %v, want true", m["success"])
	}
	if m["id"] != "123" {
		t.Errorf("id = %v, want 123", m["id"])
	}
}

func TestScanParamsJSON(t *testing.T) {
	params := ScanParams{Type: "quick"}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	expected := `{"type":"quick"}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", data, expected)
	}
}
