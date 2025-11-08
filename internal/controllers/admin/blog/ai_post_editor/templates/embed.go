package templates

import (
	"embed"
	"sync"

	"github.com/flosch/pongo2/v6"
)

//go:embed *
var files embed.FS

// Caches for parsed templates to avoid reparsing on every request.
var (
	templateCache = make(map[string]*pongo2.Template)
	cacheMutex    = &sync.RWMutex{}
)

func ToBytes(path string) ([]byte, error) {
	return files.ReadFile(path)
}
func ToString(path string) (string, error) {
	bytes, err := ToBytes(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ResourceExists(path string) bool {
	_, err := files.ReadFile(path)
	return err == nil
}

// Template renders a template from the embedded files, using a cache to avoid
// reparsing the template on every call.
//
// Example:
//
//	data := map[string]any{
//		"Title": "Hello World",
//	}
//
//	template, err := templates.Template("index.html", data)
//	if err != nil {
//		panic(err)
//	}
func Template(path string, data map[string]any) (string, error) {
	cacheMutex.RLock()
	tmpl, found := templateCache[path]
	cacheMutex.RUnlock()

	if !found {
		cacheMutex.Lock()
		// Double-check if another goroutine populated the cache while we were waiting for the lock.
		tmpl, found = templateCache[path]
		if !found {
			s, err := ToString(path)
			if err != nil {
				cacheMutex.Unlock()
				return "", err
			}
			parsedTmpl, err := pongo2.FromString(s)
			if err != nil {
				cacheMutex.Unlock()
				return "", err // Replaced panic with proper error return.
			}
			templateCache[path] = parsedTmpl
			tmpl = parsedTmpl
		}
		cacheMutex.Unlock()
	}

	// Execute the template (either from cache or newly parsed).
	out, err := tmpl.Execute(data)
	if err != nil {
		return "", err
	}

	return out, nil
}

// Tpl is a shortcut for Template that ignores errors
func Tpl(path string, data map[string]any) string {
	out, err := Template(path, data)
	if err != nil {
		return ""
	}

	return out
}
