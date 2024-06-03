// SPDX-License-Identifier: MIT

package widget

import (
	giolayout "gioui.org/layout"
	"gioui.org/widget"
	"github.com/mheremans/goui/definition"
	"github.com/mheremans/goui/types"
)

func init() {
	definition.RegisterUIElement((*List)(nil), newListFromDefinition)
}

type List struct {
	*Widget

	list    widget.List
	binding types.BindableList

	itemEventHandler ListItemEventHandlerFn
	itemRenderer     ListItemRendererFn

	ctx types.Context
}

func NewList(
	ctx types.Context,
	axis giolayout.Axis,
	alignment giolayout.Alignment,
	itemEventHandler ListItemEventHandlerFn,
	itemRenderer ListItemRendererFn,
	id ...string,
) *List {
	l := new(List)
	l.ctx = ctx
	l.Widget = NewWidget(ctx.Window(), id...)
	l.list.Axis = axis
	l.list.Alignment = alignment
	l.itemEventHandler = itemEventHandler
	l.itemRenderer = itemRenderer
	return l
}

func newListFromDefinition(
	ctx types.Context,
	data map[string]any,
) (definition.DefinitionType, error) {
	id, _ := definition.MapValueString[string](data, "id")
	axis, _ := definition.GioConstantFromMap[giolayout.Axis](data, "axis")
	alignment, _ := definition.GioConstantFromMap[giolayout.Alignment](data, "alignment")
	itemEventHandler, _ := definition.FunctionFromMap[ListItemEventHandlerFn](ctx, data, "itemEventHandler")
	itemRenderer, _ := definition.FunctionFromMap[ListItemRendererFn](ctx, data, "itemRenderer")
	scrollToEnd, _ := definition.MapValueBool[bool](data, "scrollToEnd")

	i := NewList(ctx, axis, alignment, itemEventHandler, itemRenderer, id)
	i.list.ScrollToEnd = scrollToEnd

	if binding, ok := definition.BindingFromMap[types.BindableList](
		ctx, data, "binding",
	); ok {
		i.Bind(binding)
	}

	return i, nil
}

func (l *List) Bind(binding types.BindableList) {
	if l.binding != nil {
		l.binding.Unwatch(l)
		l.binding = nil
	}

	if binding == nil {
		return
	}

	l.binding = binding
	l.binding.Watch(l)
}

func (l List) Axis() giolayout.Axis {
	return l.list.Axis
}

func (l List) Alignment() giolayout.Alignment {
	return l.list.Alignment
}

func (l List) ScrollToEnd() bool {
	return l.list.ScrollToEnd
}

func (l *List) SetAxis(axis giolayout.Axis) {
	l.list.Axis = axis
}

func (l *List) SetAlignment(alignment giolayout.Alignment) {
	l.list.Alignment = alignment
}

func (l *List) SetScrollToEnd(scrollToEnd bool) {
	l.list.ScrollToEnd = scrollToEnd
}

func (l *List) HandleEvents(ctx types.Context) {
	// Cache this cycles context, so we can use it in the Draw, where we handle
	// the events of the list items.
	l.ctx = ctx
}

func (l *List) Draw(gtx giolayout.Context) giolayout.Dimensions {
	return l.list.Layout(gtx, l.binding.Size(), func(gtx giolayout.Context, index int) giolayout.Dimensions {
		l.itemEventHandler(l.ctx, index, l.binding)
		return l.itemRenderer(gtx, index, l.binding)
	})
}

func (l *List) BindingChanged(binding types.Bindable) {
	if _, ok := binding.(types.BindableList); ok {
		l.Wnd().Invalidate()
	}
}
