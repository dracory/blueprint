# Proposal: Email Styling Constants Migration

## Overview

This proposal outlines the migration of email styling constants from the Blueprint project to the shared base email package. This will ensure consistent email styling across all Dracory projects while maintaining flexibility for customization.

## Current State Analysis

### Blueprint Email Styling

The current Blueprint project contains minimal styling constants in `internal/emails/consts.go`:

```go
const STYLE_HEADING = "margin:0px;padding:10px 0px;text-align:left;font-size:22px;"
const STYLE_PARAGRAPH = "margin:0px;padding:10px 0px;text-align:left;font-size:16px;"
const STYLE_BUTTON = "display: inline-block; padding: 10px 20px; font-size: 16px; color: white; background-color: #007BFF; text-align: center; text-decoration: none; border-radius: 5px;"
```

### Current Usage Pattern

These constants are used throughout Blueprint's email templates:

```go
// Example from user_email_invite_friend.go
h1 := hb.Heading1().
    HTML(`You have an awesome friend`).
    Style(STYLE_HEADING)

p1 := hb.Paragraph().
    HTML(`Hi ` + recipientName + `,`).
    Style(STYLE_PARAGRAPH)
```

### Existing Base Email Package

The base email package already exists at `github.com/dracory/base/email` with:
- Email sending functionality
- Template generation with responsive design
- Comprehensive styling in the template itself

## Identified Issues

1. **Inconsistent Styling**: Each project may define different email styles
2. **Limited Scope**: Only 3 basic styles defined in Blueprint
3. **Maintenance Overhead**: Styling changes require updates across multiple projects
4. **Design Inconsistency**: Different projects may have different email appearances
5. **Missing Styles**: No comprehensive style library for common email components

## Proposed Solution

### 1. Enhanced Base Email Styling Package

Create a comprehensive styling system within the existing base email package:

```go
// github.com/dracory/base/email/styles.go
package email

// Typography Styles
const (
    StyleHeading1   = "margin:0px;padding:10px 0px;text-align:left;font-size:24px;font-weight:600;color:#333333;"
    StyleHeading2   = "margin:0px;padding:8px 0px;text-align:left;font-size:20px;font-weight:600;color:#333333;"
    StyleHeading3   = "margin:0px;padding:6px 0px;text-align:left;font-size:18px;font-weight:600;color:#333333;"
    StyleParagraph = "margin:0px;padding:10px 0px;text-align:left;font-size:16px;line-height:1.6;color:#333333;"
    StyleSmall      = "margin:0px;padding:5px 0px;text-align:left;font-size:14px;color:#666666;"
)

// Button Styles
const (
    StyleButtonPrimary   = "display: inline-block; padding: 12px 24px; font-size: 16px; font-weight:600; color: white; background-color: #007BFF; text-align: center; text-decoration: none; border-radius: 6px; border: 1px solid #007BFF;"
    StyleButtonSecondary = "display: inline-block; padding: 12px 24px; font-size: 16px; font-weight:600; color: #007BFF; background-color: transparent; text-align: center; text-decoration: none; border-radius: 6px; border: 2px solid #007BFF;"
    StyleButtonSuccess   = "display: inline-block; padding: 12px 24px; font-size: 16px; font-weight:600; color: white; background-color: #28A745; text-align: center; text-decoration: none; border-radius: 6px; border: 1px solid #28A745;"
    StyleButtonDanger    = "display: inline-block; padding: 12px 24px; font-size: 16px; font-weight:600; color: white; background-color: #DC3545; text-align: center; text-decoration: none; border-radius: 6px; border: 1px solid #DC3545;"
    StyleButtonSmall     = "display: inline-block; padding: 8px 16px; font-size: 14px; font-weight:600; color: white; background-color: #007BFF; text-align: center; text-decoration: none; border-radius: 4px; border: 1px solid #007BFF;"
)

// Layout Styles
const (
    StyleContainer   = "max-width: 600px; margin: 0 auto; padding: 20px; background-color: #ffffff;"
    StyleSection     = "margin: 20px 0px; padding: 15px; background-color: #f8f9fa; border-radius: 6px;"
    StyleDivider     = "height: 1px; background-color: #dee2e6; margin: 20px 0px; border: none;"
    StyleCard        = "padding: 20px; background-color: #ffffff; border: 1px solid #dee2e6; border-radius: 8px; margin: 10px 0px;"
)

// Alert Styles
const (
    StyleAlertInfo    = "padding: 12px 16px; background-color: #D1ECF1; border: 1px solid #BEE5EB; border-radius: 4px; color: #0C5460; margin: 10px 0px;"
    StyleAlertSuccess = "padding: 12px 16px; background-color: #D4EDDA; border: 1px solid #C3E6CB; border-radius: 4px; color: #155724; margin: 10px 0px;"
    StyleAlertWarning = "padding: 12px 16px; background-color: #FFF3CD; border: 1px solid #FFEAA7; border-radius: 4px; color: #856404; margin: 10px 0px;"
    StyleAlertDanger  = "padding: 12px 16px; background-color: #F8D7DA; border: 1px solid #F5C6CB; border-radius: 4px; color: #721C24; margin: 10px 0px;"
)

// List Styles
const (
    StyleListUnordered = "margin: 10px 0px; padding-left: 20px; color: #333333;"
    StyleListOrdered   = "margin: 10px 0px; padding-left: 20px; color: #333333;"
    StyleListItem      = "margin: 5px 0px; line-height: 1.5;"
)

// Table Styles
const (
    StyleTable     = "width: 100%; border-collapse: collapse; margin: 15px 0px;"
    StyleTableHead = "background-color: #f8f9fa; border: 1px solid #dee2e6; padding: 12px; text-align: left; font-weight: 600;"
    StyleTableCell = "border: 1px solid #dee2e6; padding: 12px; text-align: left;"
)

// Utility Styles
const (
    StyleTextCenter   = "text-align: center;"
    StyleTextRight    = "text-align: right;"
    StyleTextMuted    = "color: #6c757d;"
    StyleTextPrimary  = "color: #007BFF;"
    StyleTextSuccess  = "color: #28A745;"
    StyleTextDanger   = "color: #DC3545;"
    StyleTextWarning  = "color: #FFC107;"
    StyleBgLight      = "background-color: #f8f9fa;"
    StyleBgDark       = "background-color: #343a40;"
)
```

