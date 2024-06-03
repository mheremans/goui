// SPDX-License-Identifier: MIT

package types

import (
	giolayout "gioui.org/layout"
)

type UIElement interface {
	HandleEvents(Context)
	Draw(giolayout.Context) giolayout.Dimensions

	ID() string
	SetID(string)

	Wnd() Window
}
