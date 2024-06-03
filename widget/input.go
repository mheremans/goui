// SPDX-License-Identifier: MIT

package widget

import (
	"image/color"
	"strconv"
	"strings"

	"gioui.org/io/key"
	giolayout "gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mheremans/goui/definition"
	"github.com/mheremans/goui/types"
)

func init() {
	definition.RegisterUIElement((*Input)(nil), newInputFromDefinition)
}

type InputType uint

const (
	TypeAny InputType = iota
	TypeText
	TypeNumeric
	TypeInteger
	TypeEmail
	TypeURL
	TypeTelephone
	TypePassword
)

func (t InputType) InputHint() key.InputHint {
	switch t {
	case TypeText:
		return key.HintText
	case TypeNumeric:
		return key.HintNumeric
	case TypeInteger:
		return key.HintNumeric
	case TypeEmail:
		return key.HintEmail
	case TypeURL:
		return key.HintURL
	case TypeTelephone:
		return key.HintTelephone
	case TypePassword:
		return key.HintPassword
	default:
		return key.HintAny
	}
}

func (t InputType) String() string {
	return []string{
		"Any",
		"Text",
		"Numeric",
		"Integer",
		"Email",
		"URL",
		"Telephone",
		"Password"}[t]
}

func (t InputType) DefaultMaskRune() rune {
	if t == TypePassword {
		return '*'
	}
	return 0
}

func (t InputType) DefaultFilter() string {
	switch t {
	case TypeNumeric:
		return "0123456789.-"
	case TypeInteger:
		return "0123456789-"
	case TypeEmail:
		return "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+-_~!#$%&'./=^`{}|@"
	case TypeURL:
		return "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+-*_~!?#$%&'.:,;/=()[]@"
	case TypeTelephone:
		return "0123456789-./ ()+"
	default:
		return ""
	}
}

func (t InputType) AdjustedFilter(txt string) (f string) {
	f = t.DefaultFilter()
	switch t {
	case TypeNumeric:
		if len(txt) == 0 {
			return
		}
		f = strings.Replace(f, "-", "", 1)
		if strings.Contains(txt, ".") {
			f = strings.Replace(f, ".", "", 1)
		}
	case TypeInteger:
		if len(txt) == 0 {
			return
		}
		f = strings.Replace(f, "-", "", 1)
	case TypeEmail:
		if len(txt) == 0 {
			return
		}
		if strings.Contains(txt, "@") {
			f = strings.Replace(f, "@", "", 1)
		}
	case TypeURL:
		if len(txt) == 0 {
			return
		}
		if strings.Contains(txt, "?") {
			f = strings.Replace(f, "?", "", 1)
		}
	case TypeTelephone:
		if len(txt) == 0 {
			return
		}
		f = strings.Replace(f, "+", "", 1)
		if strings.Contains(txt, "(") {
			f = strings.Replace(f, "(", "", 1)
		}
		if strings.Contains(txt, ")") {
			f = strings.Replace(f, ")", "", 1)
		}
	}
	return
}

func InputTypeFromString(s string) InputType {
	switch s {
	case "Text":
		return TypeText
	case "Numeric":
		return TypeNumeric
	case "Integer":
		return TypeInteger
	case "Email":
		return TypeEmail
	case "URL":
		return TypeURL
	case "Telephone":
		return TypeTelephone
	case "Password":
		return TypePassword
	default:
		return TypeAny
	}
}

// Input is a text input
//
// Yaml definition:
//
//	type: widget.Input
//	id: <string>				# id of the element (used to get a reference to it in code)
//	alignment: <string>			# input alignment ("Start", "End", "Middle")
//	hint: <string>				# input hint
//	singleLine: <bool>			# single line input
//	readOnly: <bool>			# read only input
//	submit: <bool>				# submit input
//	mask: <string>				# input mask
//	maxLen: <int>				# input max length
//	inputType: <string>			# input hint ("Any", "Text", "Numeric", "Integer", "Email", "URL", "Telephone", "Password")
//	wrapPolicy: <string>		# input wrap policy ("WrapHeuristically", "WrapWords", "WrapGraphemes")
//	filterCallback: <string>	# input filter callback function to adjust the filter based on the current input
//	binding: <string>			# binding reference (will be requested throught the view)
type Input struct {
	*Widget

	editor *material.EditorStyle
	input  widget.Editor

	binding *types.Binding[string]

	inputType     InputType
	inputFilterFn InputFilterFn
}

func NewInput(ctx types.Context, hint string, id ...string) *Input {
	i := new(Input)
	i.Widget = NewWidget(ctx.Window(), id...)
	editor := material.Editor(ctx.Window().Theme(), &i.input, hint)
	i.editor = &editor

	return i
}

