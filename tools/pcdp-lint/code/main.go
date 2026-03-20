package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Constants for exit codes
const (
	ExitOK              = 0
	ExitLintError       = 1
	ExitInvocationError = 2
)

// Severity levels for diagnostics
type Severity int

const (
	Error Severity = iota
	Warning
)

func (s Severity) String() string {
	switch s {
	case Error:
		return "ERROR"
	case Warning:
		return "WARNING"
	default:
		return "UNKNOWN"
	}
}

// Diagnostic represents a lint diagnostic
type Diagnostic struct {
	Severity Severity
	Section  string
	Message  string
	Line     int
}

// LintResult represents the result of linting a file
type LintResult struct {
	File        string
	Diagnostics []Diagnostic
	ExitCode    int
}

// Built-in SPDX license identifiers (subset for demonstration)
var spdxLicenses = map[string]bool{
	"Apache-2.0":          true,
	"MIT":                 true,
	"GPL-2.0-only":        true,
	"GPL-3.0-only":        true,
	"LGPL-2.1-or-later":   true,
	"BSD-2-Clause":        true,
	"BSD-3-Clause":        true,
	"ISC":                 true,
	"CC0-1.0":             true,
	"CC-BY-4.0":           true,
	"MPL-2.0":             true,
	"Unlicense":           true,
	"GPL-2.0-or-later":    true,
	"GPL-3.0-or-later":    true,
	"LGPL-2.1-only":       true,
	"LGPL-3.0-only":       true,
	"LGPL-3.0-or-later":   true,
	"AGPL-3.0-only":       true,
	"AGPL-3.0-or-later":   true,
}

// Known deployment templates
var deploymentTemplates = []string{
	"wasm", "ebpf", "kernel-module", "verified-library",
	"cli-tool", "gui-tool", "cloud-native", "backend-service",
	"library-c-abi", "enterprise-software", "academic",
	"python-tool", "enhance-existing", "manual", "template",
}

// Known verification values
var knownVerificationValues = []string{
	"none", "lean4", "fstar", "dafny", "custom",
}

// Required sections
var requiredSections = []string{
	"## META", "## TYPES", "## BEHAVIOR", "## PRECONDITIONS",
	"## POSTCONDITIONS", "## INVARIANTS", "## EXAMPLES",
}

// Required META fields
var requiredMetaFields = []string{
	"Deployment", "Verification", "Safety-Level",
	"Version", "Spec-Schema", "License",
}

