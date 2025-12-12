package utils

import (
	"log"
	"runtime/debug"
)

// Go runs a goroutine with panic recovery
func Go(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Goroutine panic recovered: %v\n%s", r, debug.Stack())
			}
		}()
		fn()
	}()
}