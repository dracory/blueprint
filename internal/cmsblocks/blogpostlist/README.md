# Blog Post List Block

A CMS block that renders a grid of blog posts with configurable display options.

## Features

- **Tag Filtering**: Automatically filters posts by tag when URL contains `/tag/{tag-slug}/`
- **Configurable Layout**: Set number of columns (1, 2, 3, 4, 6)
- **Pagination**: Optional pagination controls
- **Display Options**: Toggle images, summaries, dates
- **Excerpt Length**: Control auto-generated excerpt length

## Block Type Key

`blog_post_list`

## Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| posts_per_page | int | 12 | Number of posts to display |
| columns | int | 4 | Grid columns (1, 2, 3, 4, 6) |
| show_pagination | bool | true | Show pagination controls |
| show_images | bool | true | Display featured images |
| show_summary | bool | true | Display post summary/excerpt |
| show_date | bool | true | Display publication date |
| excerpt_length | int | 150 | Max characters for excerpts |

## Tag Filtering

When a page containing this block is accessed via a URL like:

```
https://pressplugs.co.uk/tag/journalists-journorequest-pr-media-press/
```

The block automatically:
1. Extracts the tag slug (`journalists-journorequest-pr-media-press`)
2. Finds the matching tag
3. Displays only posts with that tag

If the tag doesn't exist, an empty list is displayed.

## Usage

```go
import "project/internal/cms/blocks/blogpostlist"

blockType := blogpostlist.NewBlogPostListBlockType(blogStore)
```

## Files

- `blog_post_list_block_type.go` - Block type implementation
- `renderer.go` - HTML rendering logic
- `blog_post_list_block_type_test.go` - Unit tests
