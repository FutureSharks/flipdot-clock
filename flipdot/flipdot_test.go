package flipdot

import (
	"reflect"
	"testing"
)

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
