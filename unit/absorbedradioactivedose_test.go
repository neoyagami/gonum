// Code generated by "go generate github.com/neoyagami/gonum/unit; DO NOT EDIT.

// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unit

import (
	"fmt"
	"testing"
)

func TestAbsorbedRadioactiveDose(t *testing.T) {
	t.Parallel()
	for _, value := range []float64{-1, 0, 1} {
		var got AbsorbedRadioactiveDose
		err := got.From(AbsorbedRadioactiveDose(value).Unit())
		if err != nil {
			t.Errorf("unexpected error for %T conversion: %v", got, err)
		}
		if got != AbsorbedRadioactiveDose(value) {
			t.Errorf("unexpected result from round trip of %T(%v): got: %v want: %v", got, float64(value), got, value)
		}
		if got != got.AbsorbedRadioactiveDose() {
			t.Errorf("unexpected result from self interface method call: got: %#v want: %#v", got, value)
		}
		err = got.From(ether(1))
		if err == nil {
			t.Errorf("expected error for ether to %T conversion", got)
		}
	}
}

func TestAbsorbedRadioactiveDoseFormat(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		value  AbsorbedRadioactiveDose
		format string
		want   string
	}{
		{1.23456789, "%v", "1.23456789 Gy"},
		{1.23456789, "%.1v", "1 Gy"},
		{1.23456789, "%20.1v", "                1 Gy"},
		{1.23456789, "%20v", "       1.23456789 Gy"},
		{1.23456789, "%1v", "1.23456789 Gy"},
		{1.23456789, "%#v", "unit.AbsorbedRadioactiveDose(1.23456789)"},
		{1.23456789, "%s", "%!s(unit.AbsorbedRadioactiveDose=1.23456789 Gy)"},
	} {
		got := fmt.Sprintf(test.format, test.value)
		if got != test.want {
			t.Errorf("Format %q %v: got: %q want: %q", test.format, float64(test.value), got, test.want)
		}
	}
}
