// Package forklift provides a simple way to load Go packages.
//
// There are three types of packages:
// normal packages, test packages, and external test packages.
// A normal package is the result of a "go build" command.
// A test package is the result of a "go test" command,
// excluding test files that declare a package name with a "_test" suffix.
// An external test package is the result of a "go test" command,
// excluding test files that do not declare a package name with a "_test" suffix.
//
// To load the normal package in the current directory:
//
//	p, err := forklift.LoadPackage(".")
//
// To load the test package in the "time" package:
//
//	p, err := forklift.LoadTestPackage("time")
//
// To load the external test package in the "strings" package:
//
//	p, err := forklift.LoadExternalTestPackage("strings")
//
// The result is a [*golang.org/x/tools/go/packages.Package].
//
// Paths are passed directly to [golang.org/x/tools/go/packages.Load].
// All information is loaded.
//
// To configure the loading behavior, use [Loader].
package forklift

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Loader provides Packages for import paths.
type Loader struct {
	// Context is used if set.
	Context context.Context

	// Dir is the build system working directory. It defaults to the current one.
	Dir string

	// Env is the build system environment variables.
	Env []string

	// Flags is the build system command-line flags.
	Flags []string

	// Mode is the information to include.
	Mode packages.LoadMode
}

func loadError(err error) error {
	return fmt.Errorf("cannot load package: %v", err)
}

// ErrNotFound neans the package was not found.
var ErrNotFound = fmt.Errorf("package not found")

func handle(p *packages.Package) (*packages.Package, error) {
	if p == nil {
		return nil, ErrNotFound
	}
	var errs []error
	for _, err := range p.Errors {
		switch err.Kind {
		case packages.ListError:
			return nil, ErrNotFound
		case packages.ParseError, packages.TypeError:
			var prefix string
			if err.Pos != "" && err.Pos != "-" {
				prefix = err.Pos + ": "
			}
			errs = append(errs, errors.New(prefix+err.Msg))
		default:
			panic(err.Kind)
		}
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return p, nil
}

// LoadPackage returns the package for path.
// It returns [ErrNotFound] if the package is not found, and other errors.
func (l Loader) LoadPackage(path string) (*packages.Package, error) {
	ps, err := packages.Load(&packages.Config{Context: l.Context, Dir: l.Dir, Env: l.Env, BuildFlags: l.Flags, Mode: l.Mode}, path)
	if err != nil {
		return nil, loadError(err)
	}
	var match *packages.Package
loop:
	for _, p := range ps {
		if strings.HasSuffix(p.Name, "_test") {
			continue
		}
		for _, f := range p.GoFiles {
			if strings.HasSuffix(f, "_test.go") {
				continue loop
			}
		}
		match = p
		break
	}
	return handle(match)
}

// LoadTestPackage returns the test package for path.
// It returns [ErrNotFound] if the package is not found, and other errors.
func (l Loader) LoadTestPackage(path string) (*packages.Package, error) {
	ps, err := packages.Load(&packages.Config{Context: l.Context, Dir: l.Dir, Env: l.Env, BuildFlags: l.Flags, Mode: l.Mode, Tests: true}, path)
	if err != nil {
		return nil, loadError(err)
	}
	var match *packages.Package
loop:
	for _, p := range ps {
		if strings.HasSuffix(p.Name, "_test") {
			continue
		}
		for _, f := range p.GoFiles {
			if strings.HasSuffix(f, "_test.go") {
				match = p
				break loop
			}
		}
	}
	return handle(match)
}

// LoadExternalTestPackage returns the external test package for path.
// It returns [ErrNotFound] if the package is not found, and other errors.
func (l Loader) LoadExternalTestPackage(path string) (*packages.Package, error) {
	ps, err := packages.Load(&packages.Config{Context: l.Context, Dir: l.Dir, Env: l.Env, BuildFlags: l.Flags, Mode: l.Mode, Tests: true}, path)
	if err != nil {
		return nil, loadError(err)
	}
	var match *packages.Package
	for _, p := range ps {
		if strings.HasSuffix(p.Name, "_test") {
			match = p
			break
		}
	}
	return handle(match)
}

var mode packages.LoadMode = packages.NeedCompiledGoFiles |
	packages.NeedDeps |
	packages.NeedEmbedFiles |
	packages.NeedEmbedPatterns |
	packages.NeedFiles |
	packages.NeedImports |
	packages.NeedModule |
	packages.NeedName |
	packages.NeedSyntax |
	packages.NeedTypes |
	packages.NeedTypesInfo |
	packages.NeedTypesSizes

// LoadPackage returns the package for path.
// It returns [ErrNotFound] if the package is not found, and other errors.
func LoadPackage(path string) (*packages.Package, error) {
	return Loader{Mode: mode}.LoadPackage(path)
}

// LoadTestPackage returns the test package for path.
// It returns [ErrNotFound] if the package is not found, and other errors.
func LoadTestPackage(path string) (*packages.Package, error) {
	return Loader{Mode: mode}.LoadTestPackage(path)
}

// LoadExternalTestPackage returns the external test package for path.
// It returns [ErrNotFound] if the package is not found, and other errors.
func LoadExternalTestPackage(path string) (*packages.Package, error) {
	return Loader{Mode: mode}.LoadExternalTestPackage(path)
}
