package clock

import (
	"time"

	fonts "github.com/FutureSharks/flipdot-clock/flipdot/fonts"
	log "github.com/sirupsen/logrus"
)

// Display interface for the clock to interact with the hardware
type Display interface {
	Show(displayData [28]uint16) error
}

// Run starts the clock loop
func Run(display Display) {
	for {
		err := showTime(display)
		if err != nil {
			log.Errorf("Failed to show time: %v", err)
		}
		time.Sleep(1 * time.Minute)
	}
}

func showTime(display Display) error {
	now := time.Now()
	timeStr := now.Format("15:04")
	displayData := [28]uint16{}

	result := []uint16{}
	for _, char := range timeStr {
		fontData, err := fonts.GetCharacter(char, "small")
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
		if i+1 < len(displayData) {
			displayData[i+1] = v
		}
	}

	log.Debugf("Displaying time: %s", timeStr)

	return display.Show(displayData)
}
