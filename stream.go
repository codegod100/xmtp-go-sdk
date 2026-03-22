package xmtp

import (
	"context"
	"sync"
)

// Stream represents a stream of values that can be iterated
type Stream[T any] struct {
	ch     chan T
	err    error
	done   chan struct{}
	cancel context.CancelFunc
	once   sync.Once
}

// NewStream creates a new stream
func NewStream[T any](ctx context.Context, bufferSize int) *Stream[T] {
	ctx, cancel := context.WithCancel(ctx)
	return &Stream[T]{
		ch:     make(chan T, bufferSize),
		done:   make(chan struct{}),
		cancel: cancel,
	}
}

// Channel returns the underlying channel
func (s *Stream[T]) Channel() <-chan T {
	return s.ch
}

// Push adds a value to the stream
func (s *Stream[T]) Push(value T) {
	select {
	case s.ch <- value:
	case <-s.done:
	}
}

// Error sets an error on the stream
func (s *Stream[T]) Error(err error) {
	s.once.Do(func() {
		s.err = err
		close(s.done)
	})
}

// Close closes the stream
func (s *Stream[T]) Close() {
	s.once.Do(func() {
		close(s.done)
	})
}

// Cancel cancels the stream context
func (s *Stream[T]) Cancel() {
	if s.cancel != nil {
		s.cancel()
	}
}

// Wait waits for the stream to complete
func (s *Stream[T]) Wait() error {
	<-s.done
	return s.err
}

// Collect collects all values from the stream into a slice
func (s *Stream[T]) Collect(ctx context.Context) ([]T, error) {
	var result []T
	for {
		select {
		case v, ok := <-s.ch:
			if !ok {
				return result, s.err
			}
			result = append(result, v)
		case <-s.done:
			return result, s.err
		case <-ctx.Done():
			return result, ctx.Err()
		}
	}
}
