// Code generated by "go generate github.com/neoyagami/gonum/unit”; DO NOT EDIT.

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

// Resistance represents an electrical resistance, impedance or reactance in Ohms.
type Resistance float64

const Ohm Resistance = 1

// Unit converts the Resistance to a *Unit.
func (r Resistance) Unit() *Unit {
	return New(float64(r), Dimensions{
		CurrentDim: -2,
		LengthDim:  2,
		MassDim:    1,
		TimeDim:    -3,
	})
}

// Resistance allows Resistance to implement a Resistancer interface.
func (r Resistance) Resistance() Resistance {
	return r
}

// From converts the unit into the receiver. From returns an
// error if there is a mismatch in dimension.
func (r *Resistance) From(u Uniter) error {
	if !DimensionsMatch(u, Ohm) {
		*r = Resistance(math.NaN())
		return errors.New("unit: dimension mismatch")
	}
	*r = Resistance(u.Unit().Value())
	return nil
}

func (r Resistance) Format(fs fmt.State, c rune) {
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
		const unit = " Ω"
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
		fmt.Fprintf(fs, "%%!%c(%T=%g Ω)", c, r, float64(r))
	}
}
