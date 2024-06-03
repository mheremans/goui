// SPDX-License-Identifier: MIT

package types

import "gioui.org/widget/material"

type Window interface {
	Theme() *material.Theme
	Invalidate()
}
