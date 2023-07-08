package forklift

import (
	"testing"

	"golang.org/x/tools/go/packages"
)

func check(t *testing.T, p *packages.Package, err error, found bool) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
	if found {
		if p == nil {
			t.Error("package is nil")
		}
	} else {
		if p != nil {
			t.Error("package is not nil")
		}
	}
}

func TestLoadPackage(t *testing.T) {
	t.Parallel()
	p, err := LoadPackage("time")
	check(t, p, err, true)
	p, err = LoadPackage(".")
	check(t, p, err, true)
}

func TestLoadTestPackage(t *testing.T) {
	t.Parallel()
	p, err := LoadTestPackage("time")
	check(t, p, err, true)
}

func TestLoadExternalTestPackage(t *testing.T) {
	t.Parallel()
	p, err := LoadExternalTestPackage("time")
	check(t, p, err, true)
}
