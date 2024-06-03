// SPDX-License-Identifier: MIT

package goui

import "github.com/mheremans/goui/types"

type ViewModel struct {
	bindings map[string]types.Bindable
}

func (vm *ViewModel) Initialize() (err error) {
	return
}

func (vm *ViewModel) Destroy() (err error) {
	return
}

func (vm *ViewModel) GetBinding(name string) types.Bindable {
	if vm.bindings == nil {
		return nil
	}
	if b, ok := vm.bindings[name]; ok {
		return b
	}
	return nil
}

func (vm *ViewModel) RegisterBindings(bindings ...types.Bindable) {
	for _, b := range bindings {
		vm.RegisterBinding(b)
	}
}

func (vm *ViewModel) RegisterBinding(binding types.Bindable) {
	if vm.bindings == nil {
		vm.bindings = make(map[string]types.Bindable)
	}
	vm.bindings[binding.Name()] = binding
}
