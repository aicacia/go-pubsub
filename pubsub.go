package pubsub

import "sync"

const default_subscriber_channel_size = 1024

type Subscriber[T any] struct {
	pubsub *PubSub[T]
	C      chan T
}

func (s *Subscriber[T]) Close() {
	s.pubsub.close(s)
}

type PubSub[T any] struct {
	mutex        sync.RWMutex
	subscribers  []*Subscriber[T]
	channel_size int
}

func NewPubSub[T any]() PubSub[T] {
	return PubSub[T]{
		mutex:        sync.RWMutex{},
		channel_size: default_subscriber_channel_size,
	}
}

func (ps *PubSub[T]) SetChannelSize(channel_size int) *PubSub[T] {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	ps.channel_size = channel_size
	return ps
}

func (ps *PubSub[T]) Publish(value T) {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()
	for _, subscriber := range ps.subscribers {
		if len(subscriber.C) < cap(subscriber.C) {
			subscriber.C <- value
		}
	}
}

func (ps *PubSub[T]) Subscribe(channel_size ...int) *Subscriber[T] {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	var cs = ps.channel_size
	if len(channel_size) > 0 {
		cs = channel_size[0]
	}
	subscriber := Subscriber[T]{
		pubsub: ps,
		C:      make(chan T, cs),
	}
	ps.subscribers = append(ps.subscribers, &subscriber)
	return &subscriber
}

func (ps *PubSub[T]) Subscribers() []*Subscriber[T] {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()
	subscribers := make([]*Subscriber[T], len(ps.subscribers))
	copy(subscribers, ps.subscribers)
	return subscribers
}

func (ps *PubSub[T]) close(s *Subscriber[T]) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	close(s.C)
	index := -1
	for i, subscriber := range ps.subscribers {
		if subscriber.C == s.C {
			index = i
		}
	}
	if index != -1 {
		ps.subscribers = append(ps.subscribers[:index], ps.subscribers[index+1:]...)
	}
}
