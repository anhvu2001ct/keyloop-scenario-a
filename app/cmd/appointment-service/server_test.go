package main

import (
	"context"
	"scenario-a/internal/config"
	"scenario-a/internal/dep"
	"sync"
	"testing"
	"time"
)

func TestStartServer(t *testing.T) {
	config.MustInitForTest()
	deps := dep.Init(nil)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	done := make(chan struct{})

	go func() {
		wg.Done()
		startServer(ctx, 0, deps)
		close(done)
	}()

	wg.Wait()
	time.Sleep(200 * time.Millisecond) // wait a little for server to start

	cancel() // stop server

	select {
	case <-done:
		// no panic, success!
	case <-time.After(1000 * time.Millisecond):
		t.Fatal("app take too long to shutdown")
	}
}
