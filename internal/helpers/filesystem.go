package helpers

import "embed"

// EmbeddedFileToBytes reads a file from the embedded filesystem and returns its content as bytes
func EmbeddedFileToBytes(embeddedFileSystem embed.FS, path string) ([]byte, error) {
	return embeddedFileSystem.ReadFile(path)
}

// EmbeddedFileToString reads a file from the embedded filesystem and returns its content as a string
func EmbeddedFileToString(embeddedFileSystem embed.FS, path string) (string, error) {
	bytes, err := EmbeddedFileToBytes(embeddedFileSystem, path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
