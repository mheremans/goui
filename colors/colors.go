// SPDX-License-Identifier: MIT

package colors

import (
	"errors"
	"image/color"
	"strconv"
	"strings"

	"golang.org/x/exp/shiny/materialdesign/colornames"
)

func GetColor(str string) (color.NRGBA, error) {
	if str[0] == '#' {
		return getColorFromHex(str)
	}
	return getColorFromName(str)
}

func getColorFromHex(hex string) (color.NRGBA, error) {
	hex = strings.Replace(hex, "#", "", -1)

	chnls := make([]uint8, 4)
	chnls[3] = 255
	numIts := len(hex) / 2
	if numIts > 4 {
		numIts = 4
	}

	for i := 0; i < numIts; i++ {
		ch := hex[i*2 : i*2+2]
		v, err := strconv.ParseUint(ch, 16, 8)
		if err != nil {
			return color.NRGBA{}, err
		}
		chnls[i] = uint8(v)
	}
	return color.NRGBA{
		R: chnls[0],
		G: chnls[1],
		B: chnls[2],
		A: chnls[3],
	}, nil
}

func getColorFromName(name string) (color.NRGBA, error) {
	c, ok := colornames.Map[name]
	if !ok {
		return color.NRGBA{}, errors.New("color not found")
	}
	return color.NRGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: c.A,
	}, nil
}