func main() {
	// Handle SIGTERM and SIGINT for clean exit
	// For a short-lived CLI tool, we rely on Go runtime default behavior

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "error: missing file argument\n")
		fmt.Fprintf(os.Stderr, "usage: pcdp-lint [strict=true] <specfile.md>\n")
		fmt.Fprintf(os.Stderr, "       pcdp-lint list-templates\n")
		fmt.Fprintf(os.Stderr, "       pcdp-lint version\n")
		os.Exit(ExitInvocationError)
	}

	// Parse key=value options and commands
	var strict bool
	var filename string
	var command string

	for _, arg := range args {
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			key, value := parts[0], parts[1]
			switch key {
			case "strict":
				if value == "true" {
					strict = true
				} else if value == "false" {
					strict = false
				} else {
					fmt.Fprintf(os.Stderr, "error: strict must be true or false\n")
					os.Exit(ExitInvocationError)
				}
			default:
				fmt.Fprintf(os.Stderr, "error: unrecognised option: %s\n", key)
				os.Exit(ExitInvocationError)
			}
		} else if arg == "list-templates" {
			command = "list-templates"
		} else if arg == "version" {
			command = "version"
		} else {
			if filename != "" {
				fmt.Fprintf(os.Stderr, "error: multiple file arguments not supported\n")
				os.Exit(ExitInvocationError)
			}
			filename = arg
		}
	}

	// Handle commands
	if command == "list-templates" {
		listTemplates()
		os.Exit(ExitOK)
	}

	if command == "version" {
		fmt.Printf("pcdp-lint 0.3.7 (schema 0.3.7) spdx/3.21\n")
		os.Exit(ExitOK)
	}

	// Validate filename
	if filename == "" {
		fmt.Fprintf(os.Stderr, "error: missing file argument\n")
		os.Exit(ExitInvocationError)
	}

	if !strings.HasSuffix(filename, ".md") {
		fmt.Fprintf(os.Stderr, "error: file must have .md extension: %s\n", filename)
		os.Exit(ExitInvocationError)
	}

	// Check if file exists and is readable
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: cannot open file: %s\n", filename)
		os.Exit(ExitInvocationError)
	}

	// Perform lint
	result := lintFile(filename, strict)

	// Output diagnostics to stderr
	for _, diag := range result.Diagnostics {
		fmt.Fprintf(os.Stderr, "%s  %s:%d  [%s]  %s\n",
			diag.Severity, result.File, diag.Line, diag.Section, diag.Message)
	}

	// Output summary to stdout
	errorCount := 0
	warningCount := 0
	for _, diag := range result.Diagnostics {
		if diag.Severity == Error {
			errorCount++
		} else {
			warningCount++
		}
	}

	if result.ExitCode == ExitOK {
		if warningCount == 0 {
			fmt.Printf("✓ %s: valid\n", result.File)
		} else {
			fmt.Printf("✓ %s: valid (%d warning(s))\n", result.File, warningCount)
		}
	} else {
		if strict && errorCount == 0 && warningCount > 0 {
			fmt.Printf("✗ %s: %d error(s), %d warning(s) [strict mode]\n",
				result.File, errorCount, warningCount)
		} else {
			fmt.Printf("✗ %s: %d error(s), %d warning(s)\n",
				result.File, errorCount, warningCount)
		}
	}

	os.Exit(result.ExitCode)
}

func listTemplates() {
	templateDefaults := map[string]string{
		"wasm":                 "(template file not found)",
		"ebpf":                 "(template file not found)",
		"kernel-module":        "(template file not found)",
		"verified-library":     "(template file not found)",
		"cli-tool":             "Go",
		"gui-tool":             "(template file not found)",
		"cloud-native":         "(template file not found)",
		"backend-service":      "(template file not found)",
		"library-c-abi":        "(template file not found)",
		"enterprise-software":  "(template file not found)",
		"academic":             "(template file not found)",
		"python-tool":          "Python",
		"enhance-existing":     "(declare Language: in META)",
		"manual":               "(declare Target: in META)",
		"template":             "(template definition file, not translatable)",
	}

	for _, template := range deploymentTemplates {
		defaultLang := templateDefaults[template]
		fmt.Printf("%s  →  %s\n", template, defaultLang)
	}
}

func lintFile(filename string, strict bool) LintResult {
	result := LintResult{
		File:        filename,
		Diagnostics: []Diagnostic{},
		ExitCode:    ExitOK,
	}

	file, err := os.Open(filename)
	if err != nil {
		// This should not happen as we check file existence before calling this
		result.ExitCode = ExitInvocationError
		return result
	}
	defer file.Close()

	// Read file content
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		result.ExitCode = ExitInvocationError
		return result
	}

	// Parse sections
	sections := parseSections(lines)

	// Run validation rules
	result.Diagnostics = append(result.Diagnostics, validateRequiredSections(sections)...)
	result.Diagnostics = append(result.Diagnostics, validateMetaFields(sections)...)
	result.Diagnostics = append(result.Diagnostics, validateDeploymentTemplate(sections)...)
	result.Diagnostics = append(result.Diagnostics, validateDeprecatedFields(sections)...)
	result.Diagnostics = append(result.Diagnostics, validateVerificationField(sections)...)
	result.Diagnostics = append(result.Diagnostics, validateExamplesSection(sections, lines)...)

	// Determine exit code
	hasError := false
	hasWarning := false
	for _, diag := range result.Diagnostics {
		if diag.Severity == Error {
			hasError = true
		} else {
			hasWarning = true
		}
	}

	if hasError {
		result.ExitCode = ExitLintError
	} else if strict && hasWarning {
		result.ExitCode = ExitLintError
	}

	return result
}

