package main

import (
    "testing"
    "time"
)

func TestPublisherSub(t *testing.T){
    publisher := NewPublisher()
    receiver := make(chan []byte)
    publisher.Sub(receiver)
    if publisher.subscribers[0] != receiver {
        t.Errorf("Publisher.Sub should add argument to the internal subscribers slice")
    }
}

func TestPublisherPub(t *testing.T){
    var message []byte
    timeout := make(chan bool, 1)
    wait := make(chan bool, 1)
    publisher := NewPublisher()
    receiver := make(chan []byte)
    publisher.Sub(receiver)
    go func(){
        message = <-receiver
        wait<-true
    }()
    go func() {
        time.Sleep(1 * time.Second)
        timeout <- true
    }()
    publisher.Pub([]byte("ohai"))
    select {
    case <-wait:
        if string(message) != "ohai" {
            t.Errorf("wrong message, expected %s got %s", "ohai", string(message))
        }
    case <-timeout:
        t.Errorf("no message received after calling Publisher.Pub after 1 second")
    }
}

func TestPublisherUnsub(t *testing.T){
    publisher := NewPublisher()
    receiver := make(chan []byte)
    publisher.Sub(receiver)
    time.Sleep(100 * time.Nanosecond)
    publisher.Unsub(receiver)
    time.Sleep(100 * time.Nanosecond)
    if len(publisher.subscribers) != 0 {
      t.Errorf("subscriber should be removed when calling Publisher.Unsub")
    }
}
