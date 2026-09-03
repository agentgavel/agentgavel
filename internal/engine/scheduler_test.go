package engine

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSeedScheduler(t *testing.T) {
	seeds := DeterministicSeeds(25)
	sched := &SeedScheduler{Seeds: seeds, MaxWorkers: 4, Timeout: time.Second}
	seen := map[int64]bool{}
	var mu sync.Mutex
	outcomes := sched.Run(context.Background(), func(ctx context.Context, seed int64) error {
		mu.Lock()
		seen[seed] = true
		mu.Unlock()
		return nil
	})
	if len(outcomes) != 25 {
		t.Fatalf("outcomes %d", len(outcomes))
	}
	for _, s := range seeds {
		if !seen[s] {
			t.Fatalf("missing seed %d", s)
		}
	}
}
