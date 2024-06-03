// SPDX-License-Identifier: MIT

package widget

import (
	giolayout "gioui.org/layout"
	giowidget "gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mheremans/goui/definition"
	"github.com/mheremans/goui/types"
)

type CheckBox struct {
	*Widget

	checkbox *material.CheckBoxStyle
	check    giowidget.Bool
	binding  *types.Binding[bool]

	OnHovered      OnHoveredFn
	OnHoverEntered OnHoverEnteredFn
	OnHoverExited  OnHoverExitedFn
	OnPressed      OnPressedFn
	OnPressDown    OnPressDownFn
	OnPressUp      OnPressUpFn

	prevHoverState bool
	prevPressState bool
}

func NewCheckBox(
	ctx types.Context,
	label string,
	value bool,
	id ...string,
) *CheckBox {
	c := new(CheckBox)
	c.Widget = NewWidget(ctx.Window(), id...)
	checkBox := material.CheckBox(ctx.Window().Theme(), &c.check, label)
	c.checkbox = &checkBox
	c.check.Value = value
	return c
}

func newCheckBoxFromDefinition(
	ctx types.Context,
	data map[string]any,
) (definition.DefinitionType, error) {
	id, _ := definition.MapValueString[string](data, "id")
	label, _ := definition.MapValueString[string](data, "label")
	value, _ := definition.MapValueBool[bool](data, "value")
	c := NewCheckBox(ctx, label, value, id)
	if binding, ok := definition.BindingFromMap[*types.Binding[bool]](ctx, data, "binding"); ok {
		c.Bind(binding)
	}
	c.OnHovered, _ = definition.FunctionFromMap[OnHoveredFn](
		ctx, data, "onHovered")
	c.OnHoverEntered, _ = definition.FunctionFromMap[OnHoverEnteredFn](
		ctx, data, "onHoverEntered")
	c.OnHoverExited, _ = definition.FunctionFromMap[OnHoverExitedFn](
		ctx, data, "onHoverExited")
	c.OnPressed, _ = definition.FunctionFromMap[OnPressedFn](
		ctx, data, "onPressed")
	c.OnPressDown, _ = definition.FunctionFromMap[OnPressDownFn](
		ctx, data, "onPressDown")
	c.OnPressUp, _ = definition.FunctionFromMap[OnPressUpFn](
		ctx, data, "onPressUp")

	return c, nil
}

func (c *CheckBox) Bind(binding *types.Binding[bool]) {
	if c.binding != nil {
		c.binding.Unwatch(c)
		c.binding = nil
	}

	if binding == nil {
		return
	}

	c.binding = binding
	c.binding.Watch(c)
}

func (c CheckBox) Value() bool {
	return c.check.Value
}

func (c CheckBox) Label() string {
	return c.checkbox.Label
}

func (c *CheckBox) SetLabel(label string) {
	c.checkbox.Label = label
	c.Wnd().Invalidate()
}

func (c *CheckBox) SetValue(value bool) {
	c.check.Value = value
	c.Wnd().Invalidate()
}

func (c *CheckBox) HandleEvents(ctx types.Context) {
	if c.binding != nil {
		c.binding.Set(c.check.Value)
	}

	pressed := c.check.Pressed()
	hovered := c.check.Hovered()

	if c.OnPressDown != nil && pressed && !c.prevPressState {
		c.OnPressDown(ctx, c)
	}
	if c.OnPressUp != nil && !pressed && c.prevPressState {
		c.OnPressUp(ctx, c)
	}
	if c.OnPressed != nil && pressed {
		c.OnPressed(ctx, c)
	}
	c.prevPressState = pressed

	if c.OnHoverEntered != nil && hovered && !c.prevHoverState {
		c.OnHoverEntered(ctx, c)
	}
	if c.OnHoverExited != nil && !hovered && c.prevHoverState {
		c.OnHoverExited(ctx, c)
	}
	if c.OnHovered != nil && hovered {
		c.OnHovered(ctx, c)
	}
	c.prevHoverState = hovered
}

func (c *CheckBox) Draw(gtx giolayout.Context) giolayout.Dimensions {
	return c.checkbox.Layout(gtx)
}

func (c *CheckBox) BindingChanged(binding types.Bindable) {
	if bnd, ok := binding.(*types.Binding[bool]); ok {
		c.SetValue(bnd.Get())
	}
}
