package messages

import (
	"fmt"
)

// Call this when you want to die
func CheckErrFatal(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

// This is useless.
// I'm keeping it here to remind myself how to do it.
func CheckErrWrap(e error, wrap_string string) error {
	if e != nil {
		return fmt.Errorf(wrap_string+" :%w", e)
	}
	return nil
}
