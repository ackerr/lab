package utils

import (
	"fmt"
	"os"
)

// Check : return the error message if the err is not nil
func Check(err error) {
	if err != nil {
		Err(err)
	}
}

// Err : return the error message，then exit
func Err(msg ...interface{}) {
	fmt.Println(msg...)
	os.Exit(1)
}
