package main

import "sync"

type pubsub struct {
	mu   sync.RWMutex
	subs []chan string
}

func (ps *pubsub) publish(msg string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.subs {
		ch <- msg
	}
}

func (ps *pubsub) subscribe() <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string)
	ps.subs = append(ps.subs, ch)
	return ch
}
