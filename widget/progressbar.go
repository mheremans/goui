// SPDX-License-Identifier: MIT

package widget

import (
	giolayout "gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/mheremans/goui/definition"
	"github.com/mheremans/goui/types"
)

func init() {
	definition.RegisterUIElement((*ProgressBar)(nil), newProgressBarFromDefinition)
}

// ProgressBar is a widget that displays a progress bar.
//
// Yaml definition:
//
//	type: widget.ProgressBar
//	id: <string>		# id of the element (used to get a reference to it in code)
//	value: <number>		# initial value of the progress bar (in range 0.0 - 1.0)
//	binding: <string>	# binding reference (will be requested throught the view)
type ProgressBar struct {
	*Widget

	progressBar *material.ProgressBarStyle
	binding     *types.Binding[float32]
}

// NewProgressBar creates a new ProgressBar widget with the specified context
// and value.
//
// Parameters:
// - ctx: the types.Context used to create the widget.
// - value: the initial value of the progress bar.
//
// Returns:
// - *ProgressBar: a pointer to the newly created ProgressBar widget.
func NewProgressBar(ctx types.Context, value float32, id ...string) *ProgressBar {
	bar := material.ProgressBar(ctx.Window().Theme(), value)

	return &ProgressBar{
		Widget:      NewWidget(ctx.Window(), id...),
		progressBar: &bar,
	}
}

func newProgressBarFromDefinition(
	ctx types.Context,
	data map[string]any,
) (definition.DefinitionType, error) {
	id, _ := definition.MapValueString[string](data, "id")
	value, _ := definition.MapValueFloat[float32](data, "value")
	pb := NewProgressBar(ctx, value, id)
	if binding, ok := definition.BindingFromMap[*types.Binding[float32]](ctx, data, "binding"); ok {
		pb.Bind(binding)
	}
	return pb, nil
}

func (p *ProgressBar) Bind(binding *types.Binding[float32]) {
	if p.binding != nil {
		p.binding.Unwatch(p)
		p.binding = nil
	}

	if binding == nil {
		return
	}

	p.binding = binding
	p.binding.Watch(p)
	p.SetValue(p.binding.Get())
}

// Value returns the current value of the ProgressBar.
func (p ProgressBar) Value() float32 {
	return p.progressBar.Progress
}

// SetValue sets the value of the ProgressBar.
func (p *ProgressBar) SetValue(value float32) {
	p.progressBar.Progress = value
	p.Wnd().Invalidate()
}

// HandleEvents handles the events for the ProgressBar widget.
func (p *ProgressBar) HandleEvents(ctx types.Context) {
	if p.binding != nil {
		p.binding.Set(p.progressBar.Progress)
	}
}

// Draw renders the ProgressBar widget on the provided giolayout.Context.
//
// Parameters:
// - ctx: The giolayout.Context to draw the ProgressBar on.
//
// Returns:
// - giolayout.Dimensions: The dimensions of the drawn ProgressBar.
func (p *ProgressBar) Draw(gtx giolayout.Context) giolayout.Dimensions {
	return p.progressBar.Layout(gtx)
}

func (p *ProgressBar) BindingChanged(binding types.Bindable) {
	if bnd, ok := binding.(*types.Binding[float32]); ok {
		p.SetValue(bnd.Get())
	}
}
