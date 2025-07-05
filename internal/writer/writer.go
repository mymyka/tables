package writer

import (
	"os"
	"path/filepath"
)

func Write(root string, c map[string]string) error {
	// Define destination directory
	dest := "."

	// Create base destination directory if it doesn't exist
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	// Iterate through the map and create files
	for filename, content := range c {
		// Create directory structure: dest/root/filename/
		dirPath := filepath.Join(dest, root, filename)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}

		// Create full path: dest/root/filename/filename.go
		fullPath := filepath.Join(dirPath, filename+".go")

		// Create or overwrite file
		file, err := os.Create(fullPath)
		if err != nil {
			return err
		}

		// Write content to file
		_, err = file.WriteString(content)
		file.Close() // Close immediately after writing

		if err != nil {
			return err
		}
	}

	return nil
}
