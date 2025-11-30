package jobs

import (
	"container/heap"
	"context"
	"sync"
	"time"

	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
)

// scheduledJob holds a job payload and its scheduled execution time (unix seconds).
type scheduledJob struct {
	executeAt int64
	payload   interface{}
	index     int // heap index
}

type jobMinHeap []*scheduledJob

func (h jobMinHeap) Len() int            { return len(h) }
func (h jobMinHeap) Less(i, j int) bool  { return h[i].executeAt < h[j].executeAt }
func (h jobMinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i]; h[i].index = i; h[j].index = j }
func (h *jobMinHeap) Push(x interface{}) { *h = append(*h, x.(*scheduledJob)) }
func (h *jobMinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// InMemoryScheduler is an in-process scheduler that executes jobs via a provided handler.
// Start it with Run in a goroutine; stop with Stop.
type InMemoryScheduler struct {
	mu          sync.Mutex
	jobs        jobMinHeap
	closed      bool
	subscribers map[int]func(interface{})
	nextSubID   int
	wake        chan struct{}
}

var _ ports.Scheduler = (*InMemoryScheduler)(nil)

// NewInMemoryScheduler creates a scheduler with no-op subscribers initially.
func NewInMemoryScheduler() *InMemoryScheduler {
	s := &InMemoryScheduler{wake: make(chan struct{}, 1), subscribers: make(map[int]func(interface{}))}
	heap.Init(&s.jobs)
	return s
}

// Schedule implements ports.Scheduler.
func (s *InMemoryScheduler) Schedule(job interface{}, executeAt int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return context.Canceled
	}
	heap.Push(&s.jobs, &scheduledJob{executeAt: executeAt, payload: job})
	// best-effort wake
	select {
	case s.wake <- struct{}{}:
	default:
	}
	return nil
}

// Run processes jobs until Stop is called. Safe to run in a single goroutine.
func (s *InMemoryScheduler) Run(ctx context.Context) {
	for {
		s.mu.Lock()
		if s.closed {
			s.mu.Unlock()
			return
		}
		if s.jobs.Len() == 0 {
			s.mu.Unlock()
			select {
			case <-s.wake:
				// new job scheduled
				continue
			case <-ctx.Done():
				return
			}
		}
		next := s.jobs[0]
		now := time.Now().Unix()
		waitSec := next.executeAt - now
		if waitSec > 0 {
			s.mu.Unlock()
			select {
			case <-time.After(time.Duration(waitSec) * time.Second):
				// time reached; loop to execute
			case <-s.wake:
				// new job or stop; loop to recalc
			case <-ctx.Done():
				return
			}
			continue
		}
		// pop and execute
		heap.Pop(&s.jobs)
		payload := next.payload
		// copy subscribers to avoid holding lock during callbacks
		subs := make([]func(interface{}), 0, len(s.subscribers))
		for _, cb := range s.subscribers {
			subs = append(subs, cb)
		}
		s.mu.Unlock()
		for _, cb := range subs {
			cb(payload)
		}
	}
}

// Stop stops the scheduler loop and unblocks Run.
func (s *InMemoryScheduler) Stop() {
	s.mu.Lock()
	s.closed = true
	s.mu.Unlock()
	// best-effort wake to unblock Run selects
	select {
	case s.wake <- struct{}{}:
	default:
	}
}

// Subscribe registers a callback to receive scheduled job payloads.
// Returns an unsubscribe function.
func (s *InMemoryScheduler) Subscribe(cb func(interface{})) (unsubscribe func()) {
	s.mu.Lock()
	id := s.nextSubID
	s.nextSubID++
	s.subscribers[id] = cb
	s.mu.Unlock()
	return func() {
		s.mu.Lock()
		delete(s.subscribers, id)
		s.mu.Unlock()
	}
}

// Size returns the number of scheduled jobs (for testing/inspection).
func (s *InMemoryScheduler) Size() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.jobs.Len()
}
