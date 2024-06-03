// SPDX-License-Identifier: MIT

package widget

import (
	giolayout "gioui.org/layout"
	giowidget "gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mheremans/goui/definition"
	"github.com/mheremans/goui/icons"
	"github.com/mheremans/goui/types"
)

func init() {
	definition.RegisterUIElement((*Button)(nil), newButtonFromDefinition)
	definition.RegisterUIElement((*IconButton)(nil), newIconButtonFromDefinition)
}

// Button is a clickable button
//
// Yaml	definition:
//
//	type: widget.Button
//	id: <string>				# id of the element (used to get a reference to
//								# it in code)
//	label: <string>				# button label
//	onClicked: <string>			# clicked the button (will be called when the
//								# button is clicked)
//	onHovered: <string>			# hovering over the button (will be called for
//								# as long as the mouse is hovering over the
//								# button)
//	onHoverEntered: <string>	# entered the button area (will be called when
//								# the mouse starts hovering over the button)
//	onHoverExited: <string>		# exited the button area (will be called when
//								# the mouse stops hovering over the button)
//	onPressed					# pressing the button (will be called for as
//								# long as the mouse is pressing the button)
//	onPressDown					# pressed button down (will be called when the
//								# mouse started pressing the button down)
//	onPressUp					# released button (will be called when the mouse
//								# stops pressing the button down)
type Button struct {
	*Widget

	button    *material.ButtonStyle
	clickable giowidget.Clickable

	OnClicked      OnClickedFn
	OnHovered      OnHoveredFn
	OnHoverEntered OnHoverEnteredFn
	OnHoverExited  OnHoverExitedFn
	OnPressed      OnPressedFn
	OnPressDown    OnPressDownFn
	OnPressUp      OnPressUpFn

	prevHoverState bool
	prevPressState bool
}

func NewButton(ctx types.Context, label string, id ...string) *Button {
	b := new(Button)
	b.Widget = NewWidget(ctx.Window(), id...)
	button := material.Button(ctx.Window().Theme(), &b.clickable, label)
	b.button = &button
	return b
}

func newButtonFromDefinition(
	ctx types.Context,
	data map[string]any,
) (definition.DefinitionType, error) {
	id, _ := definition.MapValueString[string](data, "id")
	label, _ := definition.MapValueString[string](data, "label")
	b := NewButton(ctx, label, id)
	b.OnClicked, _ = definition.FunctionFromMap[OnClickedFn](
		ctx, data, "onClicked")
	b.OnHovered, _ = definition.FunctionFromMap[OnHoveredFn](
		ctx, data, "onHovered")
	b.OnHoverEntered, _ = definition.FunctionFromMap[OnHoverEnteredFn](
		ctx, data, "onHoverEntered")
	b.OnHoverExited, _ = definition.FunctionFromMap[OnHoverExitedFn](
		ctx, data, "onHoverExited")
	b.OnPressed, _ = definition.FunctionFromMap[OnPressedFn](
		ctx, data, "onPressed")
	b.OnPressDown, _ = definition.FunctionFromMap[OnPressDownFn](
		ctx, data, "onPressDown")
	b.OnPressUp, _ = definition.FunctionFromMap[OnPressUpFn](
		ctx, data, "onPressUp")

	return b, nil
}

func (b Button) Label() string {
	return b.button.Text
}

func (b *Button) SetLabel(label string) {
	b.button.Text = label
	b.Wnd().Invalidate()
}

func (b *Button) HandleEvents(ctx types.Context) {
	pressed := b.clickable.Pressed()
	hovered := b.clickable.Hovered()
	clicked := b.clickable.Clicked(ctx.Gtx())

	if b.OnPressDown != nil && pressed && !b.prevPressState {
		b.OnPressDown(ctx, b)
	}
	if b.OnPressUp != nil && !pressed && b.prevPressState {
		b.OnPressUp(ctx, b)
	}
	if b.OnPressed != nil && pressed {
		b.OnPressed(ctx, b)
	}
	b.prevPressState = pressed

	if b.OnHoverEntered != nil && hovered && !b.prevHoverState {
		b.OnHoverEntered(ctx, b)
	}
	if b.OnHoverExited != nil && !hovered && b.prevHoverState {
		b.OnHoverExited(ctx, b)
	}
	if b.OnHovered != nil && hovered {
		b.OnHovered(ctx, b)
	}
	b.prevHoverState = hovered

	if b.OnClicked != nil && clicked {
		b.OnClicked(ctx, b)
	}
}

func (b *Button) Draw(gtx giolayout.Context) giolayout.Dimensions {
	return b.button.Layout(gtx)
}

