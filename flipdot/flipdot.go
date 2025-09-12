package flipdot

import (
	"fmt"
	"math"
	"time"

	fonts "github.com/FutureSharks/flipdot-clock/flipdot/fonts"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

// DisplayOutput interface for different output methods
type DisplayOutput interface {
	Show(displayData [28]uint16) error
	Close() error
}

// Display represents a flipdot display
// it could be a physical Alfa-Zeta 14*28 display connected via serial port or a simulated display that runs in the terminal
type Display struct {
	output DisplayOutput
}

func NewDisplay(terminalMode bool, portName string, baudRate int) (*Display, error) {
	if terminalMode {
		return &Display{output: &TerminalOutput{}}, nil
	}
	return newSerialDisplay(portName, baudRate)
}

// NewSerialDisplay creates a new Display instance with serial output
func newSerialDisplay(portName string, baudRate int) (*Display, error) {
	mode := &serial.Mode{
		BaudRate: baudRate,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open port: %v", err)
	}

	output := &SerialOutput{port: port}
	return &Display{output: output}, nil
}

// Close closes the display connection
func (d *Display) Close() error {
	return d.output.Close()
}

func (d *Display) RunTestPattern() error {
	log.Debug("Running test pattern...")

	centerX := 13.5
	centerY := 6.5

	for i := 0; i < 16; i++ {
		var displayData [28]uint16
		for col := 0; col < 28; col++ {
			var columnData uint16
			for row := 0; row < 14; row++ {
				dx := float64(col) - centerX
				dy := float64(row) - centerY
				distance := math.Sqrt(dx*dx + dy*dy)
				if int(distance) == i {
					columnData |= (1 << row)
				}
			}
			displayData[col] = columnData
		}

		err := d.Show(displayData)
		if err != nil {
			return fmt.Errorf("failed to send test pattern: %v", err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	return nil
}

// ShowTime displays the current time on the 14x28 display.
func (d *Display) ShowTime() error {
	now := time.Now()
	timeStr := now.Format("15:04")
	displayData := [28]uint16{}

	result := []uint16{}
	for _, char := range timeStr {
		fontData, err := fonts.GetCharacter(char, "5x8")
		if err != nil {
			return err
		}
		// add the character
		result = append(result, fontData...)
		// add a small gap before next character
		result = append(result, uint16(0))
	}

	// add left display border
	displayData[0] = 0

	for i, v := range result {
		displayData[i+1] = v
	}

	log.Debugf("Displaying time: %s", timeStr)

	return d.Show(displayData)
}

func (d *Display) Show(displayData [28]uint16) error {
	return d.output.Show(displayData)
}

func (d *Display) ShowText(text string, scrollSpeed time.Duration, loop bool, fontSize string) error {
	for {
		allFrames := []uint16{}

		// start with a blank display
		for range 28 {
			allFrames = append(allFrames, 0)
		}

		t, err := d.prepareText(text, fontSize)
		if err != nil {
			return err
		}
		allFrames = append(allFrames, t...)
		allFrames = append(allFrames, 0)

		for {
			toSend := [28]uint16{}

			for i := range 28 {
				v := uint16(0)
				if i < len(allFrames) {
					v = allFrames[i]
				}
				toSend[i] = v
			}

			err := d.Show(toSend)
			if err != nil {
				return err
			}

			// delete first column to scroll the display right to left
			allFrames = allFrames[1:]

			if len(allFrames) == 0 {
				break
			}

			time.Sleep(scrollSpeed)
		}

		if !loop {
			break
		}
	}

	return nil
}

func (d *Display) prepareText(text string, fontSize string) ([]uint16, error) {
	result := []uint16{}
	for _, char := range text {
		letterData, err := fonts.GetCharacter(char, fontSize)
		if err != nil {
			return nil, err
		}
		result = append(result, letterData...)
		result = append(result, uint16(0))
	}

	return result, nil
}
