package utils

import (
	"fmt"
	"os"
)

// Err will return the error message
func Err(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
