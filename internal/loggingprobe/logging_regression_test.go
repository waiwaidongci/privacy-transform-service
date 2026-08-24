package loggingprobe

import (
	"sync"
	"testing"
	"time"

	"github.com/ali/go-0821/privacy-transform-service/internal/infrastructure/logging"
)

func TestEventDoesNotMutateCallerFields(t *testing.T) {
	fields := map[string]any{"request_id": "req-1"}
	logging.New().Info("handled", fields)
	if len(fields) != 1 || fields["request_id"] != "req-1" {
		t.Fatalf("caller fields mutated: %#v", fields)
	}
}

func TestConcurrentEventSharedFieldsIsRaceFree(t *testing.T) {
	fields := map[string]any{"request_id": "shared"}
	logger := logging.New()
	start := make(chan struct{})
	var workers sync.WaitGroup
	for i := 0; i < 32; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			logger.Info("parallel", fields)
		}()
	}
	close(start)
	workers.Wait()
	if len(fields) != 1 {
		t.Fatalf("shared fields polluted: %#v", fields)
	}
}

func TestWithCallerReturnsIndependentMap(t *testing.T) {
	input := map[string]any{"request_id": "req-2"}
	called := logging.WithCaller(input)
	if _, ok := input["caller"]; ok {
		t.Fatalf("WithCaller changed input: %#v", input)
	}
	called["request_id"] = "changed"
	if input["request_id"] != "req-2" {
		t.Fatalf("WithCaller returned shared map: %#v", input)
	}
}

func TestWithDurationReturnsIndependentMap(t *testing.T) {
	input := map[string]any{"request_id": "req-3"}
	timed := logging.WithDuration(input, time.Now().Add(-time.Second))
	if _, ok := input["duration_ms"]; ok {
		t.Fatalf("WithDuration changed input: %#v", input)
	}
	timed["request_id"] = "changed-again"
	if input["request_id"] != "req-3" {
		t.Fatalf("WithDuration returned shared map: %#v", input)
	}
}

func TestZeroValueLoggerIsSafe(t *testing.T) {
	var logger logging.Logger
	logger.Info("zero", map[string]any{"ok": true})
	var nilLogger *logging.Logger
	nilLogger.Error("nil", nil)
}