// Section represents a parsed section
type Section struct {
	Name    string
	StartLine int
	Content []string
}

func parseSections(lines []string) map[string]Section {
	sections := make(map[string]Section)
	var currentSection *Section

	for i, line := range lines {
		if strings.HasPrefix(line, "## ") {
			// Save previous section
			if currentSection != nil {
				sections[currentSection.Name] = *currentSection
			}
			
			// Start new section
			currentSection = &Section{
				Name:      line,
				StartLine: i + 1,
				Content:   []string{},
			}
		} else if currentSection != nil {
			currentSection.Content = append(currentSection.Content, line)
		}
	}

	// Save last section
	if currentSection != nil {
		sections[currentSection.Name] = *currentSection
	}

	return sections
}

func validateRequiredSections(sections map[string]Section) []Diagnostic {
	var diagnostics []Diagnostic
	
	for _, required := range requiredSections {
		found := false
		
		// Special handling for BEHAVIOR section
		if required == "## BEHAVIOR" {
			for sectionName := range sections {
				if sectionName == "## BEHAVIOR" ||
					strings.HasPrefix(sectionName, "## BEHAVIOR:") ||
					strings.HasPrefix(sectionName, "## BEHAVIOR/INTERNAL:") {
					found = true
					break
				}
			}
		} else {
			_, found = sections[required]
		}
		
		if !found {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Error,
				Section:  "structure",
				Message:  fmt.Sprintf("Missing required section: %s", required),
				Line:     1,
			})
		}
	}
	
	return diagnostics
}

func validateMetaFields(sections map[string]Section) []Diagnostic {
	var diagnostics []Diagnostic
	
	metaSection, exists := sections["## META"]
	if !exists {
		return diagnostics // Already caught by validateRequiredSections
	}
	
	// Parse META fields
	metaFields := parseMetaFields(metaSection.Content)
	
	// Check required fields
	for _, required := range requiredMetaFields {
		if _, exists := metaFields[required]; !exists {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Error,
				Section:  "META",
				Message:  fmt.Sprintf("Missing required META field: %s", required),
				Line:     metaSection.StartLine,
			})
		}
	}
	
	// Check for at least one Author field
	if len(metaFields["Author"]) == 0 {
		diagnostics = append(diagnostics, Diagnostic{
			Severity: Error,
			Section:  "META",
			Message:  "Missing required META field: Author (at least one Author: line required)",
			Line:     metaSection.StartLine,
		})
	}
	
	// Validate field values
	for field, values := range metaFields {
		for _, value := range values {
			if strings.TrimSpace(value) == "" {
				diagnostics = append(diagnostics, Diagnostic{
					Severity: Error,
					Section:  "META",
					Message:  fmt.Sprintf("META field %s has empty value", field),
					Line:     metaSection.StartLine,
				})
			}
		}
	}
	
	// Validate Version format
	if versions, exists := metaFields["Version"]; exists && len(versions) > 0 {
		version := versions[0]
		if !isValidSemanticVersion(version) {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Error,
				Section:  "META",
				Message:  fmt.Sprintf("Version '%s' is not valid semantic versioning. Required format: MAJOR.MINOR.PATCH (e.g. 0.1.0)", version),
				Line:     metaSection.StartLine,
			})
		}
	}
	
	// Validate Spec-Schema format
	if schemas, exists := metaFields["Spec-Schema"]; exists && len(schemas) > 0 {
		schema := schemas[0]
		if !isValidSemanticVersion(schema) {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Error,
				Section:  "META",
				Message:  fmt.Sprintf("Spec-Schema '%s' is not valid semantic versioning. Required format: MAJOR.MINOR.PATCH (e.g. 0.1.0)", schema),
				Line:     metaSection.StartLine,
			})
		}
	}
	
	// Validate License SPDX
	if licenses, exists := metaFields["License"]; exists && len(licenses) > 0 {
		license := licenses[0]
		if !isValidSPDXLicense(license) {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Error,
				Section:  "META",
				Message:  fmt.Sprintf("License '%s' is not a valid SPDX identifier. See https://spdx.org/licenses/ for valid identifiers. Compound expressions permitted (e.g. Apache-2.0 OR MIT).", license),
				Line:     metaSection.StartLine,
			})
		}
	}
	
	return diagnostics
}

