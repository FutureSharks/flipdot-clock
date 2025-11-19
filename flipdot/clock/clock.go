package clock

import (
	"math/rand"
	"time"

	fonts "github.com/FutureSharks/flipdot-clock/flipdot/fonts"
	log "github.com/sirupsen/logrus"
)

// Display interface for the clock to interact with the hardware
type Display interface {
	Show(displayData [28]uint16) error
}

// Run starts the clock loop
func Run(display Display, mode string) {
	if mode == "transition" {
		runTransitionMode(display)
		return
	}
	if mode == "default" {
		runDefaultMode(display)
		return
	}
	log.Fatalf("Invalid clock-mode value %s. Must be 'default' or 'transition'", mode)
}

func runDefaultMode(display Display) {
	for {
		err := showTime(display)
		if err != nil {
			log.Errorf("Failed to show time: %v", err)
			return
		}
		time.Sleep(1 * time.Minute)
	}
}

func showTime(display Display) error {
	now := time.Now()
	return showTimeFor(display, now)
}

func showTimeFor(display Display, t time.Time) error {
	timeStr := t.Format("15:04")
	displayData, err := renderTime(timeStr)
	if err != nil {
		return err
	}

	log.Debugf("Displaying time: %s", timeStr)
	return display.Show(displayData)
}

func renderTime(timeStr string) ([28]uint16, error) {
	displayData := [28]uint16{}
	result := []uint16{}
	for _, char := range timeStr {
		fontData, err := fonts.GetCharacter(char, "small")
		if err != nil {
			return displayData, err
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
	return displayData, nil
}

func runTransitionMode(display Display) {
	// Initial display
	err := showTime(display)
	if err != nil {
		log.Errorf("Failed to show initial time: %v", err)
	}

	for {
		now := time.Now()
		nextMinute := now.Truncate(time.Minute).Add(time.Minute)

		// Wait for next minute
		timeToWait := nextMinute.Sub(now)
		time.Sleep(timeToWait)

		newTime := time.Now()
		newTimeStr := newTime.Format("15:04")

		prevTime := newTime.Add(-1 * time.Minute)
		prevTimeStr := prevTime.Format("15:04")

		prevData, _ := renderTime(prevTimeStr)
		newData, _ := renderTime(newTimeStr)

		// Calculate diffs
		type pixel struct {
			col int
			row int
			val bool // true if bit is set (1), false if 0
		}

		var diffs []pixel

		for col := range 28 {
			prevCol := prevData[col]
			newCol := newData[col]
			if prevCol == newCol {
				continue
			}
			for row := range 14 {
				prevBit := (prevCol >> row) & 1
				newBit := (newCol >> row) & 1
				if prevBit != newBit {
					diffs = append(diffs, pixel{col: col, row: row, val: newBit == 1})
				}
			}
		}

		// Shuffle diffs
		rand.Shuffle(len(diffs), func(i, j int) {
			diffs[i], diffs[j] = diffs[j], diffs[i]
		})

		// Apply diffs one by one
		currentData := prevData
		for _, p := range diffs {
			// Update currentData
			if p.val {
				currentData[p.col] |= (1 << p.row)
			} else {
				currentData[p.col] &^= (1 << p.row)
			}

			display.Show(currentData)
			time.Sleep(1 * time.Second)
		}

		// Ensure we end up with the exact new data (in case of any drift or logic error)
		display.Show(newData)
	}
}
