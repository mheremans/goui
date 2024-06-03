// SPDX-License-Identifier: MIT

package goui

import (
	"embed"
	"fmt"

	giolayout "gioui.org/layout"

	"github.com/mheremans/goui/definition"
	"github.com/mheremans/goui/types"
	"github.com/mheremans/goui/widget"
)

type ViewScreen struct {
	fs         embed.FS
	screenName string
}

func NewViewScreen(
	fs embed.FS,
	screenName string,
) *ViewScreen {
	return &ViewScreen{
		fs:         fs,
		screenName: screenName,
	}
}

type View struct {
	*widget.Widget

	impl       types.View
	viewModel  types.ViewModel
	viewScreen *ViewScreen

	def  *definition.Definition
	root types.UIElement

	exportedFns map[string]any
}

func ConfigureView(
	view types.View,
	viewModel types.ViewModel,
	viewScreen *ViewScreen,
) *View {
	v := new(View)
	if viewScreen != nil {
		v.Widget = widget.NewWidget(nil, viewScreen.screenName)
	} else {
		v.Widget = widget.NewWidget(nil)
	}
	v.impl = view
	v.viewModel = viewModel
	v.viewScreen = viewScreen
	v.exportedFns = make(map[string]any)
	return v
}

func (v *View) ExportFunction(name string, fn any) {
	v.exportedFns[name] = fn
}

func (v *View) SetViewRoot(root types.UIElement) {
	v.root = root
}

func (v *View) ViewModel() types.ViewModel {
	return v.viewModel
}

func (v *View) Initialize(ctx types.Context) (err error) {
	ctx.SetView(v.impl)
	v.Widget.SetWnd(ctx.Window())

	v.viewModel.Initialize()

	if v.viewScreen != nil {
		v.def, err = definition.New(ctx, v.viewScreen.fs, v.viewScreen.screenName)
		if err != nil {
			err = fmt.Errorf("failed to create definition: %w", err)
			return err
		}
		v.root = v.def.Root()
	}

	return
}

func (v *View) Destroy(ctx types.Context) (err error) {
	err = v.viewModel.Destroy()
	return
}

func (v *View) HandleEvents(ctx types.Context) {
	v.root.HandleEvents(ctx)
}

func (v *View) Draw(gtx giolayout.Context) giolayout.Dimensions {
	return v.root.Draw(gtx)
}

func (v *View) DrawView(ctx types.Context) giolayout.Dimensions {
	return v.Draw(ctx.Gtx())
}

func (v *View) FindFunction(name string) any {
	if fn, ok := v.exportedFns[name]; ok {
		return fn
	}
	return nil
}

func (v *View) FindBinding(name string) types.Bindable {
	return v.viewModel.GetBinding(name)
}

func GetElementById[T any](v *View, id string) (elem T, ok bool) {
	return definition.ElementById[T](v.def, id)
}
