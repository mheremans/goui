// SPDX-License-Identifier: MIT

package types

import (
	"gioui.org/unit"
	"golang.org/x/exp/constraints"
)

type SizeConstraint interface {
	unit.Dp | constraints.Float | constraints.Integer
}

type sizeDP struct {
	Width  unit.Dp
	Height unit.Dp
}

type WindowSize = sizeDP
type SpacerSize = sizeDP

// NewWindowSize creates a new WindowSize struct with the given width and height.
//
// Parameters:
// - w: the width of the window.
// - h: the height of the window.
//
// Returns:
// - WindowSize: a struct representing the window size with the given width and height.
func NewWindowSize[T SizeConstraint](w, h T) WindowSize {
	return WindowSize{
		unit.Dp(w),
		unit.Dp(h),
	}
}

// NewSpacerSize creates a new SpacerSize struct with the given width and height.
//
// Parameters:
// - w: the width of the spacer.
// - h: the height of the spacer.
//
// Returns:
// - SpacerSize: a struct representing the spacer size with the given width and height.
func NewSpacerSize[T SizeConstraint](w, h T) SpacerSize {
	return SpacerSize{
		unit.Dp(w),
		unit.Dp(h),
	}
}
