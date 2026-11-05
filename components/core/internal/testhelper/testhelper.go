package testhelper

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func BuildBinary(name, pathprefix string) (string, error) {
	dir := filepath.Join(pathprefix, name)
	aout := filepath.Join(dir, name)

	err := exec.Command("go", "build", "-o", aout, dir).Run()
	if err != nil {
		return "", fmt.Errorf("build command: %w", err)
	}

	return aout, nil
}
