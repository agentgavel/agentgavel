package engine

import (
	"context"
	"sync"
	"time"
)

// SeedOutcome is the result of one seed run.
type SeedOutcome struct {
	Seed int64
	Err  error
}

// SeedFunc runs one seed.
type SeedFunc func(ctx context.Context, seed int64) error

// SeedScheduler runs seeds with a worker pool and per-seed timeout.
type SeedScheduler struct {
	Seeds      []int64
	MaxWorkers int
	Timeout    time.Duration
}

// Run executes fn for every seed and returns outcomes in completion order.
func (s *SeedScheduler) Run(ctx context.Context, fn SeedFunc) []SeedOutcome {
	workers := s.MaxWorkers
	if workers <= 0 {
		workers = 1
	}
	jobs := make(chan int64)
	var (
		mu       sync.Mutex
		outcomes []SeedOutcome
		wg       sync.WaitGroup
	)
	worker := func() {
		defer wg.Done()
		for seed := range jobs {
			seedCtx := ctx
			cancel := func() {}
			if s.Timeout > 0 {
				seedCtx, cancel = context.WithTimeout(ctx, s.Timeout)
			}
			err := fn(seedCtx, seed)
			cancel()
			mu.Lock()
			outcomes = append(outcomes, SeedOutcome{Seed: seed, Err: err})
			mu.Unlock()
		}
	}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}
	for _, seed := range s.Seeds {
		select {
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return outcomes
		case jobs <- seed:
		}
	}
	close(jobs)
	wg.Wait()
	return outcomes
}

// DeterministicSeeds returns [0..n).
func DeterministicSeeds(n int) []int64 {
	out := make([]int64, n)
	for i := 0; i < n; i++ {
		out[i] = int64(i)
	}
	return out
}
