package forklift

import (
	"testing"

	"golang.org/x/tools/go/packages"
)

func check(t *testing.T, p *packages.Package, err error, exists bool) {
	t.Helper()
	if exists {
		if p == nil {
			t.Error("package is nil")
		}
		if err != nil {
			t.Error("error is not nil:", err)
		}
	} else {
		if p != nil {
			t.Error("package is not nil")
		}
		if err != nil {
			t.Error("error is nil")
		}
	}
}

func TestLoadPackage(t *testing.T) {
	t.Parallel()
	p, err := LoadPackage("time")
	check(t, p, err, true)
	p, err = LoadPackage(".")
	check(t, p, err, true)
	p, err = LoadPackage("bad")
	check(t, p, err, false)
}

func TestLoadTestPackage(t *testing.T) {
	t.Parallel()
	p, err := LoadTestPackage("time")
	check(t, p, err, true)
	p, err = LoadPackage(".")
	check(t, p, err, true)
	p, err = LoadPackage("bad")
	check(t, p, err, false)
}

func TestLoadExternalTestPackage(t *testing.T) {
	t.Parallel()
	p, err := LoadExternalTestPackage("time")
	check(t, p, err, true)
	p, err = LoadExternalTestPackage(".")
	check(t, p, err, false)
	p, err = LoadPackage("bad")
	check(t, p, err, false)
}
