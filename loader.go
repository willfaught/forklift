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

// LoadPackage returns the package for importPath, or nil if it does not exist.
func (l Loader) LoadPackage(importPath string) (*packages.Package, error) {
	ps, err := packages.Load(&packages.Config{
		Context:    l.Context,
		Dir:        l.Dir,
		Env:        l.Env,
		BuildFlags: l.Flags,
		Mode:       l.Mode,
	}, importPath)
	if err != nil {
		return nil, fmt.Errorf("cannot load package %s: %v", importPath, err)
	}
	if len(ps) == 0 {
		return nil, nil
	}
	if errs := ps[0].Errors; len(errs) > 0 {
		return nil, fmt.Errorf("cannot load package %s: %v", importPath, errs[0])
	}
	return ps[0], nil
}

// LoadTestPackage returns the test package for importPath, or nil if it does not exist.
func (l Loader) LoadTestPackage(importPath string) (*packages.Package, error) {
	ps, err := packages.Load(&packages.Config{
		Context:    l.Context,
		Dir:        l.Dir,
		Env:        l.Env,
		BuildFlags: l.Flags,
		Mode:       l.Mode,
		Tests:      true,
	}, importPath)
	if err != nil {
		return nil, fmt.Errorf("cannot load package %s: %v", importPath, err)
	}
	for _, p := range ps {
		if p.PkgPath == importPath {
			for _, f := range p.GoFiles {
				if strings.HasSuffix(f, "_test.go") {
					if errs := p.Errors; len(errs) > 0 {
						return nil, fmt.Errorf("cannot load package %s: %v", importPath, errs[0])
					}
					return p, nil
				}
			}
		}
	}
	return nil, nil
}

// LoadExternalTestPackage returns the external test package for importPath, or nil if it does not exist.
func (l Loader) LoadExternalTestPackage(importPath string) (*packages.Package, error) {
	ps, err := packages.Load(&packages.Config{
		Context:    l.Context,
		Dir:        l.Dir,
		Env:        l.Env,
		BuildFlags: l.Flags,
		Mode:       l.Mode,
		Tests:      true,
	}, importPath)
	if err != nil {
		return nil, fmt.Errorf("cannot load package %s: %v", importPath, err)
	}
	for _, p := range ps {
		if strings.HasSuffix(p.Name, "_test") {
			if errs := p.Errors; len(errs) > 0 {
				return nil, fmt.Errorf("cannot load package %s: %v", importPath, errs[0])
			}
			return p, nil
		}
	}
	return nil, nil
}

var mode packages.LoadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedTypesSizes |
	packages.NeedModule |
	packages.NeedEmbedFiles |
	packages.NeedEmbedPatterns

// LoadPackage returns the package for importPath, or nil if it does not exist.
func LoadPackage(importPath string) (*packages.Package, error) {
	return Loader{Mode: mode}.LoadPackage(importPath)
}

// LoadTestPackage returns the test package for importPath, or nil if it does not exist.
func LoadTestPackage(importPath string) (*packages.Package, error) {
	return Loader{Mode: mode}.LoadTestPackage(importPath)
}

// LoadExternalTestPackage returns the external test package for importPath, or nil if it does not exist.
func LoadExternalTestPackage(importPath string) (*packages.Package, error) {
	return Loader{Mode: mode}.LoadExternalTestPackage(importPath)
}
