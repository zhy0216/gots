---
title: "Installation"
description: "How to install and set up goTS"
order: 2
category: "Getting Started"
---

# Installation

## Prerequisites

- Go 1.20 or higher

## Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/gots.git
cd gots

# Build the compiler
go build -o gots ./cmd/gots

# Optional: Add to PATH
sudo mv gots /usr/local/bin/
```

## Verify Installation

```bash
gots --version
```
