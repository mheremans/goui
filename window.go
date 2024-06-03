// SPDX-License-Identifier: MIT

package goui

import (
	"errors"
	"image/color"
	"log"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/mheremans/goui/types"

	_ "github.com/mheremans/goui/layout"
	_ "github.com/mheremans/goui/widget"
)

type CloseChan chan struct{}

type CloseHandler func()

type Window struct {
	initState *windowInitState // Temporary settings cache used to initialize the window
	closeChan CloseChan        // Channel that will notify when the window is closed

	w     *app.Window // The gio window
	theme *material.Theme
	op    op.Ops

	newView         types.View // The new view to render (replaces the old view)
	view            types.View // The view to render
	viewInitialized bool

	OnClose CloseHandler // OnClose callback
}

// NewWindow creates a new Window with the given title.
//
// It returns a pointer to the newly created Window.
func NewWindow(title string) *Window {
	wnd := &Window{
		closeChan: make(CloseChan),
		w:         new(app.Window),
		theme:     material.NewTheme(),
	}
	wnd.initState = &windowInitState{}
	wnd.initState.title = &title

	return wnd
}

// SetTitle sets the title of the Window.
//
// It takes a string parameter `title` representing the new title of the Window.
// If the Window is in the initialization state, the title is stored in the
// `initState` field.
// Otherwise the title is directly applied to the window.
func (wnd *Window) SetTitle(title string) {
	if wnd.initState != nil {
		wnd.initState.title = &title
		return
	}
	wnd.w.Option(app.Title(title))
}

// SetSize sets the size of the Window.
//
// It takes a types.WindowSize parameter `size` representing the new size of the
// Window.
// If the Window is in the initialization state, the size is stored in the
// `initState` field.
// Otherwise, the size is directly applied to the window.
func (wnd *Window) SetSize(size types.WindowSize) {
	if wnd.initState != nil {
		wnd.initState.size = &size
		return
	}
	wnd.w.Option(app.Size(size.Width, size.Height))
}

// SetMaxSize sets the maximum size of the Window.
//
// It takes a types.WindowSize parameter `size` representing the new maximum
// size of the Window.
// If the Window is in the initialization state, the maximum size is stored in
// the `initState` field.
// Otherwise, the maximum size is directly applied to the window.
func (wnd *Window) SetMaxSize(size types.WindowSize) {
	if wnd.initState != nil {
		wnd.initState.maxSize = &size
		return
	}
	wnd.w.Option(app.MaxSize(size.Width, size.Height))
}

// SetMinSize sets the minimum size of the Window.
//
// It takes a types.WindowSize parameter `size` representing the new minimum
// size of the Window.
// If the Window is in the initialization state, the minimum size is stored in
// the `initState` field.
// Otherwise, the minimum size is directly applied to the window.
func (wnd *Window) SetMinSize(size types.WindowSize) {
	if wnd.initState != nil {
		wnd.initState.minSize = &size
		return
	}
	wnd.w.Option(app.MinSize(size.Width, size.Height))
}

// SetStatusColor sets the status color of the windows Andoid status bar.
//
// It takes a color.NRGBA parameter `color` representing the new status color
// bar color  of the Android Window.
// If the Window is in the initialization state, the status color is stored in
// the `initState` field.
// Otherwise, the status color is directly applied to the window.
func (wnd *Window) SetStatusColor(color color.NRGBA) {
	if wnd.initState != nil {
		wnd.initState.statusColor = &color
		return
	}
	wnd.w.Option(app.StatusColor(color))
}

// SetNavigationColor sets the navigation color of the windows Android
// navigation bar or the address bar in browsers.
//
// It takes a color.NRGBA parameter `color` representing the new navigation
// color of the Window.
// If the Window is in the initialization state, the navigation color is stored
// in the `initState` field.
// Otherwise, the navigation color is directly applied to the window.
func (wnd *Window) SetNavigationColor(color color.NRGBA) {
	if wnd.initState != nil {
		wnd.initState.navigationColor = &color
		return
	}
	wnd.w.Option(app.NavigationColor(color))
}

