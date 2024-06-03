// SPDX-License-Identifier: MIT

package widget

import (
	"github.com/google/uuid"
	"github.com/mheremans/goui/types"
)

type Widget struct {
	wnd types.Window
	id  string
}

func NewWidget(wnd types.Window, id ...string) *Widget {
	return &Widget{
		wnd: wnd,
		id: func() string {
			if len(id) == 0 {
				return uuid.NewString()
			}
			if id[0] == "" {
				return uuid.NewString()
			}
			return id[0]
		}(),
	}
}

func (w Widget) Wnd() types.Window {
	return w.wnd
}

func (w Widget) ID() string {
	return w.id
}

func (w *Widget) SetWnd(wnd types.Window) {
	w.wnd = wnd
}

func (w *Widget) SetID(id string) {
	if id == "" {
		id = uuid.NewString()
	}
	w.id = id
}

func (w *Widget) AddChild(_ types.UIElement, _ ...float32) bool {
	return false
}
