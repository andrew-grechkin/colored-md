#!/usr/bin/env -S colored-md

# Executable Markdown Example

This is a demonstration of `colored-md`'s executable markdown feature.

## How It Works

This file has a shebang on the first line:

```
#!/usr/bin/env -S colored-md
```

When you make this file executable with `chmod +x`, you can run it directly:

```bash
./executable.md
```

The shebang line tells your shell to use `colored-md` as the interpreter. When `colored-md` processes executable files, it automatically strips the shebang line from the output, so you only see the rendered content.

## Feature Details

- **Only affects executable files**: Non-executable markdown files will show the shebang as regular content
- **Only strips colored-md shebangs**: Shebangs for other interpreters (like `#!/bin/bash`) are preserved
- **Works with any shebang format**:
  - `#!/usr/bin/env -S colored-md`
  - `#!/usr/bin/env colored-md`
  - `#!/usr/bin/colored-md`
  - `#!/path/to/colored-md`

## Example Usage

```bash
# Make this file executable
chmod +x executable.md

# Run it directly
./executable.md

# Or process it normally
colored-md executable.md
```

Both methods produce the same clean output without the shebang line.

## Benefits

This feature enables you to:

- Create standalone markdown documentation that renders itself
- Build self-documenting scripts that display formatted help text
- Distribute markdown files that work as both documents and executables
- Follow UNIX conventions where interpreter directives are hidden from output
