package flipdot

import (
	"fmt"
	"strings"
)

// TerminalOutput implements DisplayOutput for terminal display
type TerminalOutput struct{}

// Show method for TerminalOutput - matches DisplayOutput interface
func (t *TerminalOutput) Show(displayData [28]uint16) error {
	// Clear screen and move cursor to top
	fmt.Print("\033[2J\033[H")

	// Display the flipdot pattern as ASCII art
	fmt.Println("Flipdot Display Output:")
	fmt.Println("┌" + strings.Repeat("─", 28*3) + "┐")

	// Display all 14 rows (top 7 rows from lower bits, bottom 7 rows from upper bits)
	for row := 0; row < 14; row++ {
		fmt.Print("│")
		for col := 0; col < 28; col++ {
			if displayData[col]&(1<<row) != 0 {
				fmt.Print(" ● ")
			} else {
				fmt.Print("   ")
			}
		}
		fmt.Println("│")
	}

	fmt.Println("└" + strings.Repeat("─", 28*3) + "┘")
	return nil
}

// Close method for TerminalOutput
func (t *TerminalOutput) Close() error {
	return nil // Nothing to close for terminal output
}