### 2. Style Builder Functions

Provide helper functions for dynamic style creation:

```go
// Style builders for customization
func ButtonStyle(color, backgroundColor string, size ButtonSize) string
func TextStyle(color, size string, weight FontWeight) string
func AlertStyle(alertType AlertType) string
func ContainerStyle(padding, margin string) string

// Enums for type safety
type ButtonSize string
const (
    ButtonSizeSmall  ButtonSize = "small"
    ButtonSizeMedium ButtonSize = "medium"
    ButtonSizeLarge  ButtonSize = "large"
)

type AlertType string
const (
    AlertTypeInfo    AlertType = "info"
    AlertTypeSuccess AlertType = "success"
    AlertTypeWarning AlertType = "warning"
    AlertTypeDanger  AlertType = "danger"
)
```

### 3. Theme Support

Add theme support for different brand identities:

```go
// Theme definitions
type Theme struct {
    Primary   string
    Secondary string
    Success   string
    Danger    string
    Warning   string
    Light     string
    Dark      string
}

// Predefined themes
var (
    ThemeDefault = Theme{
        Primary:   "#007BFF",
        Secondary: "#6C757D",
        Success:   "#28A745",
        Danger:    "#DC3545",
        Warning:   "#FFC107",
        Light:     "#F8F9FA",
        Dark:      "#343A40",
    }
    
    ThemeDracory = Theme{
        Primary:   "#17A2B8",
        Secondary: "#6C757D",
        Success:   "#28A745",
        Danger:    "#DC3545",
        Warning:   "#FFC107",
        Light:     "#F8F9FA",
        Dark:      "#343A40",
    }
)

// Theme-based style generation
func GetTheme(theme Theme) ThemeStyles
type ThemeStyles struct {
    ButtonPrimary   string
    ButtonSecondary string
    TextPrimary     string
    // ... other themed styles
}
```

## Migration Plan

### Phase 1: Base Package Enhancement

1. **Create Styles Module**:
   ```
   github.com/dracory/base/email/
   ├── styles.go          # All style constants
   ├── builders.go        # Style builder functions
   ├── theme.go           # Theme support
   └── components.go      # Pre-built component styles
   ```

2. **Implement Core Styles**:
   - Typography styles (headings, paragraphs, text)
   - Button styles (primary, secondary, success, danger)
   - Layout styles (containers, sections, dividers)
   - Alert styles (info, success, warning, danger)
   - List and table styles
   - Utility styles (alignment, colors, backgrounds)

3. **Add Style Builders**:
   - Dynamic style creation functions
   - Type-safe enums for sizes and types
   - Theme-based style generation

### Phase 2: Blueprint Integration

1. **Update Blueprint Imports**:
   ```go
   // Before
   import "project/internal/emails"
   
   // After
   import baseEmail "github.com/dracory/base/email"
   ```

2. **Replace Style Constants**:
   ```go
   // Before
   h1.Style(emails.STYLE_HEADING)
   p1.Style(emails.STYLE_PARAGRAPH)
   
   // After
   h1.Style(baseEmail.StyleHeading1)
   p1.Style(baseEmail.StyleParagraph)
   ```

3. **Remove Deprecated Constants**:
   - Delete `internal/emails/consts.go`
   - Update all email template files

### Phase 3: Enhanced Email Components

1. **Pre-built Components**:
   ```go
   // Helper functions for common email patterns
   func HeroSection(title, subtitle string) string
   func CallToAction(text, url string) string
   func InfoBox(message string) string
   func Footer(appName string) string
   ```

2. **Template Integration**:
   - Integrate styles with existing template system
   - Provide style-customizable template options
   - Add responsive design considerations

## Benefits

