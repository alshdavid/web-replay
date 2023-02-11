package rx

import "sync"

type SubscriberFunc[T any] func(data T)

type Subject[T any] struct {
	m             sync.Mutex
	callbackIndex int
	callbacks     map[int]*SubscriberFunc[T]
}

func NewSubject[T any]() *Subject[T] {
	return &Subject[T]{
		m:             sync.Mutex{},
		callbackIndex: 0,
		callbacks:     map[int]*SubscriberFunc[T]{},
	}
}

func (s *Subject[T]) Next(data T) {
	s.m.Lock()
	for _, cb := range s.callbacks {
		if cb != nil {
			(*cb)(data)
		}
	}
	s.m.Unlock()
}

func (s *Subject[T]) Subscribe(cb SubscriberFunc[T]) *Subscription {
	s.m.Lock()
	s.callbackIndex += 1
	s.callbacks[s.callbackIndex] = &cb
	s.m.Unlock()
	return &Subscription{
		cleanupFunc: func() {
			s.m.Lock()
			s.callbacks[s.callbackIndex] = nil
			s.m.Unlock()
		},
	}
}

type Subscription struct {
	cleanupFunc func()
}

func (s *Subscription) Unsubscribe() {
	s.cleanupFunc()
}
