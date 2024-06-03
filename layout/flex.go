// SPDX-License-Identifier: MIT

package layout

import (
	giolayout "gioui.org/layout"
	"github.com/mheremans/goui/definition"
	"github.com/mheremans/goui/types"
)

func init() {
	definition.RegisterUIElement((*Flex)(nil), newFlexFromDefinition)
}

// Flex is a layout that arranges its children according to the FlexBox principle.
//
// Yaml definition:
//
//	type: layout.Flex
//	id: <string>		# id of the element (used to get a reference to it in code)
//	axis: <string>		# gio layout.Axis: "Horizontal" or "Vertical"
//	spacing: <string>	# gio layout.Spacing: "SpaceEnd", "SpaceStart", "SpaceSides", "SpaceAround", "SpaceBetween" or "SpaceEvenly"
//	alignment: <string> # gio layout.Alignment: "Start", "End", "Middle" or "Baseline"
//	children: [{}]		# list of child elements
type Flex struct {
	*Layout

	flex     *giolayout.Flex
	children []*flexChild
}

// NewFlex creates a new Flex layout with the given context, axis, spacing,
// and alignment.
//
// Parameters:
// - ctx: The context for the layout.
// - axis: The axis of the layout.
// - spacing: The spacing between elements in the layout.
// - alignment: The alignment of elements in the layout.
//
// Returns:
// - A pointer to the newly created Flex layout.
func NewFlex(
	ctx types.Context,
	axis giolayout.Axis,
	spacing giolayout.Spacing,
	alignment giolayout.Alignment,
	rigidChildren ...types.UIElement,
) *Flex {
	chldArr := make([]*flexChild, 0, len(rigidChildren))
	for _, child := range rigidChildren {
		chldArr = append(chldArr, &flexChild{element: child})
	}

	return &Flex{
		Layout: NewLayout(ctx.Window()),
		flex: &giolayout.Flex{
			Axis:      axis,
			Spacing:   spacing,
			Alignment: alignment,
		},
		children: chldArr,
	}
}

func newFlexFromDefinition(
	ctx types.Context,
	data map[string]any,
) (definition.DefinitionType, error) {
	id, _ := definition.MapValueString[string](data, "id")
	axis, _ := definition.GioConstantFromMap[giolayout.Axis](data, "axis")
	spacing, _ := definition.GioConstantFromMap[giolayout.Spacing](data, "spacing")
	alignment, _ := definition.GioConstantFromMap[giolayout.Alignment](data, "alignment")
	f := NewFlex(ctx, axis, spacing, alignment)
	f.SetID(id)
	return f, nil
}

// AddChild adds a child element to the Flex layout.
//
// Parameters:
//   - child: The child element to be added.
//   - weight: (Optional) The weight of the child element. If provided, the
//     element behaves as a flexed element; otherwise, it behaves as a rigid
//     element.
func (f *Flex) AddChild(child types.UIElement, weight ...float32) bool {
	f.children = append(f.children, &flexChild{
		element: child,
		weight: func() *float32 {
			if len(weight) > 0 && weight[0] > 0 {
				return &weight[0]
			}
			return nil
		}(),
	})
	return true
}

func (f *Flex) HandleEvents(ctx types.Context) {
	for _, child := range f.children {
		child.element.HandleEvents(ctx)
	}
}

// Draw draws the Flex layout using the provided context.
//
// Parameters:
// - ctx: The context for the layout.
//
// Returns:
// - The dimensions of the drawn layout.
func (f *Flex) Draw(gtx giolayout.Context) giolayout.Dimensions {
	children := make([]giolayout.FlexChild, 0, len(f.children))
	for _, child := range f.children {
		children = append(children, child.draw())
	}
	return f.flex.Layout(gtx, children...)
}

type flexChild struct {
	element types.UIElement
	weight  *float32
}

// draw returns a giolayout.FlexChild based on the weight of the flexChild.
//
// If the weight is nil, it returns a giolayout.Rigid with the result of
// f.element.Draw.
// Otherwise, it returns a giolayout.Flexed with the weight and the result of
// f.element.Draw.
func (f *flexChild) draw() giolayout.FlexChild {
	if f.weight == nil {
		return giolayout.Rigid(f.element.Draw)
	}
	return giolayout.Flexed(*f.weight, f.element.Draw)
}
