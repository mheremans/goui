// SPDX-License-Identifier: MIT

package widget

import "gioui.org/text"

func gioTextWrapPolicyFromString(policy string) text.WrapPolicy {
	switch policy {
	case "WrapWords":
		return text.WrapWords
	case "WrapGraphemes":
		return text.WrapGraphemes
	default:
		return text.WrapGraphemes
	}
}
