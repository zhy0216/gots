---
title: "CLI Commands"
description: "Command line interface reference"
order: 3
category: "Getting Started"
---

# CLI Commands

The goTS CLI provides several commands for building, running, and working with goTS programs.

## Commands

### `gots run`

Compile and execute a goTS program immediately.

```bash
gots run program.gts
```

### `gots build`

Compile to a native binary with optional custom output name.

```bash
gots build program.gts
gots build program.gts -o myapp
```

### `gots emit-go`

Generate Go source code without compiling to binary.

```bash
gots emit-go program.gts
gots emit-go program.gts --output=output.go
```

### `gots repl`

Start an interactive REPL for testing and experimentation.

```bash
gots repl
```
