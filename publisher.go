package main

import "sync/atomic"
import "sync"

type Publisher struct {
	subscribers            []chan []byte
	toberemovedsubscribers chan chan []byte
	tobeaddedsubscribers   chan chan []byte
	messages               chan []byte
	Publishes              int64
	m                      sync.Mutex
}

func NewPublisher() *Publisher {
	publisher := Publisher{}
	publisher.toberemovedsubscribers = make(chan chan []byte)
	publisher.tobeaddedsubscribers = make(chan chan []byte)
	publisher.messages = make(chan []byte)
	publisher.Publishes = 0
	go publisher.monitorsubscriptions()
	go publisher.monitormessages()
	return &publisher
}

func (p *Publisher) Sub(c chan []byte) {
	p.tobeaddedsubscribers <- c
}

func (p *Publisher) Unsub(c chan []byte) {
	p.toberemovedsubscribers <- c
}

func (p *Publisher) Pub(message []byte) {
	p.messages <- message
}

func (p *Publisher) monitorsubscriptions() {
	var c chan []byte

	for {
		select {
		case c = <-p.toberemovedsubscribers:
			var i = 0
			var v chan []byte
			p.m.Lock()
			for i, v = range p.subscribers {
				if v == c {
					break
				}
			}
			p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
			p.m.Unlock()
		case c = <-p.tobeaddedsubscribers:
			p.m.Lock()
			p.subscribers = append(p.subscribers, c)
			p.m.Unlock()
		}
	}
}

func (p *Publisher) monitormessages() {
	for {
		message := <-p.messages
		atomic.AddInt64(&p.Publishes, 1)
		p.m.Lock()
		for _, c := range p.subscribers {
			c <- message
		}
		p.m.Unlock()
	}
}
