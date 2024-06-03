// SPDX-License-Identifier: MIT

package views

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	giolayout "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/mheremans/goui"
	"github.com/mheremans/goui/examples/eggtimer/viewmodels"
	"github.com/mheremans/goui/types"
	"github.com/mheremans/goui/widget"
)

type TimerView struct {
	*goui.View

	input       *widget.Input
	progressBar *widget.ProgressBar
	button      *widget.Button
}

func NewTimerView() *TimerView {
	v := new(TimerView)
	v.View = goui.ConfigureView(
		v, viewmodels.NewTimer(),
		goui.NewViewScreen(Definitions, "def/timerview.def.yml"),
	)
	v.ExportFunction("drawEgg", v.drawEgg)
	v.ExportFunction("onButtonStartClicked", v.onButtonStartClicked)
	return v
}

func (v *TimerView) Initialize(ctx types.Context) (err error) {
	if err = v.View.Initialize(ctx); err != nil {
		err = fmt.Errorf("failed to initialize view: %w", err)
		return
	}

	v.FindBinding("Boiling").Watch(v)
	v.FindBinding("Time Remaining").Watch(v)
	v.input, _ = goui.GetElementById[*widget.Input](v.View, "timeInput")
	v.progressBar, _ = goui.GetElementById[*widget.ProgressBar](v.View, "progressBar")
	v.button, _ = goui.GetElementById[*widget.Button](v.View, "startButton")

	return
}

func (v *TimerView) onButtonStartClicked(ctx types.Context, button types.UIElement) {
	v.ViewModel().(*viewmodels.Timer).ToggleBoiling()
}

func (s *TimerView) BindingChanged(binding types.Bindable) {
	switch binding.Name() {
	case "Boiling":
		binding := binding.(*types.Binding[bool])
		s.UpdateButtonLabel(binding.Get())
	case "Time Remaining":
		fmt.Printf("Time Remaining: %s\n", binding.(*types.Binding[string]).Get())
	}
}

func (s *TimerView) UpdateButtonLabel(boiling bool) {
	if boiling {
		s.button.SetLabel("Stop")
	} else {
		s.button.SetLabel("Start")
	}
}

func (v *TimerView) drawEgg(gtx giolayout.Context, graphic types.UIElement) image.Point {
	progress, ok := v.FindBinding("Progress").(*types.Binding[float32])
	if !ok {
		return image.Point{}
	}

	var eggPath clip.Path
	op.Offset(image.Pt(gtx.Constraints.Max.X/2, gtx.Constraints.Max.Y/4)).Add(gtx.Ops)
	eggPath.Begin(gtx.Ops)

	max := float32(0.0)
	min := float32(gtx.Constraints.Max.Y)

	for deg := 0.0; deg <= 360.0; deg++ {
		rad := deg / 360 * 2 * math.Pi
		cosT := math.Cos(rad)
		sinT := math.Sin(rad)
		a := float64(gtx.Constraints.Max.Y) * 0.2
		b := float64(gtx.Constraints.Max.Y) * 0.3
		d := float64(gtx.Constraints.Max.Y) * 0.04
		x := a * cosT
		y := -(math.Sqrt(b*b-d*d*cosT*cosT) + d*sinT) * sinT
		p := f32.Pt(float32(x), float32(y))
		if p.Y > max {
			max = p.Y
		}
		if p.Y < min {
			min = p.Y
		}
		eggPath.LineTo(p)
	}
	eggPath.Close()
	eggArea := clip.Outline{Path: eggPath.End()}.Op()
	color := color.NRGBA{
		R: 255,
		G: uint8(239 * (1 - progress.Get())),
		B: uint8(174 * (1 - progress.Get())),
		A: 255}
	paint.FillShape(gtx.Ops, color, eggArea)

	return image.Point{Y: int(max - min)}
}