// SetFullscreen sets the fullscreen mode for the window.
//
// fullscreen: a boolean indicating whether the window should be in fullscreen
// mode or not.
func (wnd *Window) SetFullscreen(fullscreen bool) {
	if wnd.initState != nil {
		wnd.initState.fullscreen = &fullscreen
		return
	}
	if fullscreen {
		wnd.w.Option(app.Fullscreen.Option())
	} else {
		wnd.w.Option(app.Windowed.Option())
	}
}

// SetScreen sets the view for the Window.
//
// This method swaps out the original view for the new view.
func (wnd *Window) SetView(view types.View) {
	wnd.newView = view
	wnd.viewInitialized = false
}

// Theme returns the material theme associated with the Window.
func (wnd Window) Theme() *material.Theme {
	return wnd.theme
}

func (wnd Window) Invalidate() {
	wnd.w.Invalidate()
}

// Show shows the window with the given view.
//
// It takes a View parameter and returns a CloseChan and an error.
// The CloseChan is a channel that will receive a value when the window is
// closed.
// The error is non-nil if the window is already running.
func (wnd *Window) Show(view types.View) (CloseChan, error) {
	if wnd.initState == nil {
		return nil, errors.New("Window already running")
	}

	if view != nil {
		wnd.SetView(view)
	}

	go func() {
		defer func() {
			if wnd.OnClose != nil {
				wnd.OnClose()
			}
			wnd.closeChan <- struct{}{}
			close(wnd.closeChan)
		}()

		wnd.initState.initWindow(wnd)
		wnd.initState = nil

		err := wnd.eventLoop()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return wnd.closeChan, nil
}

// eventLoop is a method of the Window struct that continuously listens for
// events and performs actions based on the type of event received. It runs in
// a loop until a DestroyEvent is received, at which point it returns the error
// associated with the event.
func (wnd *Window) eventLoop() error {
	for {
		switch e := wnd.w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			ctx := NewContext(wnd, e)

			// Check if the view has changed
			if err := wnd.swapView(ctx); err != nil {
				return err
			}

			// Initialize the view if needed
			if err := wnd.initView(ctx); err != nil {
				return err
			}

			if wnd.view != nil {
				wnd.view.HandleEvents(ctx)
				wnd.view.DrawView(ctx)
			}
			e.Frame(ctx.Gtx().Ops)
		}
	}
}

func (wnd *Window) swapView(ctx types.Context) (err error) {
	if wnd.newView == nil {
		return
	}

	if wnd.newView == wnd.view {
		return
	}

	if wnd.view != nil {
		if err = wnd.view.Destroy(ctx); err != nil {
			return
		}
	}
	wnd.view = wnd.newView
	wnd.newView = nil
	wnd.viewInitialized = false
	return
}

func (wnd *Window) initView(ctx types.Context) (err error) {
	if wnd.view == nil {
		return
	}

	if wnd.viewInitialized {
		return
	}

	if err = wnd.view.Initialize(ctx); err != nil {
		return
	}
	wnd.viewInitialized = true
	return
}

type windowInitState struct {
	title           *string
	size            *types.WindowSize
	maxSize         *types.WindowSize
	minSize         *types.WindowSize
	statusColor     *color.NRGBA
	navigationColor *color.NRGBA
	fullscreen      *bool
}

// initWindow initializes the window with the provided options.
func (s windowInitState) initWindow(wnd *Window) {
	var options []app.Option

	if s.title != nil {
		options = append(options, app.Title(*s.title))
	}
	if s.size != nil {
		options = append(options, app.Size(s.size.Width, s.size.Height))
	}
	if s.maxSize != nil {
		options = append(options, app.MaxSize(s.maxSize.Width, s.maxSize.Height))
	}
	if s.minSize != nil {
		options = append(options, app.MinSize(s.minSize.Width, s.minSize.Height))
	}
	if s.statusColor != nil {
		options = append(options, app.StatusColor(*s.statusColor))
	}
	if s.navigationColor != nil {
		options = append(options, app.NavigationColor(*s.navigationColor))
	}
	if s.fullscreen != nil {
		if *s.fullscreen {
			options = append(options, app.Fullscreen.Option())
		} else {
			options = append(options, app.Windowed.Option())
		}
	}

	wnd.w.Option(options...)
}
