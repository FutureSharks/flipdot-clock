package clock

import (
	"math/rand"
	"strings"
	"time"

	fonts "github.com/FutureSharks/flipdot-clock/flipdot/fonts"
	log "github.com/sirupsen/logrus"
)

type Display interface {
	Show(displayData [28]uint16) error
}

func Run(display Display, mode string, size string) {
	if mode == "transition" {
		runTransitionMode(display, size)
		return
	}
	if mode == "default" {
		runDefaultMode(display, size)
		return
	}
	log.Fatalf("Invalid clock-mode value %s. Must be 'default' or 'transition'", mode)
}

func runDefaultMode(display Display, size string) {
	for {
		err := showCurrentTime(display, size)
		if err != nil {
			log.Errorf("Failed to show time: %v", err)
			return
		}
		time.Sleep(1 * time.Minute)
	}
}

func showCurrentTime(display Display, size string) error {
	displayData, err := renderTime(time.Now(), size)
	if err != nil {
		return err
	}

	return display.Show(displayData)
}

func renderTime(t time.Time, size string) ([28]uint16, error) {
	timeStr := t.Format("15:04")
	displayData := [28]uint16{}
	result := []uint16{}

	// some size specific fixes
	// remove the colon so the larger size can fit on the display
	if size == "2" {
		timeStr = strings.Replace(timeStr, ":", "", 1)
	}
	// add left display border
	if size == "1" {
		displayData[0] = 0
	}

	for i, char := range timeStr {
		fontData, err := fonts.GetCharacter(char, size)
		if err != nil {
			return displayData, err
		}

		if size == "2" && i == 2 {
			// add an extra 1 column gap between the hours and minutes
			result = append(result, uint16(0))
		}

		// add the character to the results
		result = append(result, fontData...)
		// add a 1 column gap before next character
		result = append(result, uint16(0))
	}

	for i, v := range result {
		if i < len(displayData) {
			if size == "1" && i == 0 {
				continue
			}
			displayData[i] = v
		}
	}

	return displayData, nil
}

func runTransitionMode(display Display, size string) {
	err := showCurrentTime(display, size)
	if err != nil {
		log.Errorf("Failed to show initial time: %v", err)
	}

	for {
		now := time.Now()
		nextMinute := now.Truncate(time.Minute).Add(time.Minute)

		timeToWait := nextMinute.Sub(now)
		time.Sleep(timeToWait)

		newTime := time.Now()
		prevTime := newTime.Add(-1 * time.Minute)

		prevData, _ := renderTime(prevTime, size)
		newData, _ := renderTime(newTime, size)

		type pixel struct {
			col int
			row int
			val bool
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

		log.Infof("Found %d diffs", len(diffs))

		// Shuffle diffs create a curious transition effect
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
			time.Sleep(calculateTransitionPixelInterval(len(diffs)))
		}

		// Ensure we end up with the exact new data (in case of any drift or logic error)
		display.Show(newData)
	}
}

func calculateTransitionPixelInterval(diffCount int) time.Duration {
	// get it done in 8 seconds
	transitionTimeLimit := 8 * time.Second

	if diffCount == 0 {
		return 0
	}

	return transitionTimeLimit / time.Duration(diffCount)
}