func parseMetaFields(content []string) map[string][]string {
	fields := make(map[string][]string)
	
	for _, line := range content {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "---") {
			continue
		}
		
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			fields[key] = append(fields[key], value)
		}
	}
	
	return fields
}

func isValidSemanticVersion(version string) bool {
	pattern := `^[0-9]+\.[0-9]+\.[0-9]+$`
	matched, _ := regexp.MatchString(pattern, version)
	return matched
}

func isValidSPDXLicense(license string) bool {
	// Simple validation - check if it's in our known list or contains OR/AND
	license = strings.TrimSpace(license)
	
	// Check if it's a simple license
	if spdxLicenses[license] {
		return true
	}
	
	// Check for compound expressions (simplified check for OR/AND)
	if strings.Contains(license, " OR ") || strings.Contains(license, " AND ") {
		// Create a regex to split on OR/AND
		re := regexp.MustCompile(` (OR|AND) `)
		parts := re.Split(license, -1)
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if !spdxLicenses[part] {
				return false
			}
		}
		return true
	}
	
	return false
}

func validateDeploymentTemplate(sections map[string]Section) []Diagnostic {
	var diagnostics []Diagnostic
	
	metaSection, exists := sections["## META"]
	if !exists {
		return diagnostics
	}
	
	metaFields := parseMetaFields(metaSection.Content)
	deployments, exists := metaFields["Deployment"]
	if !exists || len(deployments) == 0 {
		return diagnostics // Already caught by validateMetaFields
	}
	
	deployment := deployments[0]
	
	// Check for retired crypto-library
	if deployment == "crypto-library" {
		diagnostics = append(diagnostics, Diagnostic{
			Severity: Error,
			Section:  "META",
			Message:  "Deployment 'crypto-library' was retired in 0.3.6. Use 'verified-library' instead. verified-library covers all safety- and security-critical C-ABI libraries including cryptographic primitives.",
			Line:     1,
		})
		return diagnostics
	}
	
	// Check if deployment is known
	validTemplate := false
	for _, template := range deploymentTemplates {
		if deployment == template {
			validTemplate = true
			break
		}
	}
	
	if !validTemplate {
		diagnostics = append(diagnostics, Diagnostic{
			Severity: Error,
			Section:  "META",
			Message:  fmt.Sprintf("Unknown deployment template: '%s'. Run 'pcdp-lint list-templates' to see valid values.", deployment),
			Line:     metaSection.StartLine,
		})
		return diagnostics
	}
	
	// Special validation for specific deployments
	switch deployment {
	case "enhance-existing":
		if _, exists := metaFields["Language"]; !exists {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Error,
				Section:  "META",
				Message:  "Deployment 'enhance-existing' requires META field 'Language'",
				Line:     metaSection.StartLine,
			})
		} else if langs := metaFields["Language"]; len(langs) > 0 && strings.TrimSpace(langs[0]) == "" {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Error,
				Section:  "META",
				Message:  "META field 'Language' has empty value",
				Line:     metaSection.StartLine,
			})
		}
		
	case "manual":
		if _, exists := metaFields["Target"]; !exists {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Error,
				Section:  "META",
				Message:  "Deployment 'manual' requires META field 'Target' (no template available for language resolution)",
				Line:     metaSection.StartLine,
			})
		}
		
	case "python-tool":
		if safetyLevels, exists := metaFields["Safety-Level"]; exists && len(safetyLevels) > 0 {
			if safetyLevels[0] != "QM" {
				diagnostics = append(diagnostics, Diagnostic{
					Severity: Error,
					Section:  "META",
					Message:  "Deployment 'python-tool' requires Safety-Level: QM. Python is not suitable for safety-critical components.",
					Line:     metaSection.StartLine,
				})
			}
		}
		
		if verifications, exists := metaFields["Verification"]; exists && len(verifications) > 0 {
			if verifications[0] != "none" {
				diagnostics = append(diagnostics, Diagnostic{
					Severity: Error,
					Section:  "META",
					Message:  "Deployment 'python-tool' requires Verification: none. No formal verification path exists for Python.",
					Line:     metaSection.StartLine,
				})
			}
		}
		
	case "verified-library":
		if safetyLevels, exists := metaFields["Safety-Level"]; exists && len(safetyLevels) > 0 {
			if safetyLevels[0] == "QM" {
				diagnostics = append(diagnostics, Diagnostic{
					Severity: Warning,
					Section:  "META",
					Message:  "Deployment 'verified-library' with Safety-Level: QM is unusual. verified-library is intended for safety- or security-critical components. Consider using library-c-abi for general-purpose libraries.",
					Line:     metaSection.StartLine,
				})
			}
		}
	}
	
	return diagnostics
}

