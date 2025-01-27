package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileDiff represents the differences between two directories
type FileDiff struct {
	Added   []string
	Removed []string
	Changed []string
}

// ComparePaths compares two directory paths and returns lists of added, removed, and changed files
func ComparePaths(beforePath, afterPath string, ignoreFiles []string) (*FileDiff, error) {
	// Get file maps for both directories
	beforeFiles, err := getFileMap(beforePath, ignoreFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to read before directory: %w", err)
	}

	afterFiles, err := getFileMap(afterPath, ignoreFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to read after directory: %w", err)
	}

	diff := &FileDiff{
		Added:   make([]string, 0),
		Removed: make([]string, 0),
		Changed: make([]string, 0),
	}

	// Find added and changed files
	for path, afterHash := range afterFiles {
		beforeHash, exists := beforeFiles[path]
		if !exists {
			// File was added
			diff.Added = append(diff.Added, path)
		} else if beforeHash != afterHash {
			// File was changed
			diff.Changed = append(diff.Changed, path)
		}
	}

	// Find removed files
	for path := range beforeFiles {
		if _, exists := afterFiles[path]; !exists {
			diff.Removed = append(diff.Removed, path)
		}
	}

	return diff, nil
}

// getFileMap returns a map of relative file paths to their SHA256 hashes
func getFileMap(root string, ignoreFiles []string) (map[string]string, error) {
	files := make(map[string]string)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, err)
		}

		// Skip ignored files
		if shouldIgnore(info.Name(), ignoreFiles) {
			return nil
		}

		// Calculate file hash
		hash, err := hashFile(path)
		if err != nil {
			return fmt.Errorf("failed to hash file %s: %w", path, err)
		}

		files[relPath] = hash
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// shouldIgnore returns true if the filename matches any of the ignore patterns
func shouldIgnore(filename string, ignoreFiles []string) bool {
	for _, ignore := range ignoreFiles {
		if filename == ignore {
			return true
		}
	}
	return false
}

// hashFile calculates the SHA256 hash of a file
func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// HasChanges returns true if there are any differences between the directories
func (d *FileDiff) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
} 
