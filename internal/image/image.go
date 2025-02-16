package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CreatQeFatImage creates a FAT disk image containing the provided files.
// files: map of filenames to their contents
// outputPath: where to save the final disk image
// volumeName: name for the disk volume (optional, defaults to "DISK")
// sizeInMB: size of the disk image in MB (optional, defaults to 32)
func CreateFatImage(files map[string][]byte, outputPath string, volumeName string, sizeInMB int) error {
	if volumeName == "" {
		volumeName = "DISK"
	}
	if sizeInMB <= 0 {
		sizeInMB = 32
	}

	// Create a temporary directory for the files
	tempDir, err := os.MkdirTemp("", "files-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Write all files to the temporary directory
	for name, content := range files {
		filePath := filepath.Join(tempDir, name)
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", name, err)
		}
	}

	// Use hdiutil to create a FAT disk image
	cmd := exec.Command("hdiutil", "create",
		"-size", fmt.Sprintf("%dm", sizeInMB),
		"-fs", "MS-DOS",
		"-volname", volumeName,
		"-format", "UDTO",
		"-srcfolder", tempDir,
		outputPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("hdiutil failed: %w\nOutput: %s", err, output)
	}

	// hdiutil adds .cdr to the output file, we need to rename it
	if err := os.Rename(outputPath+".cdr", outputPath); err != nil {
		return fmt.Errorf("failed to rename image: %w", err)
	}

	return nil
}
