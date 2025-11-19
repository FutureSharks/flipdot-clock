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

	err := showTime(mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mock.ShowCalls) != 1 {
		t.Fatalf("expected 1 Show call, got %d", len(mock.ShowCalls))
	}

	// Verify that the display data contains time-related patterns
	displayData := mock.ShowCalls[0]

	// Check that first column is empty (border)
	if displayData[0] != 0 {
		t.Error("expected first column to be empty for border")
	}

	// Check that some columns have data (time display)
	hasData := false
	for i := 1; i < len(displayData); i++ {
		if displayData[i] != 0 {
			hasData = true
			break
		}
	}
	if !hasData {
		t.Error("expected time display to have some data")
	}
}
