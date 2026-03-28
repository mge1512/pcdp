package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TemplateRecord represents a PCD template with metadata and content
type TemplateRecord struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Language string `json:"language"`
	Content  string `json:"content"`
}

// ResourceRecord represents a PCD resource (template, prompt, or hints)
type ResourceRecord struct {
	URI     string `json:"uri"`
	Name    string `json:"name"`
	Content string `json:"content,omitempty"`
}

// Diagnostic represents a linting diagnostic with location and severity
type Diagnostic struct {
	Severity string `json:"severity"`
	Line     int    `json:"line"`
	Section  string `json:"section"`
	Message  string `json:"message"`
	Rule     string `json:"rule"`
}

// LintResult represents the result of linting a specification
type LintResult struct {
	Valid       bool         `json:"valid"`
	Errors      int          `json:"errors"`
	Warnings    int          `json:"warnings"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// ============ Filesystem Interface ============

// Filesystem interface for reading files
type Filesystem interface {
	ReadFile(path string) (string, error)
}

// OSFilesystem is the production implementation using os package
type OSFilesystem struct{}

func NewOSFilesystem() Filesystem {
	return &OSFilesystem{}
}

func (fs *OSFilesystem) ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FakeFilesystem is a test double for Filesystem
type FakeFilesystem struct {
	Files   map[string]string
	ReadErr map[string]error
}

func NewFakeFilesystem() *FakeFilesystem {
	return &FakeFilesystem{
		Files:   make(map[string]string),
		ReadErr: make(map[string]error),
	}
}

func (fs *FakeFilesystem) ReadFile(path string) (string, error) {
	if err, ok := fs.ReadErr[path]; ok {
		return "", err
	}
	if content, ok := fs.Files[path]; ok {
		return content, nil
	}
	return "", os.ErrNotExist
}

// ============ Search Path Helpers ============

// templateSearchDirs returns the ordered list of directories to search
// for template files. Later entries take precedence (last-wins merge).
func templateSearchDirs() []string {
	dirs := []string{"/usr/share/pcd/templates"}
	if dirExists("/etc/pcd/templates") {
		dirs = append(dirs, "/etc/pcd/templates")
	}
	if home, err := os.UserHomeDir(); err == nil {
		if d := filepath.Join(home, ".config", "pcd", "templates"); dirExists(d) {
			dirs = append(dirs, d)
		}
	}
	if dirExists(".pcd/templates") {
		dirs = append(dirs, ".pcd/templates")
	}
	return dirs
}

// hintsSearchDirs returns the ordered list of directories to search
// for hints files. Later entries take precedence (last-wins merge).
func hintsSearchDirs() []string {
	dirs := []string{"/usr/share/pcd/hints"}
	if dirExists("/etc/pcd/hints") {
		dirs = append(dirs, "/etc/pcd/hints")
	}
	if home, err := os.UserHomeDir(); err == nil {
		if d := filepath.Join(home, ".config", "pcd", "hints"); dirExists(d) {
			dirs = append(dirs, d)
		}
	}
	if dirExists(".pcd/hints") {
		dirs = append(dirs, ".pcd/hints")
	}
	return dirs
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// parseTemplateMeta extracts Template-For, Version, and default Language
// from a template file's META section and TEMPLATE-TABLE.
func parseTemplateMeta(content string) (name, version, language string) {
	inMeta := false
	inTable := false
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "## META" {
			inMeta = true
			continue
		}
		if inMeta && strings.HasPrefix(trimmed, "## ") {
			inMeta = false
		}
		if inMeta {
			if strings.HasPrefix(trimmed, "Template-For:") {
				name = strings.TrimSpace(strings.TrimPrefix(trimmed, "Template-For:"))
			}
			if strings.HasPrefix(trimmed, "Version:") {
				version = strings.TrimSpace(strings.TrimPrefix(trimmed, "Version:"))
			}
		}
		if trimmed == "## TEMPLATE-TABLE" {
			inTable = true
			continue
		}
		if inTable && strings.HasPrefix(trimmed, "## ") {
			inTable = false
		}
		if inTable && language == "" {
			// Match rows like: | LANGUAGE | Go | default | ... |
			parts := strings.Split(trimmed, "|")
			if len(parts) >= 4 {
				key := strings.TrimSpace(parts[1])
				val := strings.TrimSpace(parts[2])
				constraint := strings.TrimSpace(parts[3])
				if key == "LANGUAGE" && constraint == "default" {
					language = val
				}
			}
		}
	}
	return name, version, language
}

// ============ TemplateStore Interface ============

// TemplateStore interface for accessing templates and hints
type TemplateStore interface {
	ListTemplates() ([]TemplateRecord, error)
	GetTemplate(name, version string) (TemplateRecord, error)
	GetHints(key string) (string, error)
	ListHints() []string
	ListHintsKeys() ([]string, error)
}

// LayeredTemplateStore is the production implementation that reads from
// the filesystem search path hierarchy (later entries take precedence):
//
//	/usr/share/pcd/templates/    (pcd-templates package)
//	/etc/pcd/templates/          (system admin overrides)
//	~/.config/pcd/templates/     (user overrides)
//	./.pcd/templates/            (project-local)
type LayeredTemplateStore struct {
	templates map[string]map[string]TemplateRecord
	hints     map[string]string
}

func NewLayeredTemplateStore() TemplateStore {
	ts := &LayeredTemplateStore{
		templates: make(map[string]map[string]TemplateRecord),
		hints:     make(map[string]string),
	}
	for _, dir := range templateSearchDirs() {
		ts.loadTemplatesFrom(dir)
	}
	for _, dir := range hintsSearchDirs() {
		ts.loadHintsFrom(dir)
	}
	return ts
}

func (ts *LayeredTemplateStore) loadTemplatesFrom(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return // directory absent or unreadable — skip silently
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".template.md") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		content := string(data)
		name, version, language := parseTemplateMeta(content)
		if name == "" {
			name = strings.TrimSuffix(e.Name(), ".template.md")
		}
		if version == "" {
			version = "latest"
		}
		if _, ok := ts.templates[name]; !ok {
			ts.templates[name] = make(map[string]TemplateRecord)
		}
		// later dir entries overwrite earlier ones (last-wins)
		ts.templates[name][version] = TemplateRecord{
			Name:     name,
			Version:  version,
			Language: language,
			Content:  content,
		}
	}
}

func (ts *LayeredTemplateStore) loadHintsFrom(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return // directory absent or unreadable — skip silently
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".hints.md") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		key := strings.TrimSuffix(e.Name(), ".hints.md")
		ts.hints[key] = string(data) // last-wins
	}
}

func (ts *LayeredTemplateStore) ListTemplates() ([]TemplateRecord, error) {
	var results []TemplateRecord
	for _, versions := range ts.templates {
		for _, rec := range versions {
			results = append(results, rec)
		}
	}
	return results, nil
}

func (ts *LayeredTemplateStore) GetTemplate(name, version string) (TemplateRecord, error) {
	if versions, ok := ts.templates[name]; ok {
		if version == "latest" {
			for _, rec := range versions {
				return rec, nil
			}
		}
		if rec, ok := versions[version]; ok {
			return rec, nil
		}
		return TemplateRecord{}, fmt.Errorf("version %s not found for template %s", version, name)
	}
	return TemplateRecord{}, fmt.Errorf("unknown template: %s", name)
}

func (ts *LayeredTemplateStore) GetHints(key string) (string, error) {
	if content, ok := ts.hints[key]; ok {
		return content, nil
	}
	return "", fmt.Errorf("hints not found: %s", key)
}

func (ts *LayeredTemplateStore) ListHints() []string {
	var keys []string
	for k := range ts.hints {
		keys = append(keys, k)
	}
	return keys
}

func (ts *LayeredTemplateStore) ListHintsKeys() ([]string, error) {
	return ts.ListHints(), nil
}

// FakeTemplateStore is a test double for TemplateStore
type FakeTemplateStore struct {
	Templates []TemplateRecord
	Hints     map[string]string
}

func NewFakeTemplateStore() *FakeTemplateStore {
	return &FakeTemplateStore{
		Templates: []TemplateRecord{},
		Hints:     make(map[string]string),
	}
}

func (ts *FakeTemplateStore) ListTemplates() ([]TemplateRecord, error) {
	return ts.Templates, nil
}

func (ts *FakeTemplateStore) GetTemplate(name, version string) (TemplateRecord, error) {
	for _, t := range ts.Templates {
		if t.Name == name {
			if version == "latest" || version == t.Version {
				return t, nil
			}
		}
	}
	return TemplateRecord{}, fmt.Errorf("unknown template: %s", name)
}

func (ts *FakeTemplateStore) GetHints(key string) (string, error) {
	if content, ok := ts.Hints[key]; ok {
		return content, nil
	}
	return "", fmt.Errorf("hints not found: %s", key)
}

func (ts *FakeTemplateStore) ListHints() []string {
	var keys []string
	for k := range ts.Hints {
		keys = append(keys, k)
	}
	return keys
}

func (ts *FakeTemplateStore) ListHintsKeys() ([]string, error) {
	return ts.ListHints(), nil
}

// ============ PromptStore Interface ============

// PromptStore interface for accessing embedded prompts
type PromptStore interface {
	GetPrompt(name string) (string, error)
	ListPrompts() []string
}

// EmbeddedPromptStore is the production implementation
// with prompts embedded as Go string constants
type EmbeddedPromptStore struct {
	prompts map[string]string
}

func NewEmbeddedPromptStore() PromptStore {
	store := &EmbeddedPromptStore{
		prompts: make(map[string]string),
	}
	store.prompts["interview"] = promptInterview
	store.prompts["translator"] = promptTranslator
	return store
}

func (ps *EmbeddedPromptStore) GetPrompt(name string) (string, error) {
	if content, ok := ps.prompts[name]; ok {
		return content, nil
	}
	return "", fmt.Errorf("prompt not found: %s", name)
}

func (ps *EmbeddedPromptStore) ListPrompts() []string {
	return []string{"interview", "translator"}
}

// FakePromptStore is a test double for PromptStore
type FakePromptStore struct {
	Prompts map[string]string
}

func NewFakePromptStore() *FakePromptStore {
	return &FakePromptStore{
		Prompts: make(map[string]string),
	}
}

func (ps *FakePromptStore) GetPrompt(name string) (string, error) {
	if content, ok := ps.Prompts[name]; ok {
		return content, nil
	}
	return "", fmt.Errorf("prompt not found: %s", name)
}

func (ps *FakePromptStore) ListPrompts() []string {
	var names []string
	for name := range ps.Prompts {
		names = append(names, name)
	}
	return names
}
