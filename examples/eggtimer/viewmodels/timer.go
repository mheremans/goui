// SPDX-License-Identifier: MIT

package viewmodels

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mheremans/goui"
	"github.com/mheremans/goui/types"
)

type Timer struct {
	goui.ViewModel
	progress      *types.Binding[float32]
	boiling       *types.Binding[bool]
	timeRemaining *types.Binding[string]

	tickerStopChan chan struct{}
	tickerDuration time.Duration
	boilDuration   time.Duration
}

func NewTimer() *Timer {
	t := &Timer{
		progress:      types.NewBinding("Progress", float32(0)),
		boiling:       types.NewBinding("Boiling", false),
		timeRemaining: types.NewBinding("Time Remaining", ""),

		tickerStopChan: make(chan struct{}),
		tickerDuration: time.Millisecond * 40,
		boilDuration:   time.Minute * 5,
	}
	t.RegisterBindings(t.progress, t.boiling, t.timeRemaining)
	return t
}

func (t *Timer) Initialize() (err error) {
	if err = t.ViewModel.Initialize(); err != nil {
		return
	}

	t.boiling.Watch(t)
	go t.tickRoutine()
	return
}

func (t *Timer) Destroy() (err error) {
	t.tickerStopChan <- struct{}{}
	close(t.tickerStopChan)

	return t.ViewModel.Destroy()
}

func (t Timer) Progress() *types.Binding[float32] {
	return t.progress
}

func (t Timer) Boiling() *types.Binding[bool] {
	return t.boiling
}

func (t Timer) TimeRemaining() *types.Binding[string] {
	return t.timeRemaining
}

func (t *Timer) ToggleBoiling() {
	t.boiling.Set(!t.boiling.Get())
	if !t.boiling.Get() {
		t.reset()
		return
	}

	minutes, err := strconv.ParseFloat(t.timeRemaining.Get(), 64)
	fmt.Println(minutes)
	if err != nil {
		t.reset()
		return
	}
	t.boilDuration = time.Duration(float64(time.Minute) * minutes)
}

func (t *Timer) BindingChanged(binding types.Bindable) {
	switch binding.Name() {
	case "Boiling":
		binding := binding.(*types.Binding[bool])
		fmt.Printf("Boiling: %t\n", binding.Get())
	}
}

func (t *Timer) reset() {
	t.boiling.Set(false)
	t.progress.Set(0)
	t.timeRemaining.Set("")
}

func (t *Timer) tickRoutine() {
	ticker := time.NewTicker(40 * time.Millisecond)

	for {
		select {
		case <-t.tickerStopChan:
			return
		case <-ticker.C:
			if t.boiling.Get() {
				if t.progress.Get() < 1 {
					inc := t.tickerDuration.Seconds() / t.boilDuration.Seconds()
					t.progress.Set(t.progress.Get() + float32(inc))
					continue
				}
				t.reset()
				continue
			}
		}
	}
}
