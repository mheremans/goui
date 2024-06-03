// SPDX-License-Identifier: MIT

package widget

import (
	giolayout "gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mheremans/goui/definition"
	"github.com/mheremans/goui/types"
)

func init() {
	definition.RegisterUIElement((*Slider)(nil), newSliderFromDefinition)
}

type Slider struct {
	*Widget

	slider *material.SliderStyle
	float  widget.Float

	binding *types.Binding[float32]
}

func NewSlider(ctx types.Context, id ...string) *Slider {
	i := new(Slider)
	i.Widget = NewWidget(ctx.Window(), id...)
	slider := material.Slider(ctx.Window().Theme(), &i.float)
	i.slider = &slider

	return i
}

func newSliderFromDefinition(
	ctx types.Context,
	data map[string]any,
) (definition.DefinitionType, error) {
	id, _ := definition.MapValueString[string](data, "id")
	axis, _ := definition.GioConstantFromMap[giolayout.Axis](data, "axis")
	color, _ := definition.MapValueColor(data, "color")

	slider := NewSlider(ctx, id)
	slider.slider.Axis = axis
	slider.slider.Color = color

	if binding, ok := definition.BindingFromMap[*types.Binding[float32]](ctx, data, "binding"); ok {
		slider.Bind(binding)
	}
	return slider, nil
}

func (s *Slider) Bind(binding *types.Binding[float32]) {
	if s.binding != nil {
		s.binding.Unwatch(s)
		s.binding = nil
	}

	if binding == nil {
		return
	}

	s.binding = binding
	s.binding.Watch(s)
	s.SetValue(s.binding.Get())
}

func (s Slider) Value() float32 {
	return s.float.Value
}

func (s *Slider) SetValue(value float32) {
	s.float.Value = value
	s.Wnd().Invalidate()
}

func (s *Slider) HandleEvents(ctx types.Context) {
	if s.binding != nil {
		s.binding.Set(s.float.Value)
	}
}

func (s *Slider) Draw(gtx giolayout.Context) giolayout.Dimensions {
	return s.slider.Layout(gtx)
}

func (s *Slider) BindingChanged(binding types.Bindable) {
	if bnd, ok := binding.(*types.Binding[float32]); ok {
		if s.slider.Float.Value != bnd.Get() {
			s.slider.Float.Value = bnd.Get()
			s.Wnd().Invalidate()
		}
	}
}
