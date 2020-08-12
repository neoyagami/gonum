// Code generated by "go generate neoyagami/gonum/unit”; DO NOT EDIT.

// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unit

import (
	"errors"
	"fmt"
	"math"
	"unicode/utf8"
)

// Radioactivity represents a rate of radioactive decay in becquerels.
type Radioactivity float64

const Becquerel Radioactivity = 1

// Unit converts the Radioactivity to a *Unit.
func (r Radioactivity) Unit() *Unit {
	return New(float64(r), Dimensions{
		TimeDim: -1,
	})
}

// Radioactivity allows Radioactivity to implement a Radioactivityer interface.
func (r Radioactivity) Radioactivity() Radioactivity {
	return r
}

// From converts the unit into the receiver. From returns an
// error if there is a mismatch in dimension.
func (r *Radioactivity) From(u Uniter) error {
	if !DimensionsMatch(u, Becquerel) {
		*r = Radioactivity(math.NaN())
		return errors.New("unit: dimension mismatch")
	}
	*r = Radioactivity(u.Unit().Value())
	return nil
}

func (r Radioactivity) Format(fs fmt.State, c rune) {
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T(%v)", r, float64(r))
			return
		}
		fallthrough
	case 'e', 'E', 'f', 'F', 'g', 'G':
		p, pOk := fs.Precision()
		w, wOk := fs.Width()
		const unit = " Bq"
		switch {
		case pOk && wOk:
			fmt.Fprintf(fs, "%*.*"+string(c), pos(w-utf8.RuneCount([]byte(unit))), p, float64(r))
		case pOk:
			fmt.Fprintf(fs, "%.*"+string(c), p, float64(r))
		case wOk:
			fmt.Fprintf(fs, "%*"+string(c), pos(w-utf8.RuneCount([]byte(unit))), float64(r))
		default:
			fmt.Fprintf(fs, "%"+string(c), float64(r))
		}
		fmt.Fprint(fs, unit)
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%g Bq)", c, r, float64(r))
	}
}
