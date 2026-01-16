# sqlite3

## About this package

This package serves as a **database driver adapter** for Go's `ent` entity framework, enabling it to seamlessly utilize the **`modernc.org/sqlite`** driver as its SQLite backend.

### The Problem

1.  The `ent` framework typically recommends `mattn/go-sqlite3` as its SQLite driver.
2.  `mattn/go-sqlite3` is an excellent driver but relies on **CGO**.
3.  CGO can complicate the Go project's compilation process, especially for **cross-platform compilation**, leading to additional environment setup and dependency issues.

### The Solution

1.  `modernc.org/sqlite` is a **pure Go implementation** of the SQLite driver, completely eliminating CGO dependencies and resolving the aforementioned compilation challenges.
2.  However, `ent` requires a driver that can correctly parse specific Data Source Name (DSN) parameters (e.g., `_fk=1` for foreign key constraints) during code generation and database connection.
3.  **This package (`github.com/sqlite3ent/sqlite3`) acts as the crucial "bridge"**. It registers a driver named `sqlite3` that internally uses `modernc.org/sqlite`. Crucially, it understands `ent`-specific DSN parameters, ensuring `ent` functions correctly and without issues.

### Core Advantages

*   **CGO-Free**: Simplifies your project by removing CGO dependencies.
*   **Effortless Cross-Platform Compilation**: Compile your application for any Go-supported platform without complex configurations.
*   **Seamless `ent` Integration**: From `ent`'s perspective, it operates as if using a standard SQLite driver.

## Version

Current version: v1.41.0 (This will be updated automatically upon new releases)

## Installation

```bash
go get github.com/sqlite3ent/sqlite3@latest
```

## How to use

Using this package in your `ent` project is straightforward, involving two main steps:

### 1. Configure the Driver in `ent` Code Generation

You need to instruct `ent` to use this package as the driver within your project's `ent/generate.go` file.

```go
//go:build ignore

package main

import (
	"log"

	"entgo.io/ent/dialect"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	err := entc.Generate("./schema", &gen.Config{
		// Tell ent we are using an external SQL driver
		Features: []gen.Feature{
			gen.FeatureUpsert,
			gen.FeatureSQLDriver,
		},
		// Configure SQLite dialect
		Dialect: dialect.SQLite,
		// Specify the driver's import path
		Driver: "github.com/sqlite3ent/sqlite3",
	})
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
```

After configuration, run `go generate ./...` to regenerate your `ent` code.

### 2. Connect to the Database in Your Application

In your `main.go` or other application entry file, anonymously import this driver, then use `ent.Open` as usual.

```go
package main

import (
	"context"
	"log"

	"your/project/ent" // Replace with your project's ent path

	// Crucial step: Anonymous import of this driver.
	// This executes the driver's init() function,
	// registering it with Go's database/sql system.
	_ "github.com/sqlite3ent/sqlite3"
)

func main() {
	// DSN (Data Source Name) format is compatible with mattn/go-sqlite3.
	// Example: using an in-memory database with foreign keys enabled.
	dsn := "file:ent?mode=memory&cache=shared&_fk=1"

	// Use ent.Open to connect to the database.
	// ent will automatically use our registered "sqlite3" driver.
	client, err := ent.Open(dialect.SQLite, dsn)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	// Run database migrations (e.g., create tables)
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// ... Your application logic goes here ...
	log.Println("Successfully connected to SQLite without CGO!")
}
```

## LICENSE

Used BSD-3-Clause is same as `modernc.org/sqlite`
