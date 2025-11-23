package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/FutureSharks/flipdot-clock/flipdot"
	"github.com/FutureSharks/flipdot-clock/flipdot/clock"

	log "github.com/sirupsen/logrus"
)

func main() {
	portName := flag.String("serial-port", "/dev/ttyS0", "The serial port connected to the displays")
	baudRate := flag.Int("serial-baud", 57600, "The baud rate for the serial connection.")
	terminalMode := flag.Bool("terminal", false, "Display output to terminal instead of serial port.")
	testPattern := flag.Bool("test-pattern", false, "Display a test pattern and then exit")
	runClock := flag.Bool("clock", false, "Run the clock")
	clockMode := flag.String("clock-mode", "default", "The mode to run the clock in. Must be one of 'default' or 'transition'")
	clockSize := flag.String("clock-size", "1", "Size of the clock font. Must be one of '1' (small) or '2' (medium)")
	text := flag.String("text", "", "Display some text")
	textLoop := flag.Bool("text-loop", false, "Loop text continuously")
	textSize := flag.String("text-size", "3", "Size of each character. Value must be one of '1' (small), '2' (medium) or '3' (large)")
	scrollSpeed := flag.Int("text-scroll-speed", 5, "Text scroll speed. 1 is slow, 9 is fast")
	debugLogging := flag.Bool("debug", false, "Enable debug logging")
	flipDisplay := flag.Bool("flip-display", false, "Rotate the display 180 degrees")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "flipdot-clock: a small tool for displaying text or the time on a Alfa-Zeta XY5 14*28 flipdot display\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *textSize != "1" && *textSize != "2" && *textSize != "3" {
		log.Fatalf("Invalid text-size value %s. Must be '1', '2' or '3'", *textSize)
	}

	if *clockSize != "1" && *clockSize != "2" {
		log.Fatalf("Invalid clock-size value %s. Must be '1' or '2'", *clockSize)
	}

	if *scrollSpeed < 1 || *scrollSpeed > 9 {
		log.Fatalf("Invalid scroll-speed value %d. Must be between 1 and 9.", *scrollSpeed)
	}

	if *debugLogging {
		log.SetLevel(log.DebugLevel)
	}

	if *runClock && *text != "" {
		log.Fatalf("Cannot specify both clock and text")
	}

	// Create a new display instance
	display, err := flipdot.NewDisplay(*terminalMode, *portName, *baudRate, *flipDisplay)

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
	} else if *runClock {
		clock.Run(display, *clockMode, *clockSize)
	} else {
		log.Infoln("No mode selected. Use '-clock-mode' or '-text' arguments. Exiting.")
	}
}