func newInputFromDefinition(
	ctx types.Context,
	data map[string]any,
) (definition.DefinitionType, error) {
	id, _ := definition.MapValueString[string](data, "id")
	alignment, _ := definition.GioConstantFromMap[text.Alignment](data, "alignment")
	hint, _ := definition.MapValueString[string](data, "hint")
	multiLine, _ := definition.MapValueBool[bool](data, "multiLine")
	readOnly, _ := definition.MapValueBool[bool](data, "readOnly")
	submit, _ := definition.MapValueBool[bool](data, "submit")
	mask, _ := definition.MapValueString[string](data, "mask")
	maxLen, _ := definition.MapValueInt[int](data, "maxLen")
	inputType, _ := definition.MapValueString[string](data, "inputType")
	wrapPolicy, _ := definition.MapValueString[string](data, "wrapPolicy")
	inputFilterFn, _ := definition.FunctionFromMap[InputFilterFn](ctx, data, "filterCallback")

	i := NewInput(ctx, hint, id)
	i.input.Alignment = alignment
	i.input.SingleLine = !multiLine
	i.input.ReadOnly = readOnly
	i.input.Submit = submit
	i.input.MaxLen = maxLen
	i.inputType = InputTypeFromString(inputType)
	i.input.InputHint = i.inputType.InputHint()
	i.inputFilterFn = inputFilterFn

	if i.inputFilterFn != nil {
		i.input.Filter = i.inputFilterFn(ctx, i, "")
	} else {
		i.input.Filter = i.inputType.DefaultFilter()
	}

	i.input.WrapPolicy = gioTextWrapPolicyFromString(wrapPolicy)

	if mask != "" {
		i.input.Mask = rune(mask[0])
	} else {
		i.input.Mask = i.inputType.DefaultMaskRune()
	}

	if binding, ok := definition.BindingFromMap[*types.Binding[string]](
		ctx, data, "binding",
	); ok {
		i.Bind(binding)
	}
	return i, nil
}

func (i *Input) Bind(binding *types.Binding[string]) {
	if i.binding != nil {
		i.binding.Unwatch(i)
		i.binding = nil
	}

	if binding == nil {
		return
	}

	i.binding = binding
	i.binding.Watch(i)
	i.SetText(i.binding.Get())
}

func (i Input) SingleLine() bool {
	return i.input.SingleLine
}

func (i Input) ReadOnly() bool {
	return i.input.ReadOnly
}

func (i Input) Submit() bool {
	return i.input.Submit
}

func (i Input) Mask() rune {
	return i.input.Mask
}

func (i Input) MaxLen() int {
	return i.input.MaxLen
}

func (i Input) Filter() string {
	return i.input.Filter
}

func (i Input) InputHint() key.InputHint {
	return i.input.InputHint
}

func (i Input) WrapPolicy() text.WrapPolicy {
	return i.input.WrapPolicy
}

func (i Input) Text() string {
	return i.input.Text()
}

func (i Input) AsInteger() (int, error) {
	return strconv.Atoi(i.Text())
}

func (i Input) AsFloat() (float64, error) {
	return strconv.ParseFloat(i.Text(), 64)
}

func (i *Input) SetSignleLine(singleLine bool) {
	i.input.SingleLine = singleLine
}

func (i *Input) SetReadOnly(readOnly bool) {
	i.input.ReadOnly = readOnly
}

func (i *Input) SetSubmit(submit bool) {
	i.input.Submit = submit
}

func (i *Input) SetMask(mask rune) {
	i.input.Mask = mask
}

func (i *Input) SetMaxLen(maxLen int) {
	i.input.MaxLen = maxLen
}

func (i *Input) SetFilter(filter string) {
	i.input.Filter = filter
}

func (i *Input) SetInputHint(inputHint key.InputHint) {
	i.input.InputHint = inputHint
}

func (i *Input) SetWrapPolicy(wrapPolicy text.WrapPolicy) {
	i.input.WrapPolicy = wrapPolicy
}

func (i *Input) SetText(text string) {
	i.input.SetText(text)
	i.Wnd().Invalidate()
}

func (i *Input) HandleEvents(ctx types.Context) {
	txt := i.input.Text()
	if i.inputFilterFn != nil {
		i.input.Filter = i.inputFilterFn(ctx, i, txt)
	} else {
		i.input.Filter = i.inputType.AdjustedFilter(txt)
	}
	if i.binding != nil {
		i.binding.Set(txt)
	}
}

func (i *Input) Draw(gtx giolayout.Context) giolayout.Dimensions {
	border := widget.Border{
		Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
		CornerRadius: unit.Dp(3),
		Width:        unit.Dp(2),
	}
	inset := giolayout.Inset{
		Top:    unit.Dp(3),
		Bottom: unit.Dp(3),
		Left:   unit.Dp(3),
		Right:  unit.Dp(3),
	}

	return border.Layout(gtx, func(gtx giolayout.Context) giolayout.Dimensions {
		return inset.Layout(gtx, i.editor.Layout)
	})
}

func (i *Input) BindingChanged(binding types.Bindable) {
	if bnd, ok := binding.(*types.Binding[string]); ok {
		if i.input.Text() != bnd.Get() {
			i.input.SetText(bnd.Get())
			i.Wnd().Invalidate()
		}
	}
}

func inputCharNumericValidator(text string, _ types.UIElement) error {
	_, err := strconv.ParseFloat(text, 64)
	return err
}
