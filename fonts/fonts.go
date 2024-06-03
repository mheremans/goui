// SPDX-License-Identifier: MIT

package fonts

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"

	giofont "gioui.org/font"
	"gioui.org/font/opentype"
	"github.com/mheremans/goui/types"
)

var fontCache map[string][]giofont.FontFace

func init() {
	fontCache = make(map[string][]giofont.FontFace)
}

func GetOTFFont(
	ctx types.Context,
	file string,
	style giofont.Style,
	weight giofont.Weight,
) (
	font giofont.FontFace,
	err error,
) {
	fonts, err := GetOTFFonts(ctx, file)
	if err != nil {
		return
	}

	fonts = filterOnStyle(fonts, style)
	fonts = filterOnWeight(fonts, weight, 0)

	if len(fonts) == 0 {
		err = errors.New("no fonts found")
		return
	}
	font = fonts[0]
	return
}

func GetOTFFonts(
	ctx types.Context,
	file string,
) (
	fonts []giofont.FontFace,
	err error,
) {
	var ok bool

	if fonts, ok = fontCache[file]; ok {
		return
	}

	if ctx.FontsDir() == nil {
		err = errors.New("no fonts found")
		return
	}

	if fonts, err = newOTFFont(*ctx.FontsDir(), file); err != nil {
		return
	}
	fontCache[file] = fonts
	return
}

func newOTFFont(
	filesystem embed.FS,
	file string,
) (
	fonts []giofont.FontFace,
	err error,
) {
	var bytes []byte
	var fh fs.File

	if fh, err = filesystem.Open(file); err != nil {
		err = fmt.Errorf("no fonts found: %w", err)
		return
	}
	defer fh.Close()

	if bytes, err = io.ReadAll(fh); err != nil {
		err = fmt.Errorf("unable to read fontes: %w", err)
		return
	}

	return opentype.ParseCollection(bytes)
}

// Try to filter on the style
// If no styles match we try to filter on regular fonts
// If still nothing matches, we return everything
func filterOnStyle(
	fonts []giofont.FontFace,
	style giofont.Style,
) (
	res []giofont.FontFace,
) {
	res = make([]giofont.FontFace, 0, len(fonts))
	for _, f := range fonts {
		if f.Font.Style == style {
			res = append(res, f)
		}
	}
	if len(res) == 0 {
		res = filterOnStyle(fonts, giofont.Regular)
		if len(res) == 0 {
			res = fonts
		}
	}
	return res
}

func filterOnWeight(
	fonts []giofont.FontFace,
	weight giofont.Weight,
	offset int,
) (
	res []giofont.FontFace,
) {
	target := giofont.Weight(int(weight) + offset)

	res = make([]giofont.FontFace, 0, len(fonts))
	for _, f := range fonts {
		if f.Font.Weight == target {
			res = append(res, f)
		}
	}
	if offset == 0 && len(res) == 0 {
		offPlus := offset
		offMin := offset
		for len(res) == 0 && offPlus < 500 {
			offPlus += 100
			offMin -= 100

			res1 := filterOnWeight(fonts, weight, offPlus)
			res2 := filterOnWeight(fonts, weight, offMin)
			res = append(res, res1...)
			res = append(res, res2...)
		}

		if len(res) == 0 {
			res = fonts
		}
	}
	return res
}
