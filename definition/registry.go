// SPDX-License-Identifier: MIT

package definition

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mheremans/goui/types"
)

type ConstructorFn func(types.Context, map[string]any) (DefinitionType, error)

var uiElementRegistry map[string]ConstructorFn

func init() {
	uiElementRegistry = make(map[string]ConstructorFn)
}

func RegisterUIElement(e types.UIElement, constructor ConstructorFn) {
	t := reflect.TypeOf(e).Elem()
	tn := strings.TrimPrefix(t.PkgPath(), "github.com/mheremans/goui/")
	tn = tn + "." + t.Name()
	uiElementRegistry[tn] = constructor
}

func Instantiate(
	ctx types.Context,
	name string,
	data map[string]any,
) (
	res DefinitionType,
	err error,
) {
	constructor, ok := uiElementRegistry[name]
	if !ok {
		err = fmt.Errorf("no such element: %s", name)
		return
	}

	res, err = constructor(ctx, data)
	if err != nil {
		return
	}
	return
}
