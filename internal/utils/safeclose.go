package utils

import "io"

// ErrorHandler is a function that receives a non-null error and does something with it
type ErrorHandler func(error)

// SafeClose panics if the closable fails to close
func SafeClose(closable io.Closer) {
	SafeCloseWithStrategy(closable, func(err error) {
		panic(err)
	})
}

// SafeCloseWithStrategy tries to close the closable, and then calls strategy if there is an error
func SafeCloseWithStrategy(closable io.Closer, strategy ErrorHandler) {
	if err := closable.Close(); err != nil {
		strategy(err)
	}
}
