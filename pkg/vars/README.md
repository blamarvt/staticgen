# Variables Package

The `vars` package provides variable functionality for the staticgen static site generator, allowing you to define and use variables throughout your page templates.

## Overview

The `Store` struct holds a map of string key-value pairs that can be accessed within component templates using the `Var` function.

## Usage

### Creating a Variable Store

```go
import "github.com/blamarvt/staticgen/pkg/vars"

// Create a new store (empty by default)
variables := vars.NewStore()
```

### Setting Variables

```go
variables.Set("siteName", "My Awesome Site")
variables.Set("year", "2026")
variables.Set("author", "John Doe")
```

### Using Variables in Templates

Variables can be used within component templates using the `Var` function:

```html
<footer hcmlns="staticgen:components">
  <footer class="site-footer">
    <p>© {{ Var "year" }} {{ Var "siteName" }}. All rights reserved.</p>
    <p>Author: {{ Var "author" }}</p>
  </footer>
</footer>
```

### Generating Pages with Variables

The `page.Generate` function accepts a `vars.Store` as its third parameter:

```go
html, err := page.Generate(p, registry, variables)
```

## API Reference

### Store

#### NewStore() \*Store

Creates a new Store with an empty map of variables.

#### Set(key, value string)

Sets a variable value in the store.

```go
variables.Set("theme", "dark")
```

#### Get(key string) (string, bool)

Retrieves a variable value from the store. Returns the value and a boolean indicating whether the key exists.

```go
value, ok := variables.Get("theme")
if ok {
    // Key exists, use value
}
```

#### GetOrDefault(key, defaultValue string) string

Retrieves a variable value from the store, or returns a default value if not found.

```go
theme := variables.GetOrDefault("theme", "light")
```

## Template Functions

### Var(name string) string

Available in all component templates. Returns the value of the variable with the given name, or an empty string if the variable is not defined.

```html
{{ Var "siteName" }}
```

## Example

See [vars_test.go](../../tests/vars_test.go) for a complete example demonstrating variable usage.
