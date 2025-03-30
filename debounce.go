package main

import "time"

type debouncer struct {
	duration time.Duration
	timer    *time.Timer
}

func (d *debouncer) then(fn func()) {
	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.duration, fn)
}
