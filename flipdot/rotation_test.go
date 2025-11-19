package flipdot

import (
	"testing"
)

type MockOutput struct {
	Data [28]uint16
}

func (m *MockOutput) Show(displayData [28]uint16) error {
	m.Data = displayData
	return nil
}

func (m *MockOutput) Close() error {
	return nil
}

func TestDisplay_Show_Flip(t *testing.T) {
	mockOutput := &MockOutput{}
	display := &Display{output: mockOutput, flip: true}

	// Create a pattern: Top-left pixel set
	// In 14x28, top-left is col 0, row 0 (bit 0)
	inputData := [28]uint16{}
	inputData[0] = 1 // 00000000000001

	// Expected result after 180 degree rotation:
	// Bottom-right pixel set
	// Bottom-right is col 27, row 13 (bit 13)
	expectedData := [28]uint16{}
	expectedData[27] = 1 << 13 // 10000000000000

	err := display.Show(inputData)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if mockOutput.Data != expectedData {
		t.Errorf("Expected %v, got %v", expectedData, mockOutput.Data)
	}
}

func TestDisplay_Show_NoFlip(t *testing.T) {
	mockOutput := &MockOutput{}
	display := &Display{output: mockOutput, flip: false}

	inputData := [28]uint16{}
	inputData[0] = 1

	expectedData := [28]uint16{}
	expectedData[0] = 1

	err := display.Show(inputData)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if mockOutput.Data != expectedData {
		t.Errorf("Expected %v, got %v", expectedData, mockOutput.Data)
	}
}
