# Documentation System

The goTS website now uses Astro's content collections to manage documentation. All markdown files in `src/content/docs/` are automatically compiled into `/docs/*` routes.

## File Structure

```
website/
├── src/
│   ├── content/
│   │   ├── config.ts              # Content collection schema
│   │   └── docs/                  # Your markdown files go here
│   │       ├── introduction.md
│   │       ├── installation.md
│   │       ├── cli-commands.md
│   │       ├── type-system.md
│   │       └── functions-closures.md
│   ├── pages/
│   │   └── docs/
│   │       ├── index.astro        # Redirects to first doc
│   │       └── [...slug].astro    # Dynamic route handler
│   └── components/
│       └── docs/
│           └── DocsSidebar.astro  # Auto-generated sidebar
```

## Creating New Documentation

### 1. Create a new markdown file

Create a file in `src/content/docs/` with the `.md` extension:

```bash
touch src/content/docs/my-new-page.md
```

### 2. Add frontmatter

Every doc file needs frontmatter with metadata:

```markdown
---
title: "My Page Title"
description: "A brief description of this page"
order: 10
category: "Getting Started"
---

# My Page Title

Your content here...
```

### Frontmatter Fields

- **title** (required): Page title shown in sidebar and page header
- **description** (optional): Meta description for SEO
- **order** (optional): Number for ordering within a category (lower = earlier)
- **category** (optional): Group pages in sidebar sections (e.g., "Getting Started", "Core Concepts")

### 3. Write content

Use standard markdown syntax:

```markdown
## Headings

### Subheadings

Regular paragraphs with **bold** and *italic* text.

- Bullet lists
- Work great

Code blocks with syntax highlighting:

\`\`\`typescript
let x: int = 42
println(x)
\`\`\`

Inline `code` works too.
```

## Routes

Files automatically map to routes:

- `src/content/docs/introduction.md` → `/docs/introduction`
- `src/content/docs/type-system.md` → `/docs/type-system`
- `src/content/docs/advanced/closures.md` → `/docs/advanced/closures`

## Sidebar

The sidebar is automatically generated from your markdown files:

1. Pages are grouped by `category`
2. Within each category, pages are sorted by `order`
3. The current page is highlighted
4. Categories appear in the order they're encountered

## Development

```bash
# Start dev server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Styling

The documentation uses prose styling from the main layout. Code blocks get syntax highlighting from `@astrojs/prism`.

Custom styling is in `/src/pages/docs/[...slug].astro` in the `<style>` section.
