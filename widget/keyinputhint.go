// SPDX-License-Identifier: MIT

package widget

import "gioui.org/io/key"

func gioInputHintFromString(hint string) key.InputHint {
	switch hint {
	case "HintText":
		return key.HintText
	case "HintNumeric":
		return key.HintNumeric
	case "HintInteger":
		return key.HintNumeric
	case "HintEmail":
		return key.HintEmail
	case "HintURL":
		return key.HintURL
	case "HintTelephone":
		return key.HintTelephone
	case "HintPassword":
		return key.HintPassword
	default:
		return key.HintAny
	}
}
