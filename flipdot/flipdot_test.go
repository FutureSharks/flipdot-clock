package flipdot

import (
	"reflect"
	"testing"
	"time"
)

// MockDisplayOutput for testing
type MockDisplayOutput struct {
	ShowCalls []ShowCall
	ShowError error
	Closed    bool
}

type ShowCall struct {
	DisplayData [28]uint16
	Timestamp   time.Time
}

func (m *MockDisplayOutput) Show(displayData [28]uint16) error {
	m.ShowCalls = append(m.ShowCalls, ShowCall{
		DisplayData: displayData,
		Timestamp:   time.Now(),
	})
	return m.ShowError
}

func (m *MockDisplayOutput) Close() error {
	m.Closed = true
	return nil
}

func TestPrepareSerialFrames(t *testing.T) {
	s := &SerialOutput{}

	testCases := []struct {
		name     string
		input    [28]uint16
		expected [2][28]byte
	}{
		{
			name:     "all zeros",
			input:    [28]uint16{},
			expected: [2][28]byte{},
		},
		{
			name: "all ones",
			input: func() [28]uint16 {
				var data [28]uint16
				for i := range data {
					data[i] = 0b11111111111111 // 14 bits set to 1
				}
				return data
			}(),
			expected: func() [2][28]byte {
				var frames [2][28]byte
				for i := range frames[0] {
					frames[0][i] = 0b1111111 // Lower 7 bits
					frames[1][i] = 0b1111111 // Upper 7 bits
				}
				return frames
			}(),
		},
		{
			name: "single bit top",
			input: func() [28]uint16 {
				var data [28]uint16
				data[0] = 0b1 // bit 0
				return data
			}(),
			expected: func() [2][28]byte {
				var frames [2][28]byte
				frames[0][0] = 0b1
				return frames
			}(),
		},
		{
			name: "single bit bottom",
			input: func() [28]uint16 {
				var data [28]uint16
				data[0] = 0b10000000 // bit 7
				return data
			}(),
			expected: func() [2][28]byte {
				var frames [2][28]byte
				frames[1][0] = 0b1
				return frames
			}(),
		},
		{
			name: "checkerboard",
			input: func() [28]uint16 {
				var data [28]uint16
				for i := range data {
					if i%2 == 0 {
						data[i] = 0b10101010101010
					} else {
						data[i] = 0b01010101010101
					}
				}
				return data
			}(),
			expected: func() [2][28]byte {
				var frames [2][28]byte
				for i := range frames[0] {
					if i%2 == 0 {
						frames[0][i] = 0b0101010
						frames[1][i] = 0b1010101
					} else {
						frames[0][i] = 0b1010101
						frames[1][i] = 0b0101010
					}
				}
				return frames
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := s.prepareSerialFrames(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("unexpected result for %s:\nexpected: %v\nactual:   %v", tc.name, tc.expected, actual)
			}
		})
	}
}

