// SPDX-License-Identifier: MIT

package types

type BindableStructValue interface {
	Less(BindableStructValue) bool
}

type BindingWatcher interface {
	BindingChanged(Bindable)
}

type Bindable interface {
	Name() string
	Watch(BindingWatcher)
	Unwatch(BindingWatcher)
}

type BindableList interface {
	Bindable
	GetAt(int) (any, bool)
	Size() int
}

type binding struct {
	name     string
	watchers map[BindingWatcher]struct{}
}

func newBinding(name string) *binding {
	return &binding{
		name:     name,
		watchers: make(map[BindingWatcher]struct{}),
	}
}

func (b binding) Name() string {
	return b.name
}

func (b *binding) Watch(watcher BindingWatcher) {
	b.watchers[watcher] = struct{}{}
}

func (b *binding) Unwatch(watcher BindingWatcher) {
	delete(b.watchers, watcher)
}

func (b *binding) notify(binding Bindable) {
	for w := range b.watchers {
		w.BindingChanged(binding)
	}
}

type Binding[T comparable] struct {
	*binding
	value T
}

func NewBinding[T comparable](name string, value T) *Binding[T] {
	return &Binding[T]{
		binding: newBinding(name),
		value:   value,
	}
}

func (b Binding[T]) Get() T {
	return b.value
}

func (b *Binding[T]) Set(value T) {
	if value != b.value {
		b.value = value
		b.notify(b)
	}
}

type StructBinding[T BindableStructValue] struct {
	*binding
	value T
}

func NewStructBinding[T BindableStructValue](name string, value T) *StructBinding[T] {
	return &StructBinding[T]{
		binding: newBinding(name),
		value:   value,
	}
}

func (b StructBinding[T]) Get() T {
	return b.value
}

func (b *StructBinding[T]) Set(value T) {
	if value.Less(b.value) || b.value.Less(value) {
		b.value = value
		b.notify(b)
	}
}

type ListBinding[T comparable] struct {
	*binding
	list []T
}

func NewListBinding[T comparable](name string, values []T) *ListBinding[T] {
	b := &ListBinding[T]{
		binding: newBinding(name),
		list:    make([]T, 0, len(values)),
	}
	b.list = append(b.list, values...)
	return b
}

func (b ListBinding[T]) Get() []T {
	return b.list
}

func (b ListBinding[T]) GetAt(index int) (any, bool) {
	if index < 0 || index >= len(b.list) {
		return *new(T), false
	}
	return b.list[index], true
}

func (b ListBinding[T]) Size() int {
	return len(b.list)
}

func (b *ListBinding[T]) Set(values []T) {
	if b.equalValues(values) {
		return
	}

	b.list = make([]T, 0, len(values))
	b.list = append(b.list, values...)
	b.notify(b)
}

func (b ListBinding[T]) equalValues(other []T) bool {
	if len(b.list) != len(other) {
		return false
	}
	for i := range b.list {
		if b.list[i] != other[i] {
			return false
		}
	}
	return true
}

type StructListBinding[T BindableStructValue] struct {
	*binding
	list []T
}

func NewStructListBinding[T BindableStructValue](name string, values []T) *StructListBinding[T] {
	b := &StructListBinding[T]{
		binding: newBinding(name),
		list:    make([]T, 0, len(values)),
	}
	b.list = append(b.list, values...)
	return b
}

func (b StructListBinding[T]) Get() []T {
	return b.list
}

func (b StructListBinding[T]) GetAt(index int) (any, bool) {
	if index < 0 || index >= len(b.list) {
		return *new(T), false
	}
	return b.list[index], true
}

func (b StructListBinding[T]) Size() int {
	return len(b.list)
}

func (b *StructListBinding[T]) Set(values []T) {
	if b.equalValues(values) {
		return
	}

	b.list = make([]T, 0, len(values))
	b.list = append(b.list, values...)
	b.notify(b)
}

func (b StructListBinding[T]) equalValues(other []T) bool {
	if len(b.list) != len(other) {
		return false
	}
	for i := range b.list {
		if b.list[i].Less(other[i]) || other[i].Less(b.list[i]) {
			return false
		}
	}
	return true
}
