# Upgrade Guide: v0.51.0 to v0.52.0

This guide helps LLMs and developers upgrade Blueprint applications from v0.51.0 to v0.52.0.

## Overview

This is a small maintenance release. All generated short IDs now use `uid.SafeShortID()` from `github.com/dracory/uid` v1.12.0 — a 14-character, vowel-free, chronologically sortable ID — instead of `uid.ShortID()` / `GenerateShortID()`. The vowel-free alphabet avoids generating accidental English words in user-facing IDs (flash message keys, DOM element IDs, ORM record IDs). Several dependencies were also bumped.

**Key Changes:**
- `internal/helpers/flash.go` — `uid.ShortID()` → `uid.SafeShortID()`
- `internal/controllers/user/partials/page_header.go` — `uid.ShortID()` → `uid.SafeShortID()`
- `database/models/note.go` — `uid.GenerateShortID()` (from `github.com/dracory/neat/support/uid`) → `uid.SafeShortID()` (from `github.com/dracory/uid`)
- Dependency updates: `blogstore` v1.35.2→v1.36.0, `llm` v1.8.1→v1.9.0, `shopstore` v1.26.1→v1.26.2, `statsstore` v1.31.1→v1.31.3, `uid` v1.11.0→v1.12.0, `google.golang.org/api` v0.301.0→v0.304.0, `google.golang.org/genai` v1.73.0→v1.75.0

---

## ⚠️ Breaking Changes

### 1. Short IDs are 14 characters instead of 11

**Change**: `uid.SafeShortID()` returns a 14-character lowercase ID (`"05fr88tc139r8t"`) encoded in a vowel-free Base 16 alphabet, versus the 11-character Crockford Base32 output of `uid.ShortID()` / `GenerateShortID()` (`"sa4rc789wxg"`).

**Old Usage**:
```go
id := uid.ShortID()           // 11 chars, Crockford Base32
id := uid.GenerateShortID()   // neat/support/uid equivalent
```

**New Usage**:
```go
id := uid.SafeShortID()       // 14 chars, vowel-free Base 16
```

**Action Required**:
- If any column or field storing generated IDs is sized for exactly 11 characters, widen it to at least 14 (`VARCHAR(32)` or larger recommended). Blueprint's own usages (flash message cache keys, DOM IDs, the example `note` model's `orm.ShortID`) are unconstrained.
- IDs remain lowercase, URL-safe, and chronologically sortable — no format validation changes needed unless you validate `uid.IsShortID()` length explicitly.

---

## 🔄 Migration Steps

### Step 1: Update the version constant

Update `internal/config/version.go`:

```go
const Version = "0.52.0"
```

### Step 2: Update dependencies

```bash
go get -u ./...
go mod tidy
```

Or pin the dracory bumps explicitly:

```bash
go get github.com/dracory/blogstore@v1.36.0 \
       github.com/dracory/llm@v1.9.0 \
       github.com/dracory/shopstore@v1.26.2 \
       github.com/dracory/statsstore@v1.31.3 \
       github.com/dracory/uid@v1.12.0
go mod tidy
```

### Step 3: Replace ID generators with `uid.SafeShortID()`

Search your project for every place a short ID is generated and swap the call:

```go
// Before
import "github.com/dracory/uid"              // or "github.com/dracory/neat/support/uid"
id := uid.ShortID()                          // dracory/uid
id := uid.GenerateShortID()                  // neat/support/uid

// After
import "github.com/dracory/uid"
id := uid.SafeShortID()
```

Apply to:

- `internal/helpers/flash.go` — flash `message_id`
- `internal/controllers/user/partials/page_header.go` — DOM element ID
- `database/models/note.go` — `NewNote()` (also switch the import from `github.com/dracory/neat/support/uid` to `github.com/dracory/uid`)
- Any project-local code calling `uid.ShortID()`, `uid.GenerateShortID()`, or `neat.GenerateID()` for user-facing IDs

### Step 4: Verify

```bash
go mod tidy
go build ./...
go test ./...
```

---

## 🧪 Testing After Migration

1. **Unit tests**: `go test ./internal/helpers/... ./internal/controllers/user/... ./database/...`
2. **Flash messages**: trigger any flash redirect (e.g. a failed form post) and confirm the `message_id` in the URL is a 14-char lowercase ID and the message renders.
3. **Existing rows**: IDs generated before the upgrade remain valid — `SafeShortID` only changes newly generated values; no data migration is needed.

---

## 📝 Additional Notes

- **Why vowel-free?** `SafeShortID` maps every hex letter to a non-vowel (`a`→`b` … `e`→`r`, `f`→`t`), greatly reducing accidental English words in IDs shown to users or embedded in URLs.
- **Per-process uniqueness** is unchanged: IDs pack `(microsecond timestamp << 4) | counter`, so they sort chronologically but are not unguessable — never use them as security tokens.
- **`neat/support/uid` still exists** and its `GenerateShortID()` still works; blueprint simply prefers `dracory/uid.SafeShortID()` for all new IDs.

---

## 🆘 Common Issues and Solutions

### Issue: Compile error `uid.SafeShortID undefined`

**Cause**: `github.com/dracory/uid` is older than v1.12.0, or the file still imports `github.com/dracory/neat/support/uid` (which has no `SafeShortID`).

**Solution**: `go get github.com/dracory/uid@v1.12.0` and fix the import path.

### Issue: Generated IDs rejected by a strict 11-char validator or DB column

**Cause**: `SafeShortID` produces 14 characters.

**Solution**: Widen the column/validator to 14+ characters, or keep `uid.ShortID()` for that specific field.

---

## 📞 Support

- Repository: https://github.com/dracory/blueprint
- For questions, open an issue on the repository.
