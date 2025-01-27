package gem

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// GemfileLock represents a parsed Gemfile.lock file
type GemfileLock struct {
	Dependencies map[string]*Gem
	content      string
}

// Regular expressions for parsing
var (
	// Matches lines like "    rake (13.0.6)" or "    rails (7.0.8.5)"
	gemSpecRegex = regexp.MustCompile(`^\s+([^\s(]+)\s*\(([^)]+)\)`)
	// Matches lines like "  remote: https://rubygems.org/"
	sourceRegex = regexp.MustCompile(`^\s*remote:\s*(.+)`)
	// Matches section headers like "GEM" or "PATH"
	sectionRegex = regexp.MustCompile(`^(GEM|PATH|PLATFORMS|DEPENDENCIES|BUNDLED WITH)\s*$`)
)

// NewGemfileLock creates a new GemfileLock instance from file contents
func NewGemfileLock(content string) *GemfileLock {
	g := &GemfileLock{
		Dependencies: make(map[string]*Gem),
		content:      content,
	}
	g.parse()
	return g
}

// FromFile reads a Gemfile.lock at the given path and returns a new GemfileLock instance
func FromFile(path string) (*GemfileLock, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return NewGemfileLock(string(content)), nil
}

// parse processes the Gemfile.lock content and extracts dependencies
func (g *GemfileLock) parse() {
	scanner := bufio.NewScanner(strings.NewReader(g.content))
	var currentSection string
	var currentSource Source
	inSpecs := false
	inDependencies := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// Check for section headers
		if sectionMatch := sectionRegex.FindStringSubmatch(trimmedLine); sectionMatch != nil {
			currentSection = sectionMatch[1]
			inSpecs = false
			inDependencies = currentSection == "DEPENDENCIES"
			continue
		}

		// Check for source
		if sourceMatch := sourceRegex.FindStringSubmatch(line); sourceMatch != nil {
			currentSource = Source{
				Type: strings.ToLower(currentSection),
				URL:  strings.TrimSpace(sourceMatch[1]),
			}
			continue
		}

		// Look for the specs subsection
		if trimmedLine == "specs:" {
			inSpecs = true
			continue
		}

		// Only parse specs section under GEM or PATH, skip DEPENDENCIES section
		if !inSpecs || inDependencies {
			continue
		}

		// Parse gem specifications
		if matches := gemSpecRegex.FindStringSubmatch(line); matches != nil {
			name := matches[1]
			// Only store the exact version number from the specs section
			version := strings.TrimSpace(matches[2])
			// Skip if this looks like a constraint (contains ~>, >=, etc.)
			if strings.ContainsAny(version, "~<>=") {
				continue
			}
			g.Dependencies[name] = NewGem(name, version, currentSource)
		}
	}
}

// GetDependency returns a specific gem by name
func (g *GemfileLock) GetDependency(gemName string) *Gem {
	return g.Dependencies[gemName]
}

// GetAllDependencies returns all dependencies as a slice
func (g *GemfileLock) GetAllDependencies() []*Gem {
	gems := make([]*Gem, 0, len(g.Dependencies))
	for _, gem := range g.Dependencies {
		gems = append(gems, gem)
	}
	return gems
} 
