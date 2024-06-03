// SPDX-License-Identifier: MIT

package types

import giolayout "gioui.org/layout"

type View interface {
	UIElement

	Initialize(Context) error
	Destroy(Context) error

	DrawView(Context) giolayout.Dimensions

	FindFunction(string) any
	FindBinding(string) Bindable
}
