// SPDX-License-Identifier: MIT

package layout

import (
	giolayout "gioui.org/layout"
	"gioui.org/unit"
	"github.com/mheremans/goui/definition"
	"github.com/mheremans/goui/types"
)

func init() {
	definition.RegisterUIElement((*Inset)(nil), newInsetFromDefinition)
}

// Inset that adds padding around a child.
//
// Yaml definition:
//
//	type: layout.Inset
//	id: <string>		# id of the element (used to get a reference to it in code)
//	top: <number>		# top padding (in Dp units)
//	bottom: <number>	# bottom padding (in Dp units)
//	left: <number>		# left padding (in Dp units)
//	right: <number>		# right padding (in Dp units)
//	child: {}			# child element
type Inset struct {
	*Layout

	inset *giolayout.Inset
	child types.UIElement
}

func NewInset[T types.SizeConstraint](
	ctx types.Context,
	top, bottom, left, right T,
	child ...types.UIElement,
) *Inset {
	return &Inset{
		Layout: NewLayout(ctx.Window()),
		inset: &giolayout.Inset{
			Top:    unit.Dp(top),
			Bottom: unit.Dp(bottom),
			Left:   unit.Dp(left),
			Right:  unit.Dp(right),
		},
		child: func() types.UIElement {
			if len(child) == 0 {
				return nil
			}
			return child[0]
		}(),
	}
}

func newInsetFromDefinition(
	ctx types.Context,
	data map[string]any,
) (definition.DefinitionType, error) {
	id, _ := definition.MapValueString[string](data, "id")
	top, _ := definition.MapValueFloat[unit.Dp](data, "top")
	bottom, _ := definition.MapValueFloat[unit.Dp](data, "bottom")
	left, _ := definition.MapValueFloat[unit.Dp](data, "left")
	right, _ := definition.MapValueFloat[unit.Dp](data, "right")
	i := NewInset(ctx, top, bottom, left, right)
	i.SetID(id)
	return i, nil
}

func (i *Inset) SetChild(child types.UIElement) {
	i.child = child
}

func (i *Inset) AddChild(child types.UIElement, _ ...float32) bool {
	i.SetChild(child)
	return true
}

func (i *Inset) HandleEvents(ctx types.Context) {
	if i.child != nil {
		i.child.HandleEvents(ctx)
	}
}

func (i *Inset) Draw(gtx giolayout.Context) giolayout.Dimensions {
	return i.inset.Layout(gtx, i.child.Draw)
}
