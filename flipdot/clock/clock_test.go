package clock

import (
	"testing"
)

// MockDisplay for testing
type MockDisplay struct {
	ShowCalls [][28]uint16
}

func (m *MockDisplay) Show(displayData [28]uint16) error {
	m.ShowCalls = append(m.ShowCalls, displayData)
	return nil
}

func TestShowTime(t *testing.T) {
	mock := &MockDisplay{}

	err := showCurrentTime(mock, "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mock.ShowCalls) != 1 {
		t.Fatalf("expected 1 Show call, got %d", len(mock.ShowCalls))
	}

	// Verify that the display data contains time-related patterns
	displayData := mock.ShowCalls[0]

	// Check that first column has data (no border)
	if displayData[0] == 0 {
		t.Error("expected first column to have data (no border)")
	}

	// Check that some columns have data (time display)
	hasData := false
	hasData = false
	for i := 0; i < len(displayData); i++ {
		if displayData[i] != 0 {
			hasData = true
			break
		}
	}
	if !hasData {
		t.Error("expected time display to have some data")
	}
}

func TestShowTimeLarge(t *testing.T) {
	mock := &MockDisplay{}

	err := showCurrentTime(mock, "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mock.ShowCalls) != 1 {
		t.Fatalf("expected 1 Show call, got %d", len(mock.ShowCalls))
	}

	displayData := mock.ShowCalls[0]
	hasData := false
	hasData = false
	for i := 0; i < len(displayData); i++ {
		if displayData[i] != 0 {
			hasData = true
			break
		}
	}
	if !hasData {
		t.Error("expected large time display to have some data")
	}
}
