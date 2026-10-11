# Packages in Go

## 1. Packages, visibility and imports

A package is one directory of Go files that share a name and a scope, and the case of a name's first letter decides whether other packages can see it.

```go
// File: shop/pricing/pricing.go
package pricing // every file in this directory declares the same package

const TaxRate = 0.18 // exported: starts with an uppercase letter

func Total(amount float64) float64 { return amount + tax(amount) } // exported

func tax(amount float64) float64 { return amount * TaxRate } // unexported: package-private
```

```go
// File: shop/main.go
package main

import (
    "fmt" // standard library

    "github.com/google/uuid" // third party

    "github.com/asha/shop/pricing" // your own module: module path plus directory
)

func main() {
    fmt.Println(pricing.Total(100)) // 118
    fmt.Println(uuid.NewString())
    // pricing.tax(100)             // compile error: tax is not exported
}
```

### Import forms

```go
import str "strings"         // alias: call str.ToUpper
import _ "github.com/lib/pq" // blank: run the package's init only, here to register a driver
import . "math"              // dot: use Sqrt with no prefix; avoid outside tests
```

- `package main` with a `main` function builds an executable. Every other package is a library.
- There are only two visibility levels: exported and package-private. The rule also applies to struct fields and methods.
- All files of a package see each other's unexported names, with no import needed.
- By convention the package name equals the last element of its import path.
- Package names are short, lowercase, single words: `http`, `json`, `order`. No underscores or mixed case.
- An unused import is a compile error. The `goimports` tool adds and removes imports for you.

## 2. Modules and dependencies

A module is a set of packages versioned and released together, defined by a `go.mod` file at its root.

```
module github.com/asha/shop            // the module path: prefix of every import in it

go 1.27                                // language version this module is written for

require (
    github.com/google/uuid v1.6.0      // direct dependencies
    golang.org/x/sync v0.10.0
)

require golang.org/x/text v0.21.0 // indirect: needed by a dependency, not imported by you

replace github.com/asha/lib => ../lib  // use a local copy, for development or a fork
```

| Command | Purpose |
| --- | --- |
| `go mod init github.com/asha/shop` | Create `go.mod` for a new module |
| `go get github.com/google/uuid@v1.6.0` | Add, upgrade or downgrade one dependency |
| `go get -u ./...` | Upgrade all dependencies to their latest minor or patch versions |
| `go mod tidy` | Add what is missing and remove what is unused |
| `go list -m all` | List every module in the build |
| `go mod why <module>` | Explain why a module is needed |
| `go mod vendor` | Copy dependencies into a `vendor` directory |
| `go install example.com/tool@latest` | Build and install a command-line tool |
| `go work init ./a ./b` | Develop several local modules together (Go 1.18) |

- Versions follow semantic versioning: `vMAJOR.MINOR.PATCH`.
- From major version 2 onward, the major version is part of the import path, as in `github.com/x/lib/v2`.
- `go.sum` records a checksum for every dependency, so a changed or tampered download is rejected. Commit both files.
- Go uses minimal version selection: it picks the lowest version that satisfies every requirement, which makes builds reproducible without a separate lock file.
- Dependencies are fetched through a module proxy and stored in a shared module cache, not inside your project.
- The `go` line also selects language behaviour. A module that says `go 1.21` keeps the old loop-variable semantics even on a newer toolchain.

## 3. Project layout and the internal directory

Go prescribes no project structure, but most services converge on `cmd` for executables and `internal` for everything that should stay private.

```
shop/
  go.mod
  go.sum
  cmd/
    api/main.go       one directory per executable
    worker/main.go
  internal/
    order/            business logic, importable only from inside shop
      order.go
      order_test.go   tests sit next to the code
    payment/
    store/
  pkg/
    client/           optional: code other modules are meant to import
```

- The compiler enforces `internal`: only code rooted at the parent of that directory can import it.
- Keep `main` thin. It should read configuration, wire dependencies together and start the program.
- Organize packages by what they do (`order`, `payment`), not by what they contain (`models`, `utils`, `helpers`).
- Start with a flat layout and add directories when the code asks for them. Deep trees early on are a common mistake.
- Avoid stutter in names: prefer `order.Service` to `order.OrderService`.
- A test file may declare `package order_test` to test only the exported API, as an outside caller would.

## 4. init and initialization order

Each package is initialized exactly once before `main` runs: first its imports, then its package-level variables, then its `init` functions.

```go
package config

var settings = load() // step 2: package-level variables

func init() { // step 3: init functions
    if settings == nil {
        panic("config: failed to load")
    }
}
```

1. Imported packages are initialized first, dependencies before the packages that use them.
2. Package-level variables are set, in dependency order.
3. The package's `init` functions run, in the order the files are given to the compiler.
4. When every package is ready, `main.main` starts.

- `init` takes no arguments, returns nothing, and cannot be called from code.
- A package, and even one file, may have several `init` functions.
- A blank import such as `import _ "net/http/pprof"` exists purely to run that package's `init`.
- Keep `init` small. Hidden work there is hard to test, hard to order, and slows startup.
- Prefer an explicit constructor such as `NewStore(cfg)` that `main` calls, so dependencies are visible.

## 5. Gotchas

The import cycle is the package problem you will meet most, and interviewers like to ask how you would break one.

```
order  imports  payment
payment  imports  order      compile error: import cycle not allowed
```

```go
// Fix: the consumer defines a small interface, so it no longer imports the other package
package order

type Charger interface {
    Charge(orderID string, amount int) error
}

type Service struct{ charger Charger }

// package payment implements Charge and may import order; order never imports payment
```

- Two ways to break a cycle: move the shared types into a third package, or depend on an interface as above.
- Package-level variables are global state. Tests that share them interfere with each other, so pass dependencies in instead.
- Everything exported is a promise to your users. Keep names unexported until something outside needs them.
- A package named `utils` or `common` grows without limit and attracts cycles. Name packages after what they provide.
- Run `go mod tidy` before committing, and commit `go.sum` together with `go.mod`.
- Private repositories need `GOPRIVATE` set, so the public proxy and checksum database are skipped for them.
- GOPATH mode is legacy. Modules have been the default since Go 1.16.

## 6. Rapid-fire questions

One-line answers for a quick revision.

| Question | Short answer |
| --- | --- |
| What is a package? | A directory of Go files compiled together under one name. |
| What is a module? | A versioned collection of packages with a `go.mod` at its root. |
| How do you export a name? | Start it with an uppercase letter. |
| What is special about `package main`? | It builds an executable and must contain `func main`. |
| What is `go.mod`? | The module's path, Go version and dependency requirements. |
| What is `go.sum`? | Checksums that verify every downloaded dependency. |
| What does `go mod tidy` do? | Adds missing dependencies and removes unused ones. |
| What is a blank import for? | Running a package's `init` without using its names. |
| What is the `internal` directory? | Packages that only the surrounding module tree may import. |
| Are import cycles allowed? | No; break them with a third package or an interface. |
| When does `init` run? | Once per package, after its variables are set and before `main`. |
| What is the order of initialization? | Imports, then package variables, then `init`, then `main`. |
| How are major versions 2 and above imported? | With the version in the path, such as `/v2`. |
| What is minimal version selection? | Choosing the lowest version that satisfies all requirements. |
| What is `go.work`? | A workspace file for developing several modules together. |
