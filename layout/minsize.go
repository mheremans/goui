// SPDX-License-Identifier: MIT

package layout

import (
	giolayout "gioui.org/layout"
	"gioui.org/unit"
	"github.com/mheremans/goui/definition"
	"github.com/mheremans/goui/types"
)

func init() {
	definition.RegisterUIElement((*MinSize)(nil), newMinSizeFromDefinition)
}

type MinSize struct {
	*Layout

	minWidth  unit.Dp
	minHeight unit.Dp
	child     types.UIElement
}

func NewMinSize[T types.SizeConstraint](
	ctx types.Context,
	minWidth, minHeight T,
	child ...types.UIElement,
) *MinSize {
	return &MinSize{
		Layout:    NewLayout(ctx.Window()),
		minWidth:  unit.Dp(minWidth),
		minHeight: unit.Dp(minHeight),
		child: func() types.UIElement {
			if len(child) == 0 {
				return nil
			}
			return child[0]
		}(),
	}
}

func newMinSizeFromDefinition(
	ctx types.Context,
	data map[string]any,
) (definition.DefinitionType, error) {
	id, _ := definition.MapValueString[string](data, "id")
	minWidth, _ := definition.MapValueFloat[unit.Dp](data, "minWidth")
	minHeight, _ := definition.MapValueFloat[unit.Dp](data, "minHeight")
	ms := NewMinSize(ctx, minWidth, minHeight)
	ms.SetID(id)
	return ms, nil
}

func (ms *MinSize) SetChild(child types.UIElement) {
	ms.child = child
}

func (ms *MinSize) AddChild(child types.UIElement, _ ...float32) bool {
	ms.SetChild(child)
	return true
}

func (ms *MinSize) HandleEvents(ctx types.Context) {
	if ms.child != nil {
		ms.child.HandleEvents(ctx)
	}
}

func (ms *MinSize) Draw(gtx giolayout.Context) giolayout.Dimensions {
	if ms.minWidth > 0 && gtx.Constraints.Min.X < int(ms.minWidth) {
		gtx.Constraints.Min.X = int(ms.minWidth)
		if gtx.Constraints.Max.X < gtx.Constraints.Min.X {
			gtx.Constraints.Max.X = gtx.Constraints.Min.X
		}
	}
	if ms.minHeight > 0 && gtx.Constraints.Min.Y < int(ms.minHeight) {
		gtx.Constraints.Min.Y = int(ms.minHeight)
		if gtx.Constraints.Max.Y < gtx.Constraints.Min.Y {
			gtx.Constraints.Max.Y = gtx.Constraints.Min.Y
		}
	}
	return ms.child.Draw(gtx)

	/*dim := ms.child.Draw(gtx)
	if ms.minWidth > 0 && unit.Dp(dim.Size.X) < ms.minWidth {
		dim.Size.X = int(ms.minWidth)
	}
	if ms.minHeight > 0 && unit.Dp(dim.Size.Y) < ms.minHeight {
		dim.Size.Y = int(ms.minHeight)
	}
	return dim*/
}
