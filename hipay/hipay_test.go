package hipay

import (
	"testing"
	"time"

	"resty.dev/v3"
)

func TestNew(t *testing.T) {
	h := New("https://api.hipay.mn", "token", "entity_id")
	if h == nil {
		t.Fatal("Expected New() to return a Hipay instance, got nil")
	}

	_, ok := h.(*hipay)
	if !ok {
		t.Fatal("Expected New() to return a *hipay instance")
	}
}

// TestNew_nilOptionNoPanic ensures passing nil in options does not panic.
func TestNew_nilOptionNoPanic(t *testing.T) {
	h := New("https://api.hipay.mn", "token", "entity_id", nil, nil)
	if h == nil {
		t.Fatal("Expected New() to return a Hipay instance")
	}
}

// TestNew_WithClient applies a custom client and ensures no panic when client is nil.
func TestNew_WithClient(t *testing.T) {
	// WithClient(nil) must not panic; client is ignored when nil.
	h := New("https://api.hipay.mn", "token", "entity_id", WithClient(nil))
	if h == nil {
		t.Fatal("Expected New() to return a Hipay instance")
	}

	// WithClient(non-nil) should use the provided client.
	custom := resty.New().SetTimeout(5 * time.Second)
	h2 := New("https://api.hipay.mn", "token", "entity_id", WithClient(custom))
	if h2 == nil {
		t.Fatal("Expected New() to return a Hipay instance")
	}
	impl, ok := h2.(*hipay)
	if !ok {
		t.Fatal("Expected *hipay")
	}
	if impl.client != custom {
		t.Error("Expected WithClient(custom) to set client to custom")
	}
}

// TestFormatError_nilDetailNoPanic ensures formatError does not panic when details contains nil.
func TestFormatError_nilDetailNoPanic(t *testing.T) {
	h := &hipay{}
	// API could return "details": [null]; slice has len 1 but details[0] is nil.
	err := h.formatError("desc", "msg", []*Details{nil})
	if err == nil {
		t.Fatal("Expected formatError to return non-nil error")
	}
	if err.Error() != "msg" {
		t.Errorf("Expected error message 'msg', got %q", err.Error())
	}
}

func TestFormatError_withValidDetail(t *testing.T) {
	h := &hipay{}
	d := &Details{Field: "amount", Issue: "invalid"}
	err := h.formatError("desc", "msg", []*Details{d})
	if err == nil {
		t.Fatal("Expected non-nil error")
	}
	expect := "msg: amount - invalid"
	if err.Error() != expect {
		t.Errorf("Expected %q, got %q", expect, err.Error())
	}
}
