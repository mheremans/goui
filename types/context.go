// SPDX-License-Identifier: MIT

package types

import (
	"embed"

	"gioui.org/layout"
)

type Context interface {
	Window() Window
	View() View
	Gtx() layout.Context

	FontsDir() *embed.FS
	SetFontsDir(*embed.FS)

	SetView(View)
}
