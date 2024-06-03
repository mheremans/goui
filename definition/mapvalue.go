// SPDX-License-Identifier: MIT

package definition

import (
	"fmt"
	"image/color"
	"strconv"

	giofont "gioui.org/font"
	giolayout "gioui.org/layout"
	"github.com/mheremans/goui/colors"
	"github.com/mheremans/goui/fonts"
	"github.com/mheremans/goui/types"
	"golang.org/x/exp/constraints"
)

type gioConst interface {
	fmt.Stringer
	constraints.Integer
}

func MapValue[T any](
	data map[string]any,
	key string,
) (
	res T,
	ok bool,
) {
	v, ok := data[key]
	if !ok {
		return
	}
	res, ok = v.(T)
	return
}

func MapValueInt[T constraints.Integer](
	data map[string]any,
	key string,
) (
	res T,
	ok bool,
) {
	v, ok := data[key]
	if !ok {
		return
	}
	res, ok = asNumber[T](v)
	return
}

func MapValueString[T ~string](
	data map[string]any,
	key string,
) (
	res T,
	ok bool,
) {
	v, ok := data[key]
	if !ok {
		return
	}
	tmp, ok := v.(string)
	if !ok {
		return
	}
	res = T(tmp)
	return
}

func MapValueFloat[T constraints.Float](
	data map[string]any,
	key string,
) (
	res T,
	ok bool,
) {
	v, ok := data[key]
	if !ok {
		return
	}
	res, ok = asNumber[T](v)
	return
}

func MapValueBool[T ~bool](
	data map[string]any,
	key string,
) (
	res T,
	ok bool,
) {
	v, ok := data[key]
	if !ok {
		return
	}
	tmp, ok := v.(bool)
	if !ok {
		return
	}
	res = T(tmp)
	return
}

func MapValueColor(
	data map[string]any,
	key string,
) (
	res color.NRGBA,
	ok bool,
) {
	var err error

	v, ok := data[key]
	if !ok {
		return
	}
	tmp, ok := v.(string)
	if !ok {
		return
	}
	res, err = colors.GetColor(tmp)
	if err != nil {
		ok = false
		return
	}
	return
}

func MapValueFont(
	ctx types.Context,
	data map[string]any,
	fontKey string,
	fontStyleKey string,
	fontWeightKey string,
) (
	res giofont.FontFace,
	ok bool,
) {
	var err error

	v, ok := data[fontKey]
	if !ok {
		return
	}
	font, ok := v.(string)
	if !ok {
		return
	}

	fontStyle, _ := GioConstantFromMap[giofont.Style](data, fontStyleKey)
	fontWeight, _ := GioConstantFromMap[giofont.Weight](data, fontWeightKey)

	res, err = fonts.GetOTFFont(ctx, font, fontStyle, fontWeight)
	if err != nil {
		ok = false
		return
	}
	return
}

func GioConstantFromMap[T gioConst](
	data map[string]any,
	key string,
) (
	res T,
	ok bool,
) {
	sv, ok := MapValueString[string](data, key)
	if !ok {
		return
	}

	res, ok = findGioConstant[T](sv)
	return
}

func FunctionFromMap[T any](
	ctx types.Context,
	data map[string]any,
	key string,
) (
	res T,
	ok bool,
) {
	sv, ok := MapValueString[string](data, key)
	if !ok {
		return
	}
	fn := ctx.View().FindFunction(sv)
	if fn == nil {
		ok = false
		return
	}
	res, ok = fn.(T)
	return
}

func BindingFromMap[T types.Bindable](
	ctx types.Context,
	data map[string]any,
	key string,
) (
	res T,
	ok bool,
) {
	sv, ok := MapValueString[string](data, key)
	if !ok {
		return
	}
	bnd := ctx.View().FindBinding(sv)
	if bnd == nil {
		ok = false
		return
	}
	res, ok = bnd.(T)
	return
}

func findGioConstant[T gioConst](name string) (res T, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			// Work around gio bug
			if name == "SpaceBetween" {
				res = T(giolayout.SpaceBetween)
				ok = true
				return
			}

			ok = false
		}
	}()

	var t T
	for {
		if t.String() == name {
			ok = true
			res = t
			return
		}
		t++
	}

}

func asNumber[T constraints.Integer | constraints.Float](v any) (T, bool) {
	switch v := v.(type) {
	case int:
		return T(v), true
	case float64:
		return T(v), true
	case string:
		if f, err := strconv.ParseFloat(v, 64); err != nil {
			return T(f), true
		}
		if i, err := strconv.ParseInt(v, 10, 64); err != nil {
			return T(i), true
		}
	}
	return T(0), false
}
