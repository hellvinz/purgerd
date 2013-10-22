package main

type Publisher struct {
	subscribers            []chan []byte
	toberemovedsubscribers chan chan []byte
	unsubmonitored         bool
	Publishes              int
}

func (p *Publisher) Sub(c chan []byte) {
	p.subscribers = append(p.subscribers, c)
	if !p.unsubmonitored {
		go p.monitorunsub()
	}
}

func (p *Publisher) Unsub(c chan []byte) {
	p.toberemovedsubscribers <- c
}

func (p *Publisher) Pub(message []byte) {
	p.Publishes += 1
	for _, c := range p.subscribers {
		c <- message
	}
}

func (p *Publisher) dounsub(c chan []byte) {
	var i = 0
	var v chan []byte
	for i, v = range p.subscribers {
		if v == c {
			break
		}
	}
	p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
}

func (p *Publisher) monitorunsub() {
	p.unsubmonitored = true
	p.toberemovedsubscribers = make(chan chan []byte)
	for {
		c := <-p.toberemovedsubscribers
		p.dounsub(c)
	}
}
