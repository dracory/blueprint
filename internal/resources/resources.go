package resources

import (
	"bytes"
	"embed"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

//go:embed *
var files embed.FS

func ToBytes(path string) ([]byte, error) {
	return files.ReadFile(path)
}

func ToString(path string) (string, error) {
	b, err := ToBytes(path)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func ResourceExists(path string) bool {
	_, err := files.ReadFile(path)
	return err == nil
}

func Resource(path string) (string, error) {
	str, err := files.ReadFile(path)
	if err != nil {
		log.Println("Resource: " + path + " NOT FOUND")
		return "", err
	}
	return string(str), nil
}

func ResourceWithParams(path string, params map[string]string) string {
	parsed := template.Must(template.ParseFS(files, path))
	var tpl bytes.Buffer
	if err := parsed.Execute(&tpl, params); err != nil {
		log.Println(err)
		return ""
	}

	return tpl.String()
}

func ImageToBase64String(path string) string {
	data, _ := files.ReadFile(path)

	// Check if it's an SVG file based on extension
	isSvg := filepath.Ext(path) == ".svg"

	var mimeType string
	if isSvg {
		mimeType = "image/svg+xml"
	} else {
		mimeType = http.DetectContentType(data)
	}

	base64Encoding := ""

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/bmp":
		base64Encoding += "data:image/bmp;base64,"
	case "image/gif":
		base64Encoding += "data:image/gif;base64,"
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	case "image/webp":
		base64Encoding += "data:image/webp;base64,"
	case "image/svg+xml":
		base64Encoding += "data:image/svg+xml;base64,"
	default:
		base64Encoding += "data:image/*;base64,"
	}

	base64Encoding += base64.StdEncoding.EncodeToString(data)
	return base64Encoding
}
