package main

import "sync/atomic"

type Publisher struct {
	subscribers            []chan []byte
	toberemovedsubscribers chan chan []byte
	tobeaddedsubscribers chan chan []byte
	Publishes              int64
}

func NewPublisher() *Publisher {
    publisher := Publisher{}
    publisher.toberemovedsubscribers = make(chan chan []byte)
    publisher.tobeaddedsubscribers = make(chan chan []byte)
    publisher.Publishes = 0
    go publisher.monitorsubscription()
    return &publisher
}

func (p *Publisher) Sub(c chan []byte) {
	p.tobeaddedsubscribers <- c
}

func (p *Publisher) Unsub(c chan []byte) {
	p.toberemovedsubscribers <- c
}

func (p *Publisher) Pub(message []byte) {
	atomic.AddInt64(&p.Publishes,1)
	for _, c := range p.subscribers {
		c <- message
	}
}

func (p *Publisher) monitorsubscription() {
    var c chan []byte

	for {
        select {
        case c = <-p.toberemovedsubscribers:
            var i = 0
            var v chan []byte
            for i, v = range p.subscribers {
                if v == c {
                    break
                }
            }
            p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
        case c = <-p.tobeaddedsubscribers:
            p.subscribers = append(p.subscribers, c)
        }
	}
}
