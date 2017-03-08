package resource

import (
	"fmt"
	"os"

	"github.com/mitchellh/colorstring"
)

// Fatal -
func Fatal(message string) {
	fmt.Fprintf(os.Stderr, colorstring.Color(message))
	os.Exit(1)
}

// Fatalf -
func Fatalf(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, colorstring.Color(message), args...)
	os.Exit(1)
}
