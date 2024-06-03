// SPDX-License-Identifier: MIT

package goui

import (
	"embed"

	"gioui.org/app"
	"gioui.org/layout"
	"github.com/mheremans/goui/types"
)

type Context struct {
	window   *Window
	view     types.View
	gtx      layout.Context
	fontsDir *embed.FS
}

// NewContext creates a new Context with the given window and frame event.
//
// Parameters:
// - window: a pointer to the Window struct.
// - e: the app.FrameEvent.
//
// Returns:
// - a pointer to the newly created Context.
func NewContext(window *Window, e app.FrameEvent) *Context {
	return &Context{
		window: window,
		gtx:    app.NewContext(&window.op, e),
	}
}

// Window returns the Window associated with the Context.
func (c Context) Window() types.Window {
	return c.window
}

func (c Context) View() types.View {
	return c.view
}

// Gtx returns the layout.Context associated with the Context.
func (c Context) Gtx() layout.Context {
	return c.gtx
}

func (c *Context) FontsDir() *embed.FS {
	return c.fontsDir
}

func (c *Context) SetFontsDir(fs *embed.FS) {
	c.fontsDir = fs
}

func (c *Context) SetView(v types.View) {
	c.view = v
}
