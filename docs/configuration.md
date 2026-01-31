# Configuration Loading

The staticgen tool now supports loading variables from a `.staticgen.yml` configuration file.

## Configuration File Format

Create a `.staticgen.yml` file in your project root:

```yaml
variables:
  siteName: "My Static Site"
  author: "Your Name"
  year: "2026"
  # Add any other variables you need
```

## Usage

### Command Line

The configuration file is loaded automatically when running staticgen:

```bash
# Uses .staticgen.yml by default
staticgen

# Or specify a custom config file
staticgen --config myconfig.yml
```

### In Code

Variables from the config file are automatically loaded into the variable store:

```go
import "github.com/blamarvt/staticgen/pkg/vars"

// Load config
config, err := vars.LoadConfig(".staticgen.yml")
if err != nil {
    log.Fatal(err)
}

// Create store and load variables
variables := vars.NewStore()
variables.LoadFromConfig(config)

// Variables are now available in templates
```

### In Templates

Use variables in your component templates with the `Var` function:

```html
<footer>
  <p>© {{ Var "year" }} {{ Var "siteName" }}</p>
  <p>By {{ Var "author" }}</p>
</footer>
```

## Behavior

- If the config file is not found, staticgen continues with an empty variable store
- If the config file exists but has invalid YAML, a warning is logged
- Variables can be overridden programmatically after loading from config
- All variables are strings

## Example

See the [integration test](../../tests/vars_test.go) for a complete example of loading and using config variables.
