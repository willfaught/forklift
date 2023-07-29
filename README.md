# forklift

Package forklift provides a simple way to load Go packages.

There are three types of packages: normal packages, test packages, and external test packages. A normal package is the result of a "go build" command. A test package is the result of a "go test" command, excluding test files that declare a package name with a "\_test" suffix. An external test package is the result of a "go test" command, excluding test files that do not declare a package name with a "\_test" suffix.

To load the normal package in the current directory:

	p, err := forklift.LoadPackage(".")

To load the test package in the "time" package:

	p, err := forklift.LoadTestPackage("time")

To load the external test package in the "strings" package:

	p, err := forklift.LoadExternalTestPackage("strings")

The result is a [\*golang.org/x/tools/go/packages.Package](https://pkg.go.dev/golang.org/x/tools/go/packages#Package).

Paths are passed directly to [golang.org/x/tools/go/packages.Load](https://pkg.go.dev/golang.org/x/tools/go/packages#Load). All information is loaded.

To configure the loading behavior, use [Loader](https://pkg.go.dev/github.com/willfaught/forklift#Loader).
