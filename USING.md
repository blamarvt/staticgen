# Using staticgen

## What is staticgen?

Staticgen is a static site generator that uses a component-based architecture to generate HTML pages. It allows you to:

- **Define reusable components** with templates that can accept dynamic attributes
- **Compose pages** by combining components with static HTML
- **Mix component and plain HTML** - only template what needs to be dynamic
- **Generate static HTML files** that can be hosted on any web server

This is especially useful for building sites with HTMX or other lightweight frameworks where you want server-side rendering without a full framework.

## Component Architecture

Staticgen uses two main concepts:

1. **Component Definitions** - Reusable templates (like functions)
2. **Component Instances** - Specific uses of those templates with attributes (like function calls)

## Example Project Structure

```
my-site/
├── components/              # Component definitions
│   ├── greeting.hcml       # Reusable greeting component
│   ├── container.hcml      # Layout container
│   └── form.hcml           # Form component
├── pages/                   # Page files that use components
│   ├── index.hcml          # Homepage
│   ├── about.hcml          # About page
│   └── users/
│       ├── add.hcml        # Add user form
│       └── list.hcml       # User list
└── output/                  # Generated HTML files
```

## Example Component Definition

**components/greeting.hcml:**

```xml
<?hcml version="1.0" encoding="UTF-8"?>
<greeting hcmlns="staticgen:components">
    <div class="greeting">
        <h1>Hello, {{ .Name }}!</h1>
        <p>{{ .Message }}</p>
    </div>
</greeting>
```

## Example Page

**pages/index.hcml:**

```xml
<page hcmlns="staticgen"
      hcmlns:component="staticgen:components"
      title="Welcome"
      path="/index.html">
    <component:greeting name="World" message="Welcome to our site!" />

    <div class="content">
        <p>This is plain HTML - no template needed!</p>
    </div>

    <component:greeting name="User" message="Thanks for visiting!" />
</page>
```

## Generated Output

The above page would generate HTML like:

```html
<!DOCTYPE html>
<html>
  <head>
    <title>Welcome</title>
  </head>
  <body>
    <div class="greeting">
      <h1>Hello, World!</h1>
      <p>Welcome to our site!</p>
    </div>

    <div class="content">
      <p>This is plain HTML - no template needed!</p>
    </div>

    <div class="greeting">
      <h1>Hello, User!</h1>
      <p>Thanks for visiting!</p>
    </div>
  </body>
</html>
```

## Running staticgen

```bash
# Build the project
make build

# Generate HTML from your pages
./bin/staticgen
```

The tool will:

1. Load all component definitions from the configured components directory
2. Process each page file
3. Generate static HTML files in the output directory
