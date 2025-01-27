package gem

import (
	"encoding/json"
	"os"
)

// VersionChange represents a gem that has changed versions
type VersionChange struct {
	Name  string
	Before *Gem
	After  *Gem
}

// VersionChangeJSON represents the JSON structure for serializing a VersionChange
type VersionChangeJSON struct {
	Name         string  `json:"name"`
	BeforeGem    GemJSON `json:"before"`
	AfterGem     GemJSON `json:"after"`
}

// GemfileDiff represents the differences between two Gemfile.lock files
type GemfileDiff struct {
	Added         []*Gem
	Removed       []*Gem
	VersionChanges []VersionChange
}

// DiffJSON represents the JSON structure for serializing a GemfileDiff
type DiffJSON struct {
	Added         []GemJSON           `json:"added"`
	Removed       []GemJSON           `json:"removed"`
	VersionChanges []VersionChangeJSON `json:"version_changes"`
}

// GemJSON represents the JSON structure for serializing a Gem
type GemJSON struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Source  Source `json:"source"`
}

// NewGemfileDiff creates a GemfileDiff by comparing two Gemfile.lock files
func NewGemfileDiff(beforePath, afterPath string) (*GemfileDiff, error) {
	before, err := FromFile(beforePath)
	if err != nil {
		return nil, err
	}

	after, err := FromFile(afterPath)
	if err != nil {
		return nil, err
	}

	return CompareLockfiles(before, after), nil
}

// CompareLockfiles compares two GemfileLock instances and returns their differences
func CompareLockfiles(before, after *GemfileLock) *GemfileDiff {
	diff := &GemfileDiff{
		Added:          make([]*Gem, 0),
		Removed:        make([]*Gem, 0),
		VersionChanges: make([]VersionChange, 0),
	}

	// Find added and changed gems
	for name, afterGem := range after.Dependencies {
		beforeGem := before.Dependencies[name]
		if beforeGem == nil {
			// Gem was added
			diff.Added = append(diff.Added, afterGem)
		} else if beforeGem.Version != afterGem.Version {
			// Version changed
			diff.VersionChanges = append(diff.VersionChanges, VersionChange{
				Name:   name,
				Before: beforeGem,
				After:  afterGem,
			})
		}
	}

	// Find removed gems
	for name, beforeGem := range before.Dependencies {
		if after.Dependencies[name] == nil {
			diff.Removed = append(diff.Removed, beforeGem)
		}
	}

	return diff
}

// HasChanges returns true if there are any differences between the two Gemfile.lock files
func (d *GemfileDiff) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.VersionChanges) > 0
}

// GetAddedGems returns the slice of gems that were added
func (d *GemfileDiff) GetAddedGems() []*Gem {
	return d.Added
}

// GetRemovedGems returns the slice of gems that were removed
func (d *GemfileDiff) GetRemovedGems() []*Gem {
	return d.Removed
}

// GetVersionChanges returns the slice of version changes
func (d *GemfileDiff) GetVersionChanges() []VersionChange {
	return d.VersionChanges
}

// SaveToJSON writes the diff to a JSON file at the specified path
func (d *GemfileDiff) SaveToJSON(path string) error {
	// Convert to JSON-friendly structure
	diffJSON := DiffJSON{
		Added:         make([]GemJSON, len(d.Added)),
		Removed:       make([]GemJSON, len(d.Removed)),
		VersionChanges: make([]VersionChangeJSON, len(d.VersionChanges)),
	}

	// Convert Added gems
	for i, gem := range d.Added {
		diffJSON.Added[i] = GemJSON{
			Name:    gem.Name,
			Version: gem.Version,
			Source:  gem.Source,
		}
	}

	// Convert Removed gems
	for i, gem := range d.Removed {
		diffJSON.Removed[i] = GemJSON{
			Name:    gem.Name,
			Version: gem.Version,
			Source:  gem.Source,
		}
	}

	// Convert Version changes
	for i, change := range d.VersionChanges {
		diffJSON.VersionChanges[i] = VersionChangeJSON{
			Name: change.Name,
			BeforeGem: GemJSON{
				Name:    change.Before.Name,
				Version: change.Before.Version,
				Source:  change.Before.Source,
			},
			AfterGem: GemJSON{
				Name:    change.After.Name,
				Version: change.After.Version,
				Source:  change.After.Source,
			},
		}
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(diffJSON, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(path, data, 0644)
}

// LoadFromJSON creates a new GemfileDiff from a JSON file
func LoadFromJSON(path string) (*GemfileDiff, error) {
	// Read JSON file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON
	var diffJSON DiffJSON
	if err := json.Unmarshal(data, &diffJSON); err != nil {
		return nil, err
	}

	// Create new GemfileDiff
	diff := &GemfileDiff{
		Added:         make([]*Gem, len(diffJSON.Added)),
		Removed:       make([]*Gem, len(diffJSON.Removed)),
		VersionChanges: make([]VersionChange, len(diffJSON.VersionChanges)),
	}

	// Convert Added gems
	for i, gemJSON := range diffJSON.Added {
		diff.Added[i] = NewGem(gemJSON.Name, gemJSON.Version, gemJSON.Source)
	}

	// Convert Removed gems
	for i, gemJSON := range diffJSON.Removed {
		diff.Removed[i] = NewGem(gemJSON.Name, gemJSON.Version, gemJSON.Source)
	}

	// Convert Version changes
	for i, change := range diffJSON.VersionChanges {
		diff.VersionChanges[i] = VersionChange{
			Name:   change.Name,
			Before: NewGem(change.BeforeGem.Name, change.BeforeGem.Version, change.BeforeGem.Source),
			After:  NewGem(change.AfterGem.Name, change.AfterGem.Version, change.AfterGem.Source),
		}
	}

	return diff, nil
} 
