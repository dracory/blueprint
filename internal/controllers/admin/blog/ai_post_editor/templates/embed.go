package templates

import (
	"embed"

	"github.com/flosch/pongo2/v6"
)

//go:embed *
var files embed.FS

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

func Template(path string, data map[string]any) (string, error) {
	tmp, err := ToString(path)
	if err != nil {
		return "", err
	}

	tmpl, err := pongo2.FromString(tmp)
	if err != nil {
		panic(err)
	}

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
