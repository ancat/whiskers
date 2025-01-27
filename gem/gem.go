package gem

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Source represents where a gem can be fetched from
type Source struct {
	Type string // e.g., "git", "rubygems"
	URL  string // e.g., "https://rubygems.org" or git repository URL
}

// Gem represents a Ruby gem with its basic metadata
type Gem struct {
	Name    string
	Version string
	Source  Source
}

// NewGem creates a new Gem instance
func NewGem(name, version string, source Source) *Gem {
	return &Gem{
		Name:    name,
		Version: version,
		Source:  source,
	}
}

// String returns a string representation of the gem
func (g *Gem) String() string {
	return g.Name + " (" + g.Version + ")"
}

// IsFromRubyGems returns true if the gem is from the default RubyGems source
func (g *Gem) IsFromRubyGems() bool {
	return g.Source.URL == "https://rubygems.org/"
}

// GetDownloadURL returns the URL to download the gem
func (g *Gem) GetDownloadURL() string {
	if g.IsFromRubyGems() {
		return fmt.Sprintf("https://rubygems.org/gems/%s-%s.gem", g.Name, g.Version)
	}
	return ""
}

// DownloadAndExtract downloads the gem file and extracts it to the specified directory
func (g *Gem) DownloadAndExtract(baseDir string) error {
	if !g.IsFromRubyGems() {
		return fmt.Errorf("downloading is only supported for RubyGems.org gems")
	}

	// Create the target directory (baseDir/name-version)
	targetDir := filepath.Join(baseDir, fmt.Sprintf("%s-%s", g.Name, g.Version))
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	// Download the gem file
	url := g.GetDownloadURL()
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download gem from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download gem from %s: status %d", url, resp.StatusCode)
	}

	// Create a temporary file to store the downloaded gem
	tempFile, err := os.CreateTemp("", "*.gem")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy the downloaded content to the temporary file
	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		return fmt.Errorf("failed to save downloaded gem: %w", err)
	}

	// Rewind the temp file for reading
	if _, err := tempFile.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to rewind temporary file: %w", err)
	}

	// Open the .gem file as a tar archive
	tarReader := tar.NewReader(tempFile)

	// Find and extract data.tar.gz
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		if header.Name == "data.tar.gz" {
			// Create a temporary file for data.tar.gz
			dataTarGz, err := os.CreateTemp("", "data.tar.gz")
			if err != nil {
				return fmt.Errorf("failed to create temporary file for data.tar.gz: %w", err)
			}
			defer os.Remove(dataTarGz.Name())
			defer dataTarGz.Close()

			// Copy data.tar.gz content
			if _, err := io.Copy(dataTarGz, tarReader); err != nil {
				return fmt.Errorf("failed to copy data.tar.gz: %w", err)
			}

			// Rewind data.tar.gz for reading
			if _, err := dataTarGz.Seek(0, 0); err != nil {
				return fmt.Errorf("failed to rewind data.tar.gz: %w", err)
			}

			// Open data.tar.gz
			gzReader, err := gzip.NewReader(dataTarGz)
			if err != nil {
				return fmt.Errorf("failed to create gzip reader: %w", err)
			}
			defer gzReader.Close()

			// Extract the contents
			dataTarReader := tar.NewReader(gzReader)
			for {
				header, err := dataTarReader.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					return fmt.Errorf("failed to read data tar header: %w", err)
				}

				// Create the file path
				path := filepath.Join(targetDir, header.Name)

				switch header.Typeflag {
				case tar.TypeDir:
					if err := os.MkdirAll(path, 0755); err != nil {
						return fmt.Errorf("failed to create directory %s: %w", path, err)
					}
				case tar.TypeReg:
					// Ensure the parent directory exists
					if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
						return fmt.Errorf("failed to create parent directory for %s: %w", path, err)
					}

					// Create the file
					file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
					if err != nil {
						return fmt.Errorf("failed to create file %s: %w", path, err)
					}

					// Copy the contents
					if _, err := io.Copy(file, dataTarReader); err != nil {
						file.Close()
						return fmt.Errorf("failed to write file %s: %w", path, err)
					}
					file.Close()
				}
			}
			return nil
		}
	}

	return fmt.Errorf("data.tar.gz not found in gem file")
}

// DefaultSource returns the default RubyGems source
func DefaultSource() Source {
	return Source{
		Type: "rubygems",
		URL:  "https://rubygems.org/",
	}
} 
