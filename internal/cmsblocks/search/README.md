# Search Block

A CMS block that provides search functionality across CMS pages and blog posts. Features a search box, results list with type badges, and pagination.

## Features

### 🔍 **Search Capabilities**
- **Pages search**: Searches CMS page titles and content
- **Blog posts search**: Searches blog post titles and content
- **Configurable**: Toggle pages/posts search on/off
- **Real-time**: Works with URL query parameters (`?q=searchterm`)

### 📄 **Results Display**
- **Type badges**: Visual distinction between Pages and Blog Posts
- **Result summary**: Shows excerpt from matched content
- **Clickable links**: Navigate directly to the content
- **Date display**: Shows publication date for blog posts

### 📑 **Pagination**
- **Page numbers**: Navigate through multiple result pages
- **Previous/Next**: Quick navigation buttons
- **URL-based**: Page state preserved in URL (`?page=2`)

## Usage

### Basic Setup

```go
// Create Search block with both stores
searchBlock := search.NewSearchBlockType(cmsStore, blogStore)

// Register the block type
cmsstore.RegisterCustomBlockType(searchBlock)
```

### Configuration Example

```go
// Configure with meta fields
block.SetMeta("placeholder", "Search our site...")
block.SetMeta("results_per_page", "10")
block.SetMeta("show_pages", "true")
block.SetMeta("show_posts", "true")
```

## Admin Configuration

| Meta Key | Description | Default |
|----------|-------------|---------|
| `placeholder` | Search input placeholder text | "Search..." |
| `results_per_page` | Number of results to display | "10" |
| `show_pages` | Include CMS pages in search | "true" |
| `show_posts` | Include blog posts in search | "true" |

## URL Parameters

The search block reads from and writes to URL parameters:

- `?q=searchterm` - The search query
- `?page=2` - Pagination page number (1-based)

Example: `/search?q=press+release&page=2`

## Search Behavior

### What Gets Searched

**CMS Pages:**
- Page title
- Page content (HTML stripped)
- Page alias
- Meta description

**Blog Posts:**
- Post title
- Post content (HTML stripped)
- Post summary/excerpt
- Post slug

### Search Algorithm
- Case-insensitive matching
- Partial word matching supported
- Results combined from pages and posts
- Duplicate filtering (same content not shown twice)

## Generated HTML

```html
<section id="SectionSearch" style="background:#fff; padding:50px 0px 80px 0px;">
  <div class="container">
    <!-- Search Box -->
    <div style="text-align: center;">
      <form method="GET" style="display: flex; max-width: 800px; margin: 0 auto;">
        <input type="search" name="q" class="form-control form-control-lg" placeholder="Search...">
        <button type="submit" class="btn btn-lg" style="background:#794FC6;">Search</button>
      </form>
    </div>

    <!-- Results -->
    <div style="max-width: 900px; margin: 0 auto;">
      <div style="margin-bottom: 30px;">Found 5 result(s) for "searchterm"</div>

      <!-- Result Item -->
      <div style="padding: 25px; border-bottom: 1px solid #e9ecef;">
        <span style="background: #794FC6; color: #fff; padding: 3px 10px; border-radius: 4px;">Blog Post</span>
        <a href="/blog/123/post-slug" style="color: #794FC6; font-size: 22px;">Post Title</a>
        <p style="font-size: 14px; color: #6c757d;">12 Jan, 2024</p>
        <p style="font-size: 16px;">Summary text from the post...</p>
      </div>
    </div>

    <!-- Pagination -->
    <nav style="margin-top: 30px;">
      <ul class="pagination justify-content-center">
        <li class="page-item"><a class="page-link" href="?q=searchterm">&laquo; Previous</a></li>
        <li class="page-item active"><span class="page-link">1</span></li>
        <li class="page-item"><a class="page-link" href="?q=searchterm&page=2">2</a></li>
        <li class="page-item"><a class="page-link" href="?q=searchterm&page=2">Next &raquo;</a></li>
      </ul>
    </nav>
  </div>
</section>
```

## Styling

The search block uses inline styles and Bootstrap classes:

### Key CSS Classes
- `.form-control` - Search input styling
- `.btn` - Search button styling
- `.pagination` - Pagination controls
- `.page-item` / `.page-link` - Pagination items

### Custom Colors
- Primary purple: `#794FC6` (title links, badges)
- Teal accent: `#1ba1b6` (page badges)
- Gray text: `#6c757d` (meta information)

## Integration

### With CMS Page

1. Create a CMS page (e.g., `/search`)
2. Add the Search block to the page
3. Configure search options in block settings
4. The search form will submit to the same page

### With Navigation

Add a search link to your navbar:

```go
searchLink := hb.A().
    Class("nav-link").
    Href("/search").
    Child(hb.I().Class("bi bi-search"))
```

## Files

- `search_block_type.go` - Main block implementation
- `renderer.go` - HTML rendering logic
- `example_usage.go` - Usage examples
- `search_block_type_test.go` - Unit tests
- `README.md` - This documentation

## Dependencies

- `github.com/dracory/cmsstore` - Core CMS interfaces
- `github.com/dracory/blogstore` - Blog post search
- `github.com/dracory/form` - Admin form fields
- `github.com/dracory/hb` - HTML builder library
- `github.com/dracory/bs` - Bootstrap HTML components

## Registration

Add this to your CMS block registration (already done in `cms_add_block_types.go`):

```go
import "project/internal/cms/blocks/search"

func CmsAddBlockTypes(registry RegistryInterface) {
    // ... other blocks ...

    // Register the Search block
    searchBlock := search.NewSearchBlockType(
        registry.GetCmsStore(),
        registry.GetBlogStore(),
    )
    cmsstore.RegisterCustomBlockType(searchBlock)
}
```

## Testing

Run the tests:

```bash
go test ./internal/cms/blocks/search/...
```