### 1. Design Consistency
- **Unified Brand Identity**: All Dracory projects use consistent email styling
- **Professional Appearance**: Comprehensive style library ensures polished emails
- **Responsive Design**: Mobile-optimized styles built-in

### 2. Developer Experience
- **Faster Development**: Pre-built styles and components speed up email creation
- **Type Safety**: Builder functions and enums prevent style errors
- **Documentation**: Clear style definitions and usage examples

### 3. Maintainability
- **Single Source of Truth**: All styling centralized in base package
- **Easy Updates**: Style changes benefit all projects automatically
- **Version Control**: Style evolution tracked in one location

### 4. Flexibility
- **Theme Support**: Projects can customize colors while maintaining structure
- **Extensible**: Easy to add new styles and components
- **Backward Compatible**: Existing email templates can be updated gradually

## Implementation Details

### Style Naming Convention

Follow consistent naming patterns:
- `Style[Component][Variant]` (e.g., `StyleButtonPrimary`, `StyleAlertSuccess`)
- `Style[Element][Size]` (e.g., `StyleHeading1`, `StyleButtonSmall`)
- `Style[Property][Value]` (e.g., `StyleTextCenter`, `StyleBgLight`)

### CSS Best Practices

- **Responsive Design**: Mobile-first approach with media queries
- **Email Client Compatibility**: Use inline styles for maximum compatibility
- **Accessibility**: Proper contrast ratios and semantic styling
- **Progressive Enhancement**: Graceful degradation in older email clients

### Performance Considerations

- **Minimal CSS**: Use only necessary styles to reduce email size
- **Compression**: Optimize styles for faster loading
- **Caching**: Template-level style caching for improved performance

## File Structure

### Before (Blueprint)
```
internal/emails/
├── consts.go                    (6 lines - 3 basic styles)
├── [7 email template files using constants]
└── [repetitive style definitions]
```

### After (Base Package)
```
github.com/dracory/base/email/
├── styles.go                    (150+ lines - comprehensive styles)
├── builders.go                  (80+ lines - style builders)
├── theme.go                     (60+ lines - theme support)
├── components.go                (100+ lines - pre-built components)
├── send.go                      (existing)
├── template.go                  (existing)
└── README.md                    (updated documentation)

# Blueprint (simplified)
internal/emails/
├── [7 email template files using base styles]
└── [no style constants - uses base package]
```

## Usage Examples

### Basic Usage
```go
import baseEmail "github.com/dracory/base/email"

// Use predefined styles
h1 := hb.Heading1().
    HTML("Welcome to Our Service").
    Style(baseEmail.StyleHeading1)

button := hb.Hyperlink().
    Text("Get Started").
    Href("https://example.com").
    Style(baseEmail.StyleButtonPrimary)
```

### Theme-Based Usage
```go
// Apply custom theme
theme := baseEmail.ThemeDracory
styles := baseEmail.GetTheme(theme)

// Use themed styles
button := hb.Hyperlink().
    Text("Get Started").
    Href("https://example.com").
    Style(styles.ButtonPrimary)
```

### Dynamic Style Building
```go
// Create custom button style
customButton := baseEmail.ButtonStyle("#FFFFFF", "#17A2B8", baseEmail.ButtonSizeMedium)

// Create custom alert
alert := baseEmail.AlertStyle(baseEmail.AlertTypeInfo)
```

### Component Usage
```go
// Use pre-built components
hero := baseEmail.HeroSection("Welcome!", "Get started with our amazing service")
cta := baseEmail.CallToAction("Sign Up Now", "https://example.com/signup")
```

## Success Metrics

1. **Adoption Rate**: All new Dracory projects use base email styles within 3 months
2. **Consistency Score**: 90%+ visual consistency across project emails
3. **Development Speed**: 50% reduction in email template creation time
4. **Maintenance Reduction**: 80% fewer style-related bug reports

## Timeline

- **Phase 1**: 1-2 weeks (Base package enhancement)
- **Phase 2**: 1 week (Blueprint integration)
- **Phase 3**: 1-2 weeks (Enhanced components)

**Total Estimated Time**: 3-5 weeks

## Risks and Mitigations

### Risk: Email Client Compatibility
- **Mitigation**: Extensive testing across major email clients
- **Fallback**: Provide fallback styles for older clients

### Risk: Breaking Changes
- **Mitigation**: Maintain backward compatibility during transition
- **Migration Path**: Provide clear upgrade guide and examples

### Risk: Style Bloat
- **Mitigation**: Implement style tree-shaking and modular imports
- **Optimization**: Regular performance audits and optimization

## Conclusion

This migration will establish a comprehensive, consistent email styling system for all Dracory projects. The enhanced base email package will provide professional, responsive email templates while maintaining flexibility for customization.

The proposed solution balances consistency with flexibility, ensuring that all Dracory projects benefit from professional email design while retaining the ability to express their unique brand identity through themes and custom styles.

By centralizing email styling in the base package, we reduce maintenance overhead, improve developer experience, and ensure a cohesive brand presence across all Dracory communications.
