package stopwatch

import "time"

type Stopwatch struct {
	start   time.Time
	stop    time.Time
	running bool
}

// NewStopwatch creates a new, unstarted stopwatch.
func NewStopwatch() *Stopwatch {
	return &Stopwatch{}
}

// Start begins the stopwatch. If already running, it does nothing.
func (s *Stopwatch) Start() {
	if !s.running {
		s.start = time.Now()
		s.running = true
	}
}

// Stop halts the stopwatch. If not running, it does nothing.
func (s *Stopwatch) Stop() {
	if s.running {
		s.stop = time.Now()
		s.running = false
	}
}

// Elapsed returns the duration since the stopwatch started.
// If the stopwatch is currently running, it calculates the duration up to the current time.
// If stopped, it returns the duration between start and stop.
func (s *Stopwatch) Elapsed() time.Duration {
	if s.running {
		return time.Since(s.start)
	}
	return s.stop.Sub(s.start)
}

// Reset clears the stopwatch, setting its state back to unstarted.
func (s *Stopwatch) Reset() {
	s.start = time.Time{}
	s.stop = time.Time{}
	s.running = false
}
