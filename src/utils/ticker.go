package utils

import (
	"time"
	"sync"
)

//Copied from DocSense ^^

//Ticker fires a function in background after specified duration of inactivity. Restart() is used to reset the countdown.
//Only one instance of function is allowed to run at any given time.
type Ticker struct {
	lastRestart time.Time
	fireAfter   time.Duration
	fireFunc    func() ()
	tick        *time.Ticker
	end         chan bool
	active      bool
	mutex       sync.RWMutex
}

//Ticker fires a function in background after specified duration of inactivity. Restart() is used to reset the countdown.
//Only one instance of function is allowed to run at any given time.
func New(fireAfter time.Duration, fireFunc func() ()) *Ticker {
	return &Ticker{
		lastRestart: time.Now(),
		fireAfter:   fireAfter,
		fireFunc:    fireFunc,
		active:      false,
	}
}

func (t *Ticker) start() {
	t.end = make(chan bool)
	t.tick = time.NewTicker(t.fireAfter)
	t.active = true
	go func() {
		select {
		case <-t.tick.C:
			t.fire()
			return

		case <-t.end:
			return
		}
	}()
}

//Restart should be called every time, request has been received
func (t *Ticker) Restart() {
	t.lastRestart = time.Now()
	if t.active {
		t.Stop()
		t.start()
	} else {
		t.start()
	}
}

//Stop the executions.
func (t *Ticker) Stop() {
	t.end <- true
	t.active = false
	t.tick.Stop()
}

func (t *Ticker) fire() {
	go func() {
		t.mutex.Lock()
		t.fireFunc()
		t.mutex.Unlock()
	}()
	t.active = false
	t.tick.Stop()
}
