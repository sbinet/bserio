// Copyright 2018 The go-boostio Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binser_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/go-boostio/boostio/binser"
)

func TestRead(t *testing.T) {
	type animal struct {
		Name  string
		Legs  int16
		Tails int8
	}

	f, err := os.Open("testdata/data.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	dec := binser.NewDecoder(f)
	for _, tc := range []struct {
		name string
		want interface{}
	}{
		{"bool-false", false},
		{"bool-true", true},
		{"int8", int8(0x11)},
		{"int16", int16(0x2222)},
		{"int32", int32(0x33333333)},
		{"int64", int64(0x4444444444444444)},
		{"uint8", uint8(0xff)},
		{"uint16", uint16(0x2222)},
		{"uint32", uint32(0x3333333)},
		{"uint64", uint64(0x444444444444444)},
		{"float32", float32(2.2)},
		{"float64", 3.3},
		{"[3]uint8", [3]uint8{0x11, 0x22, 0x33}},
		{"[]uint8", []uint8{0x11, 0x22, 0x33, 0xff}},
		{"[]byte", []byte("hello")},
		{"string", "hello"},
		{"map[string]string", map[string]string{"eins": "un", "zwei": "deux", "drei": "trois"}},
		{"struct", animal{"pet", 4, 1}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rv := reflect.New(reflect.TypeOf(tc.want)).Elem()
			if rv.Kind() == reflect.Map {
				rv.Set(reflect.MakeMap(rv.Type()))
			}
			err := dec.Decode(rv.Addr().Interface())
			if err != nil {
				t.Fatalf("could not read %q: %v", tc.name, err)
			}
			if got, want := rv.Interface(), tc.want; !reflect.DeepEqual(got, want) {
				t.Fatalf("got=%#v (%T)\nwant=%#v (%T)", got, got, want, want)
			}
		})
	}
}

type animal struct {
	Name  string
	Legs  int16
	Tails int8
}

func (a *animal) UnmarshalBoost(r *binser.Reader) error {
	/*vers*/ _ = r.ReadU32()
	/*flag*/ _ = r.ReadU8()
	a.Name = r.ReadString()
	a.Legs = r.ReadI16()
	a.Tails = r.ReadI8()
	return r.Err()
}

func TestUnmarshaler(t *testing.T) {
	f, err := os.Open("testdata/data.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	dec := binser.NewDecoder(f)
	for _, tc := range []struct {
		name string
		want interface{}
	}{
		{"bool-false", false},
		{"bool-true", true},
		{"int8", int8(0x11)},
		{"int16", int16(0x2222)},
		{"int32", int32(0x33333333)},
		{"int64", int64(0x4444444444444444)},
		{"uint8", uint8(0xff)},
		{"uint16", uint16(0x2222)},
		{"uint32", uint32(0x3333333)},
		{"uint64", uint64(0x444444444444444)},
		{"float32", float32(2.2)},
		{"float64", 3.3},
		{"[3]uint8", [3]uint8{0x11, 0x22, 0x33}},
		{"[]uint8", []uint8{0x11, 0x22, 0x33, 0xff}},
		{"[]byte", []byte("hello")},
		{"string", "hello"},
		{"map[string]string", map[string]string{"eins": "un", "zwei": "deux", "drei": "trois"}},
		{"struct", animal{"pet", 4, 1}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rv := reflect.New(reflect.TypeOf(tc.want)).Elem()
			if rv.Kind() == reflect.Map {
				rv.Set(reflect.MakeMap(rv.Type()))
			}
			err := dec.Decode(rv.Addr().Interface())
			if err != nil {
				t.Fatalf("could not read %q: %v", tc.name, err)
			}
			if got, want := rv.Interface(), tc.want; !reflect.DeepEqual(got, want) {
				t.Fatalf("got=%#v (%T)\nwant=%#v (%T)", got, got, want, want)
			}
		})
	}
}