package semgrep

import (
	"fmt"
	"strings"
)

// Finding represents a single semgrep finding
type Finding struct {
	RuleID  string `json:"check_id"`
	Message string `json:"message"`
	Lines   string `json:"lines"`
	Line    int    `json:"line"`
	Path    string `json:"path"`
}

// NewFinding creates a Finding from a semgrep result
func NewFinding(result map[string]interface{}) *Finding {
	extra := result["extra"].(map[string]interface{})
	start := result["start"].(map[string]interface{})

	return &Finding{
		RuleID:  result["check_id"].(string),
		Message: extra["message"].(string),
		Lines:   extra["lines"].(string),
		Line:    int(start["line"].(float64)),
		Path:    result["path"].(string),
	}
}

// Equals checks if two findings are equivalent (same rule and lines)
func (f *Finding) Equals(other *Finding) bool {
	return f.RuleID == other.RuleID && f.Lines == other.Lines
}

// Display returns a formatted string representation of the finding
func (f *Finding) Display() string {
	return fmt.Sprintf("  [%s] line %d: %s\n    %s",
		f.RuleID,
		f.Line,
		f.Message,
		f.Lines)
}

// Rebase updates the path to be relative to the given base directory
func (f *Finding) Rebase(baseDir string) error {
	if !strings.HasPrefix(f.Path, baseDir) {
		return fmt.Errorf("base doesn't match: path %s does not start with %s", f.Path, baseDir)
	}

	f.Path = strings.TrimPrefix(f.Path, baseDir+"/")
	return nil
}

// RelativePath returns the path relative to the given base directory
func (f *Finding) RelativePath(baseDir string) string {
	return strings.TrimPrefix(f.Path, baseDir+"/")
} 
