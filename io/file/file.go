package file

import (
	"fmt"
	"io"
	"os"
)

// Rename renames the file from src to dst. If src and dst can't be renamed
// regularly, the data is copied from src to dst. dst will be overwritten
// if it already exists. src will be removed after all data has been copied
// successfully. Both files exist during copying.
func Rename(src, dst string) error {
	// First try to rename the file
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// If renaming the file fails, copy the data
	Copy(src, dst)

	if err := os.Remove(src); err != nil {
		os.Remove(dst)
		return fmt.Errorf("failed to remove source file: %w", err)
	}

	return nil
}

// Copy copies a file from src to dst.
func Copy(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}

	destination, err := os.Create(dst)
	if err != nil {
		source.Close()
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		source.Close()
		os.Remove(dst)
		return fmt.Errorf("failed to copy data from source to destination: %w", err)
	}

	source.Close()

	return nil
}
