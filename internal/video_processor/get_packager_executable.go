package videoprocessor

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed embed/packager-win-x64.exe
var packagerFS []byte

func GetPackagerExecutable() (string, error) {
	packagerName := "packager-win-x64.exe"

	// Extract the embedded packager to a temporary location
	tempPath := filepath.Join(os.TempDir(), packagerName)
	if err := os.WriteFile(tempPath, packagerFS, 0755); err != nil {
		return "", err
	}
	return tempPath, nil
}
