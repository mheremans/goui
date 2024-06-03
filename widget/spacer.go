// SPDX-License-Identifier: MIT

package widget

import (
	giolayout "gioui.org/layout"
	"gioui.org/unit"
	"github.com/mheremans/goui/definition"
	"github.com/mheremans/goui/types"
)

func init() {
	definition.RegisterUIElement((*Spacer)(nil), newSpacerFromDefinition)
}

// Spacer is a widget that can be used to add spacing between other widgets.
//
// Yaml definition:
//
//	type: widget.Spacer
//	id: <string>		# id of the element (used to get a reference to it in code)
//	width: <number>		# width of the spacer (in Dp units)
//	height: <number>	# height of the spacer (in Dp units)
type Spacer struct {
	*Widget

	spacer *giolayout.Spacer
}

// NewSpacer creates a new Spacer widget with the given size.
//
// Parameters:
// - ctx: the context in which the widget is being created.
// - size: the size of the spacer widget.
//
// Returns:
// - *Spacer: a pointer to the newly created Spacer widget.
func NewSpacer(ctx types.Context, size types.SpacerSize, id ...string) *Spacer {
	return &Spacer{
		Widget: NewWidget(ctx.Window(), id...),
		spacer: &giolayout.Spacer{Width: size.Width, Height: size.Height},
	}
}

func newSpacerFromDefinition(
	ctx types.Context,
	data map[string]any,
) (definition.DefinitionType, error) {
	id, _ := definition.MapValueString[string](data, "id")
	width, _ := definition.MapValueFloat[unit.Dp](data, "width")
	height, _ := definition.MapValueFloat[unit.Dp](data, "height")
	return NewSpacer(ctx, types.NewSpacerSize(width, height), id), nil
}

// HandleEvents handles the events for the Spacer widget.
//
// Parameters:
// - ctx: the context in which the input is being handled.
func (s *Spacer) HandleEvents(ctx types.Context) {

}

// Draw renders the Spacer widget on the given giolayout.Context.
//
// Parameters:
// - gtx: the giolayout.Context on which to render the Spacer.
//
// Returns:
// - giolayout.Dimensions: the dimensions of the rendered Spacer.
func (s *Spacer) Draw(gtx giolayout.Context) giolayout.Dimensions {
	return s.spacer.Layout(gtx)
}

func (s *Spacer) AddChild(_ types.UIElement, _ ...float32) bool {
	return false
}