func (b *Button) AddChild(_ types.UIElement, _ ...float32) bool {
	return false
}

// Button is a clickable icon button
//
// Yaml	definition:
//
//	type: widget.Button
//	id: <string>				# id of the element (used to get a reference to
//								# it in code)
//	icon: <string>				# button icon
//	descriptoin: <string>		# button description
//	onClicked: <string>			# clicked the button (will be called when the
//								# button is clicked)
//	onHovered: <string>			# hovering over the button (will be called for
//								# as long as the mouse is hovering over the
//								# button)
//	onHoverEntered: <string>	# entered the button area (will be called when
//								# the mouse starts hovering over the button)
//	onHoverExited: <string>		# exited the button area (will be called when
//								# the mouse stops hovering over the button)
//	onPressed					# pressing the button (will be called for as
//								# long as the mouse is pressing the button)
//	onPressDown					# pressed button down (will be called when the
//								# mouse started pressing the button down)
//	onPressUp					# released button (will be called when the mouse
//								# stops pressing the button down)
type IconButton struct {
	*Widget

	button    *material.IconButtonStyle
	clickable giowidget.Clickable

	OnClicked      OnClickedFn
	OnHovered      OnHoveredFn
	OnHoverEntered OnHoverEnteredFn
	OnHoverExited  OnHoverExitedFn
	OnPressed      OnPressedFn
	OnPressDown    OnPressDownFn
	OnPressUp      OnPressUpFn

	prevHoverState bool
	prevPressState bool
}

func NewIconButton(
	ctx types.Context,
	icon string,
	description string,
	id ...string,
) *IconButton {
	b := new(IconButton)
	button := material.IconButton(
		ctx.Window().Theme(), &b.clickable, icons.Icon(icon), description)
	b.Widget = NewWidget(ctx.Window(), id...)
	b.wnd = ctx.Window()
	b.button = &button
	return b
}

func newIconButtonFromDefinition(
	ctx types.Context,
	data map[string]any,
) (definition.DefinitionType, error) {
	id, _ := definition.MapValueString[string](data, "id")
	icon, _ := definition.MapValueString[string](data, "icon")
	description, _ := definition.MapValueString[string](data, "description")
	b := NewIconButton(ctx, icon, description, id)
	b.OnClicked, _ = definition.FunctionFromMap[OnClickedFn](
		ctx, data, "onClicked")
	b.OnHovered, _ = definition.FunctionFromMap[OnHoveredFn](
		ctx, data, "onHovered")
	b.OnHoverEntered, _ = definition.FunctionFromMap[OnHoverEnteredFn](
		ctx, data, "onHoverEntered")
	b.OnHoverExited, _ = definition.FunctionFromMap[OnHoverExitedFn](
		ctx, data, "onHoverExited")
	b.OnPressed, _ = definition.FunctionFromMap[OnPressedFn](
		ctx, data, "onPressed")
	b.OnPressDown, _ = definition.FunctionFromMap[OnPressDownFn](
		ctx, data, "onPressDown")
	b.OnPressUp, _ = definition.FunctionFromMap[OnPressUpFn](
		ctx, data, "onPressUp")
	return b, nil
}

func (b IconButton) Description() string {
	return b.button.Description
}

func (b *IconButton) SetIcon(icon string) {
	b.button.Icon = icons.Icon(icon)
	b.wnd.Invalidate()
}

func (b *IconButton) SetDescription(description string) {
	b.button.Description = description
	b.wnd.Invalidate()
}

func (b *IconButton) HandleEvents(ctx types.Context) {
	pressed := b.clickable.Pressed()
	hovered := b.clickable.Hovered()
	clicked := b.clickable.Clicked(ctx.Gtx())

	if b.OnPressDown != nil && pressed && !b.prevPressState {
		b.OnPressDown(ctx, b)
	}
	if b.OnPressUp != nil && !pressed && b.prevPressState {
		b.OnPressUp(ctx, b)
	}
	if b.OnPressed != nil && pressed {
		b.OnPressed(ctx, b)
	}
	b.prevPressState = pressed

	if b.OnHoverEntered != nil && hovered && !b.prevHoverState {
		b.OnHoverEntered(ctx, b)
	}
	if b.OnHoverExited != nil && !hovered && b.prevHoverState {
		b.OnHoverExited(ctx, b)
	}
	if b.OnHovered != nil && hovered {
		b.OnHovered(ctx, b)
	}
	b.prevHoverState = hovered

	if b.OnClicked != nil && clicked {
		b.OnClicked(ctx, b)
	}
}

func (b *IconButton) Draw(gtx giolayout.Context) giolayout.Dimensions {
	return b.button.Layout(gtx)
}
