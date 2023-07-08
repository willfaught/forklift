// Package forklift provides a streamlined experience for parsing Go packages.
package forklift

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Package is file information for a Go package.
type Package struct {
	Files     []*ast.File
	Positions *token.FileSet
}

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
}

// LoadPackage returns the package for importPath.
func (l Loader) LoadPackage(importPath string) (*Package, error) {
	fset := token.NewFileSet()
	config := &packages.Config{
		Context:    l.Context,
		Dir:        l.Dir,
		Env:        l.Env,
		BuildFlags: l.Flags,
		Fset:       fset,
		Mode:       packages.NeedSyntax,
	}
	ps, err := packages.Load(config, importPath)
	if err != nil {
		return nil, fmt.Errorf("cannot load package %s: %v", importPath, err)
	}
	if len(ps) == 0 {
		return nil, nil
	}
	return &Package{Files: ps[0].Syntax, Positions: fset}, nil
}

// LoadTestPackage returns the test package for importPath.
func (l Loader) LoadTestPackage(importPath string) (*Package, error) {
	fset := token.NewFileSet()
	config := &packages.Config{
		Context:    l.Context,
		Dir:        l.Dir,
		Env:        l.Env,
		BuildFlags: l.Flags,
		Fset:       fset,
		Tests:      true,
		Mode:       packages.NeedFiles | packages.NeedName | packages.NeedSyntax,
	}
	ps, err := packages.Load(config, importPath)
	if err != nil {
		return nil, fmt.Errorf("cannot load package %s: %v", importPath, err)
	}
	for _, p := range ps {
		if p.PkgPath == importPath {
			for _, f := range p.GoFiles {
				if strings.HasSuffix(f, "_test.go") {
					return &Package{Files: p.Syntax, Positions: fset}, nil
				}
			}
		}
	}
	return nil, nil
}

// LoadExternalTestPackage returns the external test package for importPath.
func (l Loader) LoadExternalTestPackage(importPath string) (*Package, error) {
	fset := token.NewFileSet()
	config := &packages.Config{
		Context:    l.Context,
		Dir:        l.Dir,
		Env:        l.Env,
		BuildFlags: l.Flags,
		Fset:       fset,
		Tests:      true,
		Mode:       packages.NeedFiles | packages.NeedName | packages.NeedSyntax,
	}
	ps, err := packages.Load(config, importPath)
	if err != nil {
		return nil, fmt.Errorf("cannot load package %s: %v", importPath, err)
	}
	pkgPath := importPath + "_test"
	for _, p := range ps {
		if p.PkgPath == pkgPath {
			return &Package{Files: p.Syntax, Positions: fset}, nil
		}
	}
	return nil, nil
}

// LoadPackage returns the package for importPath.
func LoadPackage(importPath string) (*Package, error) {
	return Loader{}.LoadPackage(importPath)
}

// LoadTestPackage returns the test package for importPath.
func LoadTestPackage(importPath string) (*Package, error) {
	return Loader{}.LoadTestPackage(importPath)
}

// LoadExternalTestPackage returns the external test package for importPath.
func LoadExternalTestPackage(importPath string) (*Package, error) {
	return Loader{}.LoadExternalTestPackage(importPath)
}
