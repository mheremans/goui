// SPDX-License-Identifier: MIT

package layout

import (
	"github.com/google/uuid"
	"github.com/mheremans/goui/types"
)

// Layout implementation of types.UIElement
type Layout struct {
	wnd types.Window
	id  string
}

func NewLayout(wnd types.Window, id ...string) *Layout {
	return &Layout{
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

func (l Layout) Wnd() types.Window {
	return l.wnd
}

func (l Layout) ID() string {
	return l.id
}

func (l *Layout) SetWnd(wnd types.Window) {
	l.wnd = wnd
}

func (l *Layout) SetID(id string) {
	if id == "" {
		id = uuid.NewString()
	}
	l.id = id
}