func validateDeprecatedFields(sections map[string]Section) []Diagnostic {
	var diagnostics []Diagnostic
	
	metaSection, exists := sections["## META"]
	if !exists {
		return diagnostics
	}
	
	metaFields := parseMetaFields(metaSection.Content)
	
	// Check for deprecated Target field
	if _, exists := metaFields["Target"]; exists {
		deployments := metaFields["Deployment"]
		if len(deployments) == 0 || deployments[0] != "manual" {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Warning,
				Section:  "META",
				Message:  "META field 'Target' is deprecated since v0.3.0. Target language is derived from the deployment template. Remove 'Target', or switch to Deployment: manual if explicit language control is required.",
				Line:     metaSection.StartLine,
			})
		}
	}
	
	// Check for deprecated Domain field
	if _, exists := metaFields["Domain"]; exists {
		diagnostics = append(diagnostics, Diagnostic{
			Severity: Warning,
			Section:  "META",
			Message:  "META field 'Domain' is deprecated since v0.3.0. Use 'Deployment' instead.",
			Line:     metaSection.StartLine,
		})
	}
	
	return diagnostics
}

func validateVerificationField(sections map[string]Section) []Diagnostic {
	var diagnostics []Diagnostic
	
	metaSection, exists := sections["## META"]
	if !exists {
		return diagnostics
	}
	
	metaFields := parseMetaFields(metaSection.Content)
	verifications, exists := metaFields["Verification"]
	if !exists || len(verifications) == 0 {
		return diagnostics
	}
	
	verification := verifications[0]
	validVerification := false
	for _, known := range knownVerificationValues {
		if verification == known {
			validVerification = true
			break
		}
	}
	
	if !validVerification {
		diagnostics = append(diagnostics, Diagnostic{
			Severity: Warning,
			Section:  "META",
			Message:  fmt.Sprintf("Unknown verification value: '%s'. Known values: none, lean4, fstar, dafny, custom. Custom verification backends are permitted; verify the value is intentional.", verification),
			Line:     metaSection.StartLine,
		})
	}
	
	return diagnostics
}

