package utils

import (
	"sync"
	"testing"
	"time"
)

func TestGo_Success(t *testing.T) {
	var executed bool
	var wg sync.WaitGroup
	
	wg.Add(1)
	Go(func() {
		executed = true
		wg.Done()
	})
	
	// Wait for goroutine to complete
	wg.Wait()
	
	if !executed {
		t.Error("expected goroutine to execute")
	}
}

func TestGo_PanicRecovery(t *testing.T) {
	var wg sync.WaitGroup
	
	wg.Add(1)
	Go(func() {
		defer wg.Done()
		panic("test panic")
	})
	
	// Wait for goroutine to complete
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	
	// Should complete without crashing the test
	select {
	case <-done:
		// Success - panic was recovered
	case <-time.After(time.Second):
		t.Fatal("goroutine didn't complete - panic may not have been recovered")
	}
}