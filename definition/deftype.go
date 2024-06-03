// SPDX-License-Identifier: MIT

package definition

import "github.com/mheremans/goui/types"

type DefinitionType interface {
	types.UIElement
	AddChild(types.UIElement, ...float32) bool
}
