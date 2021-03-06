package main

import (
	"github.com/Grafikart/subsearch/opensubtitle"
	"testing"
)

func EqualSlice(t *testing.T, a []string, b []string) {
	t.Helper()
	if len(a) != len(b) {
		t.Errorf("%v !== %v", a, b)
	}
	for i, v := range a {
		if v != b[i] {
			t.Errorf("%v[%v] !== %v[%v]", a, v, b, b[i])
		}
	}
}

func TestGetKeys(t *testing.T) {
	m := map[string]*opensubtitle.Subtitle{
		"Key":  {},
		"Key2": {},
	}
	k := getKeys(m)
	EqualSlice(t, k, []string{"Key", "Key2"})
}
