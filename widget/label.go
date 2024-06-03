// SPDX-License-Identifier: MIT

package widget

import (
	giolayout "gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/mheremans/goui/definition"
	"github.com/mheremans/goui/types"
)

type LabelFormat uint8

const (
	Text LabelFormat = iota
	H1
	H2
	H3
	H4
	H5
	H6
	Subtitle1
	Subtitle2
	Body1
	Body2
	Caption
	Overline
)

func (f LabelFormat) String() string {
	return [...]string{
		"Text",
		"H1",
		"H2",
		"H3",
		"H4",
		"H5",
		"H6",
		"Subtitle1",
		"Subtitle2",
		"Body1",
		"Body2",
		"Caption",
		"Overline"}[f]
}

func (f LabelFormat) initializer() func(*material.Theme, string) material.LabelStyle {
	switch f {
	case H1:
		return material.H1
	case H2:
		return material.H2
	case H3:
		return material.H3
	case H4:
		return material.H4
	case H5:
		return material.H5
	case H6:
		return material.H6
	case Subtitle1:
		return material.Subtitle1
	case Subtitle2:
		return material.Subtitle2
	case Body1:
		return material.Body1
	case Body2:
		return material.Body2
	case Caption:
		return material.Caption
	case Overline:
		return material.Overline
	default:
		return func(th *material.Theme, txt string) material.LabelStyle {
			return material.Label(th, th.TextSize, txt)
		}
	}
}

func LabelFormatFromString(s string) LabelFormat {
	switch s {
	case "H1":
		return H1
	case "H2":
		return H2
	case "H3":
		return H3
	case "H4":
		return H4
	case "H5":
		return H5
	case "H6":
		return H6
	case "Subtitle1":
		return Subtitle1
	case "Subtitle2":
		return Subtitle2
	case "Body1":
		return Body1
	case "Body2":
		return Body2
	case "Caption":
		return Caption
	case "Overline":
		return Overline
	default:
		return Text
	}
}

func init() {
	definition.RegisterUIElement((*Label)(nil), newLabelFromDefinition)
}

type Label struct {
	*Widget

	label *material.LabelStyle

	binding *types.Binding[string]
}

func NewLabel(ctx types.Context, txt string, format LabelFormat, id ...string) *Label {
	l := new(Label)
	l.Widget = NewWidget(ctx.Window(), id...)
	label := format.initializer()(ctx.Window().Theme(), txt)
	l.label = &label

	return l
}

func newLabelFromDefinition(
	ctx types.Context,
	data map[string]any,
) (definition.DefinitionType, error) {
	id, _ := definition.MapValueString[string](data, "id")
	txt, _ := definition.MapValueString[string](data, "text")
	labelFormat, _ := definition.MapValueString[string](data, "labelFormat")
	alignment, _ := definition.GioConstantFromMap[text.Alignment](data, "alignment")
	maxLines, _ := definition.MapValueInt[int](data, "maxLines")
	wrapPolicy, _ := definition.MapValueString[string](data, "wrapPolicy")
	truncator, _ := definition.MapValueString[string](data, "truncator")
	lineHeight, _ := definition.MapValueFloat[unit.Sp](data, "lineHeight")
	lineHeightScale, _ := definition.MapValueFloat[float32](data, "lineHeightScale")

	l := NewLabel(ctx, txt, LabelFormatFromString(labelFormat), id)

	if font, ok := definition.MapValueFont(ctx, data, "font", "fontStyle", "fontWeight"); ok {
		l.label.Font = font.Font
	}

	if color, ok := definition.MapValueColor(data, "color"); ok {
		l.label.Color = color
	}
	if selectionColor, ok := definition.MapValueColor(data, "selectionColor"); ok {
		l.label.SelectionColor = selectionColor
	}

	l.label.Alignment = alignment
	l.label.MaxLines = maxLines
	l.label.WrapPolicy = gioTextWrapPolicyFromString(wrapPolicy)
	l.label.Truncator = truncator
	l.label.LineHeight = lineHeight
	l.label.LineHeightScale = lineHeightScale

	if binding, ok := definition.BindingFromMap[*types.Binding[string]](
		ctx, data, "binding",
	); ok {
		l.Bind(binding)
	}

	return l, nil
}

func (l *Label) Bind(binding *types.Binding[string]) {
	if l.binding != nil {
		l.binding.Unwatch(l)
		l.binding = nil
	}

	if binding == nil {
		return
	}

	l.binding = binding
	l.binding.Watch(l)
	l.SetText(l.binding.Get())
}

func (l Label) Text() string {
	return l.label.Text
}

func (l Label) TextSize() unit.Sp {
	return l.label.TextSize
}

func (l *Label) SetText(txt string) {
	l.label.Text = txt
	l.Wnd().Invalidate()
}

func (l *Label) SetTextSize(size unit.Sp) {
	l.label.TextSize = size
	l.Wnd().Invalidate()
}

func (l *Label) HandleEvents(ctx types.Context) {
	if l.binding != nil {
		l.binding.Set(l.Text())
	}
}

func (l *Label) Draw(gtx giolayout.Context) giolayout.Dimensions {
	return l.label.Layout(gtx)
}

func (l *Label) BindingChanged(binding types.Bindable) {
	if bnd, ok := binding.(*types.Binding[string]); ok {
		if l.label.Text != bnd.Get() {
			l.label.Text = bnd.Get()
			l.Wnd().Invalidate()
		}
	}
}
