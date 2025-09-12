package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/FutureSharks/flipdot-clock/flipdot"

	log "github.com/sirupsen/logrus"
)

func main() {
	portName := flag.String("serial-port", "/dev/ttyS0", "The serial port connected to the displays")
	baudRate := flag.Int("serial-baud", 57600, "The baud rate for the serial connection.")
	terminalMode := flag.Bool("terminal", false, "Display output to terminal instead of serial port.")
	testPattern := flag.Bool("test-pattern", false, "Display a test pattern and then exit")
	clock := flag.Bool("clock", false, "Run the clock")
	text := flag.String("text", "", "Display some text")
	textLoop := flag.Bool("text-loop", false, "Loop text continuously")
	textSize := flag.String("text-size", "14x9", "Size of each character. Value must be one of 14x9 or 5x8")
	scrollSpeed := flag.Int("text-scroll-speed", 5, "Text scroll speed. 1 is slow, 9 is fast")
	debugLogging := flag.Bool("debug", false, "Enable debug logging")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "flipdot-clock: a small tool for displaying text or the time on a Alfa-Zeta XY5 14*28 flipdot display\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *textSize != "14x9" && *textSize != "5x8" {
		log.Fatalf("Invalid text-size value %s. Must be 14x9 or 5x8", *textSize)
	}

	if *scrollSpeed < 1 || *scrollSpeed > 9 {
		log.Fatalf("Invalid scroll-speed value %d. Must be between 1 and 9.", *scrollSpeed)
	}

	if *debugLogging {
		log.SetLevel(log.DebugLevel)
	}

	// Create a new display instance
	display, err := flipdot.NewDisplay(*terminalMode, *portName, *baudRate)

	if err != nil {
		log.Fatalf("Failed to create display: %v", err)
	}
	defer display.Close()

	if *testPattern {
		err = display.RunTestPattern()
		if err != nil {
			log.Fatalf("Failed to run test pattern: %v", err)
		}
	} else if *text != "" {
		sleepDuration := time.Duration(190-(*scrollSpeed*20)) * time.Millisecond
		err = display.ShowText(*text, sleepDuration, *textLoop, *textSize)
		if err != nil {
			log.Fatalf("Failed to show text: %v", err)
		}
	} else if *clock {
		for {
			err = display.ShowTime()
			if err != nil {
				log.Fatalf("Failed to show time: %v", err)
			}
			time.Sleep(1 * time.Minute)
		}
	} else {
		log.Infoln("No mode selected. Exiting.")
	}
}
