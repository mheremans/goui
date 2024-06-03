// SPDX-License-Identifier: MIT

package icons

import (
	"gioui.org/widget"
	materialicons "golang.org/x/exp/shiny/materialdesign/icons"
)

var icons map[string]int

func init() {
	_ = materialicons.AVAVTimer
	icons = make(map[string]int)
	for i, icon := range list {
		icons[icon.name] = i
	}
}

func Icon(name string) *widget.Icon {
	data := findIcon(name)
	if data == nil {
		return nil
	}
	w, err := widget.NewIcon(data)
	if err != nil {
		return nil
	}
	return w
}

func findIcon(name string) []byte {
	i, ok := icons[name]
	if !ok {
		return nil
	}
	return list[i].data
}
