// SPDX-License-Identifier: MIT

package definition

import (
	"embed"
	"fmt"
	"io"
	"io/fs"

	"github.com/mheremans/goui/types"
	"gopkg.in/yaml.v3"
)

type Definition struct {
	root  DefinitionType
	index map[string]DefinitionType
}

func (d Definition) Root() types.UIElement {
	return d.root
}

func (d Definition) ElementById(id string) (elem types.UIElement, ok bool) {
	elem, ok = d.index[id]
	return
}

func ElementById[T any](def *Definition, id string) (elem T, ok bool) {
	tmp, ok := def.index[id]
	if !ok {
		return
	}
	elem, ok = tmp.(T)
	return
}

func New(
	ctx types.Context,
	filesystem embed.FS,
	name string,
) (def *Definition, err error) {
	var bytes []byte
	var fh fs.File

	if fh, err = filesystem.Open(name); err != nil {
		err = fmt.Errorf("no definition found: %w", err)
		return
	}
	defer fh.Close()

	if bytes, err = io.ReadAll(fh); err != nil {
		err = fmt.Errorf("unable to read definition: %w", err)
		return
	}

	defMap := make(map[string]interface{})
	if err = yaml.Unmarshal(bytes, defMap); err != nil {
		err = fmt.Errorf("definition has syntax error: %w", err)
		return
	}

	def, err = createLayout(ctx, defMap)
	if err != nil {
		err = fmt.Errorf("failed to create definition: %w", err)
		return
	}
	return
}

func createLayout(
	ctx types.Context,
	defMap map[string]any,
) (
	def *Definition,
	err error,
) {
	def = &Definition{
		index: make(map[string]DefinitionType),
	}

	root, _, _, childDefs, err := createElement(ctx, defMap)
	if err != nil {
		err = fmt.Errorf("failed to create layout: %w", err)
		return
	}

	for _, childDef := range childDefs {
		var child DefinitionType
		var childId string
		var childWeight *float32
		var ok bool

		child, childId, childWeight, err = createChild(ctx, def, childDef)
		if err != nil {
			err = fmt.Errorf("failed to create layout: %w", err)
			return
		}

		if childWeight != nil {
			ok = root.AddChild(child, *childWeight)
		} else {
			ok = root.AddChild(child)
		}
		if !ok {
			err = fmt.Errorf("failed to add child to layout")
			return
		}
		if childId != "" {
			def.index[childId] = child
		}
	}
	def.root = root
	return
}

func createChild(
	ctx types.Context,
	def *Definition,
	defMap map[string]any,
) (
	elem DefinitionType,
	id string,
	weight *float32,
	err error,
) {
	elem, id, weight, childDefs, err := createElement(ctx, defMap)
	if err != nil {
		err = fmt.Errorf("failed to create child: %w", err)
		return
	}

	for _, childDef := range childDefs {
		var child DefinitionType
		var childId string
		var childWeight *float32
		var ok bool

		child, childId, childWeight, err = createChild(ctx, def, childDef)
		if err != nil {
			err = fmt.Errorf("failed to create child: %w", err)
			return
		}

		if childWeight != nil {
			ok = elem.AddChild(child, *childWeight)
		} else {
			ok = elem.AddChild(child)
		}
		if !ok {
			err = fmt.Errorf("failed to add child to layout")
			return
		}
		if childId != "" {
			def.index[childId] = child
		}
	}
	return
}

func createElement(
	ctx types.Context,
	defMap map[string]interface{},
) (
	elem DefinitionType,
	id string,
	weight *float32,
	childDefinitions []map[string]interface{},
	err error,
) {
	typeName, ok := MapValueString[string](defMap, "type")
	if !ok {
		err = fmt.Errorf("element has no type")
		return
	}

	id, _ = MapValueString[string](defMap, "id")
	if w, ok := MapValueFloat[float32](defMap, "weight"); ok {
		weight = &w
	}

	elem, err = Instantiate(ctx, typeName, defMap)
	if err != nil {
		return
	}

	if children, ok := defMap["children"]; ok {
		childDefinitions = make([]map[string]any, 0)
		for _, chld := range children.([]any) {
			childDefinitions = append(
				childDefinitions, chld.(map[string]any))
		}
	}
	if child, ok := defMap["child"]; ok {
		childDefinitions = []map[string]any{child.(map[string]any)}
	}
	return
}
