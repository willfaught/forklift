// Package forklift provides a streamlined experience for parsing Go packages.
package forklift

import (
	"context"
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

// LoadPackage returns the package for path, or nil if it does not exist.
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
	if match == nil {
		return nil, nil
	}
	for _, err := range match.Errors {
		if err.Kind == packages.ListError {
			return nil, nil
		}
	}
	return match, nil
}

// LoadTestPackage returns the test package for path, or nil if it does not exist.
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
	if match == nil {
		return nil, nil
	}
	for _, err := range match.Errors {
		if err.Kind == packages.ListError {
			return nil, nil
		}
	}
	return match, nil
}

// LoadExternalTestPackage returns the external test package for path, or nil if it does not exist.
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
	if match == nil {
		return nil, nil
	}
	for _, err := range match.Errors {
		if err.Kind == packages.ListError {
			return nil, nil
		}
	}
	return match, nil
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

// LoadPackage returns the package for path, or nil if it does not exist.
func LoadPackage(path string) (*packages.Package, error) {
	return Loader{Mode: mode}.LoadPackage(path)
}

// LoadTestPackage returns the test package for path, or nil if it does not exist.
func LoadTestPackage(path string) (*packages.Package, error) {
	return Loader{Mode: mode}.LoadTestPackage(path)
}

// LoadExternalTestPackage returns the external test package for path, or nil if it does not exist.
func LoadExternalTestPackage(path string) (*packages.Package, error) {
	return Loader{Mode: mode}.LoadExternalTestPackage(path)
}
