package operations

import (
	"fmt"
	"regexp"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// validNamePattern validates template names: lowercase alphanumeric, starting with a-z or 0-9.
var validNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

// supportedLanguages are the languages supported by the init scaffolding.
var supportedLanguages = []string{"go", "typescript", "python"}

// InitInfo holds parameters for initializing a template project.
type InitInfo struct {
	Name     string // Template project name
	Language string // Programming language
	Path     string // Output directory (defaults to ./<name>)
}

// Init initializes a new template project with scaffolded files.
// Both --name and --language are required (no interactive prompts).
func Init(info InitInfo) {
	name := info.Name
	language := info.Language
	path := info.Path

	if name == "" {
		sbClient.PrintError("--name is required")
		return
	}
	if language == "" {
		sbClient.PrintError("--language is required (supported: go, typescript, python)")
		return
	}

	// Validate name
	if !validNamePattern.MatchString(name) {
		sbClient.PrintError("invalid template name %q (must match: [a-z0-9][a-z0-9_-]*)", name)
		return
	}

	// Validate language
	validLang := false
	for _, l := range supportedLanguages {
		if language == l {
			validLang = true
			break
		}
	}
	if !validLang {
		sbClient.PrintError("unsupported language %q (supported: go, typescript, python)", language)
		return
	}

	if path == "" {
		path = "./" + name
	}

	fmt.Printf("Initializing %s template %q in %s...\n", language, name, path)
	if err := scaffold(name, language, path); err != nil {
		sbClient.PrintError("scaffold failed: %v", err)
		return
	}
	sbClient.PrintSuccess("Template %s initialized successfully!", name)
}