func validateExamplesSection(sections map[string]Section, lines []string) []Diagnostic {
	var diagnostics []Diagnostic
	
	examplesSection, exists := sections["## EXAMPLES"]
	if !exists {
		return diagnostics // Already caught by validateRequiredSections
	}
	
	// Parse examples
	examples := parseExamples(examplesSection.Content, examplesSection.StartLine)
	
	if len(examples) == 0 {
		diagnostics = append(diagnostics, Diagnostic{
			Severity: Error,
			Section:  "EXAMPLES",
			Message:  "EXAMPLES section contains no example blocks. Each example requires EXAMPLE:, GIVEN:, WHEN:, THEN: markers.",
			Line:     examplesSection.StartLine,
		})
		return diagnostics
	}
	
	// Validate each example
	for _, example := range examples {
		if !example.HasGiven {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Error,
				Section:  "EXAMPLES",
				Message:  fmt.Sprintf("Example '%s' missing GIVEN: marker", example.Name),
				Line:     example.Line,
			})
		}
		
		if !example.HasWhen {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Error,
				Section:  "EXAMPLES",
				Message:  fmt.Sprintf("Example '%s' missing WHEN: marker", example.Name),
				Line:     example.Line,
			})
		}
		
		if !example.HasThen {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Error,
				Section:  "EXAMPLES",
				Message:  fmt.Sprintf("Example '%s' missing THEN: marker", example.Name),
				Line:     example.Line,
			})
		}
		
		// Check for empty blocks
		if example.HasGiven && len(example.GivenContent) == 0 {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Warning,
				Section:  "EXAMPLES",
				Message:  fmt.Sprintf("Example '%s' has empty GIVEN block", example.Name),
				Line:     example.Line,
			})
		}
		
		if example.HasWhen && len(example.WhenContent) == 0 {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Warning,
				Section:  "EXAMPLES",
				Message:  fmt.Sprintf("Example '%s' has empty WHEN block", example.Name),
				Line:     example.Line,
			})
		}
		
		if example.HasThen && len(example.ThenContent) == 0 {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: Warning,
				Section:  "EXAMPLES",
				Message:  fmt.Sprintf("Example '%s' has empty THEN block", example.Name),
				Line:     example.Line,
			})
		}
	}
	
	return diagnostics
}

// Example represents a parsed example block
type Example struct {
	Name         string
	Line         int
	HasGiven     bool
	HasWhen      bool
	HasThen      bool
	GivenContent []string
	WhenContent  []string
	ThenContent  []string
}

func parseExamples(content []string, startLine int) []Example {
	var examples []Example
	var currentExample *Example
	var currentBlock string
	
	for i, line := range content {
		line = strings.TrimSpace(line)
		lineNumber := startLine + i + 1
		
		if strings.HasPrefix(line, "EXAMPLE:") {
			// Save previous example
			if currentExample != nil {
				examples = append(examples, *currentExample)
			}
			
			// Start new example
			name := strings.TrimSpace(strings.TrimPrefix(line, "EXAMPLE:"))
			currentExample = &Example{
				Name:         name,
				Line:         lineNumber,
				GivenContent: []string{},
				WhenContent:  []string{},
				ThenContent:  []string{},
			}
			currentBlock = ""
		} else if currentExample != nil {
			if strings.HasPrefix(line, "GIVEN:") {
				currentExample.HasGiven = true
				currentBlock = "given"
			} else if strings.HasPrefix(line, "WHEN:") {
				currentExample.HasWhen = true
				currentBlock = "when"
			} else if strings.HasPrefix(line, "THEN:") {
				currentExample.HasThen = true
				currentBlock = "then"
			} else if line != "" && currentBlock != "" {
				switch currentBlock {
				case "given":
					currentExample.GivenContent = append(currentExample.GivenContent, line)
				case "when":
					currentExample.WhenContent = append(currentExample.WhenContent, line)
				case "then":
					currentExample.ThenContent = append(currentExample.ThenContent, line)
				}
			}
		}
	}
	
	// Save last example
	if currentExample != nil {
		examples = append(examples, *currentExample)
	}
	
	return examples
}