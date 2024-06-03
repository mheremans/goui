// SPDX-License-Identifier: MIT

package types

type ViewModel interface {
	Initialize() error
	Destroy() error

	GetBinding(string) Bindable
}
