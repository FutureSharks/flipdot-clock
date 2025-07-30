package flipdot

import (
	"go.bug.st/serial"
)

var (
	// The two 7x28 displays must have DIP switch #1 positions 0-5 set to match the values here
	addresses = []byte{0x01, 0x02}
)

// SerialOutput implements DisplayOutput for serial communication
type SerialOutput struct {
	port serial.Port
}

// Write method for SerialOutput
// For writing the whole 14*28 display
func (s *SerialOutput) write(frames [2][28]byte) error {
	for i := range 2 {
		err := s.writeSingleDisplay(frames[i], addresses[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// Close method for SerialOutput
func (s *SerialOutput) Close() error {
	return s.port.Close()
}

func (s *SerialOutput) prepareSerialFrames(data [28]uint16) ([2][28]byte, error) {
	var frames [2][28]byte
	var topDisplayFrame [28]byte
	var bottomDisplayFrame [28]byte

	for col := range 28 {
		bottomByte := data[col] >> 7
		bottomByteShiftedAgain := bottomByte << 7
		topByte := data[col] &^ bottomByteShiftedAgain
		topDisplayFrame[col] = byte(topByte)
		bottomDisplayFrame[col] = byte(bottomByte)
	}

	frames[0] = topDisplayFrame
	frames[1] = bottomDisplayFrame

	return frames, nil
}

// writeSingleDisplay method for SerialOutput
// writes all pixels of one of the two connected 7*28 displays
func (s *SerialOutput) writeSingleDisplay(data [28]byte, address byte) error {
	var frame []byte

	// from the protocol document: all frames start with this value
	frame = append(frame, 0x80)
	// from the protocol document: enables REFRESH mode, the display will update pixels as the data is received
	frame = append(frame, 0x83)
	// from the protocol document: needs to match the value set on the display controller board DIP switch #1 positions 0-5
	frame = append(frame, address)
	// the 28 columns of the display
	frame = append(frame, data[:]...)
	// from the protocol document: all frames end with this value
	frame = append(frame, 0x8F)

	_, err := s.port.Write(frame)
	if err != nil {
		return err
	}

	return nil
}

func (s *SerialOutput) Show(displayData [28]uint16) error {
	frames, err := s.prepareSerialFrames(displayData)
	if err != nil {
		return err
	}

	return s.write(frames)
}
