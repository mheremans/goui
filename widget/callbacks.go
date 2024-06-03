// SPDX-License-Identifier: MIT

package widget

import (
	"image"

	giolayout "gioui.org/layout"
	"github.com/mheremans/goui/types"
)

// Hooks

// GraphicFn is a function that draws on the screen
type GraphicFn = func(giolayout.Context, types.UIElement) image.Point

// Handleres

// OnClickedFn is called when the button is clicked
type OnClickedFn = func(types.Context, types.UIElement)
type OnHoveredFn = func(types.Context, types.UIElement)
type OnHoverEnteredFn = func(types.Context, types.UIElement)
type OnHoverExitedFn = func(types.Context, types.UIElement)
type OnPressedFn = func(types.Context, types.UIElement)
type OnPressDownFn = func(types.Context, types.UIElement)
type OnPressUpFn = func(types.Context, types.UIElement)

// Validators

type InputFilterFn = func(types.Context, types.UIElement, string) string

// Item Renderers
type ListItemEventHandlerFn = func(types.Context, int, types.BindableList)
type ListItemRendererFn = func(giolayout.Context, int, types.BindableList) giolayout.Dimensions