// Test Display creation
func TestNewDisplay(t *testing.T) {
	t.Run("terminal mode", func(t *testing.T) {
		display, err := NewDisplay(true, "", 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if display == nil {
			t.Fatal("expected display to be created")
		}
		defer display.Close()

		// Verify it's using TerminalOutput
		_, ok := display.output.(*TerminalOutput)
		if !ok {
			t.Error("expected TerminalOutput for terminal mode")
		}
	})
}

// Test Display with mock output
func TestDisplayWithMock(t *testing.T) {
	mock := &MockDisplayOutput{}
	display := &Display{output: mock}

	t.Run("Show method", func(t *testing.T) {
		testData := [28]uint16{1, 2, 3, 4, 5}
		err := display.Show(testData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(mock.ShowCalls) != 1 {
			t.Fatalf("expected 1 Show call, got %d", len(mock.ShowCalls))
		}

		if !reflect.DeepEqual(mock.ShowCalls[0].DisplayData, testData) {
			t.Errorf("expected %v, got %v", testData, mock.ShowCalls[0].DisplayData)
		}
	})

	t.Run("Close method", func(t *testing.T) {
		err := display.Close()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !mock.Closed {
			t.Error("expected Close to be called on output")
		}
	})
}

// Test ShowTime method
func TestDisplayShowTime(t *testing.T) {
	mock := &MockDisplayOutput{}
	display := &Display{output: mock}

	err := display.ShowTime()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mock.ShowCalls) != 1 {
		t.Fatalf("expected 1 Show call, got %d", len(mock.ShowCalls))
	}

	// Verify that the display data contains time-related patterns
	displayData := mock.ShowCalls[0].DisplayData

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

// Test RunTestPattern method
func TestDisplayRunTestPattern(t *testing.T) {
	mock := &MockDisplayOutput{}
	display := &Display{output: mock}

	err := display.RunTestPattern()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have 16 frames for the test pattern
	if len(mock.ShowCalls) != 16 {
		t.Fatalf("expected 16 Show calls for test pattern, got %d", len(mock.ShowCalls))
	}

	// Verify that each frame is different (expanding circles)
	for i := 1; i < len(mock.ShowCalls); i++ {
		if reflect.DeepEqual(mock.ShowCalls[i-1].DisplayData, mock.ShowCalls[i].DisplayData) {
			t.Errorf("frame %d and %d should be different", i-1, i)
		}
	}
}

// Test ShowText method
func TestDisplayShowText(t *testing.T) {
	mock := &MockDisplayOutput{}
	display := &Display{output: mock}

	t.Run("simple text no loop", func(t *testing.T) {
		mock.ShowCalls = nil // Reset
		err := display.ShowText("Hi", 1*time.Millisecond, false, "5x8")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should have multiple calls for scrolling
		if len(mock.ShowCalls) == 0 {
			t.Fatal("expected at least one Show call")
		}
	})

	t.Run("invalid font size", func(t *testing.T) {
		err := display.ShowText("Hi", 1*time.Millisecond, false, "invalid")
		if err == nil {
			t.Fatal("expected error for invalid font size")
		}
	})

	t.Run("unsupported character", func(t *testing.T) {
		err := display.ShowText("🚀", 1*time.Millisecond, false, "5x8")
		if err == nil {
			t.Fatal("expected error for unsupported character")
		}
	})
}

// Test prepareText method
func TestDisplayPrepareText(t *testing.T) {
	display := &Display{output: &MockDisplayOutput{}}

	t.Run("valid text 5x8", func(t *testing.T) {
		result, err := display.prepareText("Hi", "5x8")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) == 0 {
			t.Fatal("expected non-empty result")
		}

		// Should have data for 'H', gap, 'i', gap
		// H is 5 columns, i is 5 columns, 2 gaps = 12 total
		expectedMinLength := 12
		if len(result) < expectedMinLength {
			t.Errorf("expected at least %d columns, got %d", expectedMinLength, len(result))
		}
	})

	t.Run("valid text 14x9", func(t *testing.T) {
		result, err := display.prepareText("A", "14x9")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) == 0 {
			t.Fatal("expected non-empty result")
		}

		// A in 14x9 font is 9 columns + 1 gap = 10 total
		expectedLength := 10
		if len(result) != expectedLength {
			t.Errorf("expected %d columns, got %d", expectedLength, len(result))
		}
	})

	t.Run("invalid font size", func(t *testing.T) {
		_, err := display.prepareText("A", "invalid")
		if err == nil {
			t.Fatal("expected error for invalid font size")
		}
	})
}

// Test TerminalOutput
func TestTerminalOutput(t *testing.T) {
	terminal := &TerminalOutput{}

	t.Run("Show method", func(t *testing.T) {
		// Create test pattern
		testData := [28]uint16{}
		testData[0] = 0b11111111111111 // All bits set in first column
		testData[1] = 0b10101010101010 // Alternating pattern in second column

		err := terminal.Show(testData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Note: This test mainly ensures no panic/error occurs
		// Visual output testing would require capturing stdout
	})

	t.Run("Close method", func(t *testing.T) {
		err := terminal.Close()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// Test SerialOutput methods beyond prepareSerialFrames
func TestSerialOutputMethods(t *testing.T) {
	// Note: These tests focus on the logic, not actual serial communication

	t.Run("prepareSerialFrames edge cases", func(t *testing.T) {
		s := &SerialOutput{}

		// Test with maximum values
		input := [28]uint16{}
		for i := range input {
			input[i] = 0b11111111111111 // All 14 bits set
		}

		frames, err := s.prepareSerialFrames(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// All bytes should be 0b1111111 (7 bits set)
		for i := range frames[0] {
			if frames[0][i] != 0b1111111 {
				t.Errorf("expected 0b1111111 in top frame[%d], got %b", i, frames[0][i])
			}
			if frames[1][i] != 0b1111111 {
				t.Errorf("expected 0b1111111 in bottom frame[%d], got %b", i, frames[1][i])
			}
		}
	})
}
