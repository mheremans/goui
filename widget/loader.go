// SPDX-License-Identifier: MIT

package widget

import (
	"image/color"

	giolayout "gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/mheremans/goui/definition"
	"github.com/mheremans/goui/types"
)

func init() {
	definition.RegisterUIElement((*Loader)(nil), newLoaderFromDefinition)
}

type Loader struct {
	*Widget

	loader *material.LoaderStyle
}

func NewLoader(ctx types.Context, id ...string) *Loader {
	loader := material.Loader(ctx.Window().Theme())

	return &Loader{
		Widget: NewWidget(ctx.Window(), id...),
		loader: &loader,
	}
}

func newLoaderFromDefinition(
	ctx types.Context,
	data map[string]any,
) (definition.DefinitionType, error) {
	id, _ := definition.MapValueString[string](data, "id")
	l := NewLoader(ctx, id)

	if color, ok := definition.MapValueColor(data, "color"); ok {
		l.loader.Color = color
	}
	return l, nil
}

func (l Loader) Color() color.NRGBA {
	return l.loader.Color
}

func (l *Loader) SetColor(c color.NRGBA) {
	l.loader.Color = c
}

func (l *Loader) HandleEvents(ctx types.Context) {
	// Do nothing
}

func (l *Loader) Draw(gtx giolayout.Context) giolayout.Dimensions {
	return l.loader.Layout(gtx)
}
