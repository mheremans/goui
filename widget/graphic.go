// SPDX-License-Identifier: MIT

package widget

import (
	"fmt"

	giolayout "gioui.org/layout"
	"github.com/mheremans/goui/definition"
	"github.com/mheremans/goui/types"
)

func init() {
	definition.RegisterUIElement((*Graphic)(nil), newGraphicFromDefinition)
}

// Graphic is a widget that draws on the screen by a user supplied draw function.
//
// Yaml definition:
//
//	type: widget.Graphic
//	id: <string>			# id of the element (used to get a reference to it in code)
//	drawFunction: <string>	# draw function reference (will be requested throught the view)
type Graphic struct {
	*Widget

	drawFn GraphicFn
}

func NewGraphic(ctx types.Context, drawFn GraphicFn, id ...string) *Graphic {
	g := new(Graphic)
	g.Widget = NewWidget(ctx.Window(), id...)
	g.drawFn = drawFn
	return g
}

func newGraphicFromDefinition(
	ctx types.Context,
	data map[string]any,
) (definition.DefinitionType, error) {
	id, _ := definition.MapValueString[string](data, "id")
	drawFn, ok := definition.FunctionFromMap[GraphicFn](ctx, data, "drawFunction")
	if !ok {
		return nil, fmt.Errorf("no draw function")
	}
	return NewGraphic(ctx, drawFn, id), nil
}

func (g *Graphic) HandleEvents(ctx types.Context) {
}

func (g *Graphic) Draw(gtx giolayout.Context) giolayout.Dimensions {
	d := g.drawFn(gtx, g)
	return giolayout.Dimensions{Size: d}
}

func (g *Graphic) AddChild(_ types.UIElement, _ ...float32) bool {
	return false
}
