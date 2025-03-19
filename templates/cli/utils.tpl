package internal

import (
	"fmt"
	"strings"
)

// Constants for terminal colors
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
)

// PrintInfo prints an info message to the console
func PrintInfo(format string, a ...interface{}) {
	fmt.Printf(ColorBlue+"INFO: "+format+ColorReset+"\n", a...)
}

// PrintSuccess prints a success message to the console
func PrintSuccess(format string, a ...interface{}) {
	fmt.Printf(ColorGreen+"SUCCESS: "+format+ColorReset+"\n", a...)
}

// PrintWarning prints a warning message to the console
func PrintWarning(format string, a ...interface{}) {
	fmt.Printf(ColorYellow+"WARNING: "+format+ColorReset+"\n", a...)
}

// PrintError prints an error message to the console
func PrintError(format string, a ...interface{}) {
	fmt.Printf(ColorRed+"ERROR: "+format+ColorReset+"\n", a...)
}

// StringInSlice checks if a string is in a slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.EqualFold(b, a) {
			return true
		}
	}
	return false
} 