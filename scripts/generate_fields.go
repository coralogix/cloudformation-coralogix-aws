package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type CloudFormationTemplate struct {
	Parameters map[string]Parameter `yaml:"Parameters"`
}

type Parameter struct {
	Type          string   `yaml:"Type"`
	Description   string   `yaml:"Description"`
	Default       any      `yaml:"Default"`
	AllowedValues []string `yaml:"AllowedValues"`
}

type FieldsConfig struct {
	Fields []Field `yaml:"fields"`
}

type Field struct {
	Name         string `yaml:"name"`
	Visible      bool   `yaml:"visible"`
	Required     bool   `yaml:"required"`
	Predefined   bool   `yaml:"predefined"`
	Type         string `yaml:"type"`
	DefaultValue any    `yaml:"default_value"`
	DisplayName  string `yaml:"display_name"`
	TooltipText  string `yaml:"tooltip_text"`
	Placeholder  string `yaml:"placeholder"`
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run generate_fields.go <integration_name> <new_version> <integration_definitions_path>")
		fmt.Println("Example: go run generate_fields.go aws-shipper-lambda 1.3.12 /path/to/integration-definitions")
		os.Exit(1)
	}

	integrationName := os.Args[1]
	newVersion := os.Args[2]
	integrationDefinitionsPath := os.Args[3]

	// Determine integration path
	var integrationPath string
	switch integrationName {
	case "aws-shipper-lambda":
		integrationPath = filepath.Join(integrationDefinitionsPath, "integrations", "shared", "aws-shipper")
	case "firehose-logs":
		integrationPath = filepath.Join(integrationDefinitionsPath, "integrations", "shared", "firehose-logs")
	default:
		integrationPath = filepath.Join(integrationDefinitionsPath, "integrations", integrationName)
	}

	// Find the most recent version directory
	latestVersion, err := findLatestVersion(integrationPath)
	if err != nil {
		fmt.Printf("Error finding latest version: %v\n", err)
		os.Exit(1)
	}

	latestFieldsFile := filepath.Join(latestVersion, "fields.yaml")
	newVersionDir := filepath.Join(integrationPath, "v"+newVersion)
	newFieldsFile := filepath.Join(newVersionDir, "fields.yaml")
	newTemplateFile := filepath.Join(newVersionDir, "template.yaml")

	// Create new version directory
	if err := os.MkdirAll(newVersionDir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Copy or create fields.yaml
	var existingFields FieldsConfig
	if _, err := os.Stat(latestFieldsFile); err == nil {
		// Copy existing fields.yaml
		data, err := os.ReadFile(latestFieldsFile)
		if err != nil {
			fmt.Printf("Error reading latest fields file: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(newFieldsFile, data, 0644); err != nil {
			fmt.Printf("Error copying fields file: %v\n", err)
			os.Exit(1)
		}

		// Parse existing fields
		if err := yaml.Unmarshal(data, &existingFields); err != nil {
			fmt.Printf("Error parsing existing fields: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Copied fields.yaml from %s to %s\n", latestFieldsFile, newFieldsFile)
	} else {
		// Create basic fields.yaml
		existingFields = FieldsConfig{
			Fields: []Field{
				{
					Name:         "ApiKey",
					Visible:      true,
					Required:     true,
					Predefined:   false,
					Type:         "api_key",
					DefaultValue: "",
					DisplayName:  "API Key",
					TooltipText:  "Coralogix Send-Your-Data API Key",
					Placeholder:  "Your API key",
				},
			},
		}

		if err := writeFieldsConfig(newFieldsFile, existingFields); err != nil {
			fmt.Printf("Error creating basic fields file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created basic fields.yaml at %s\n", newFieldsFile)
	}

	// Check for new parameters in template
	if _, err := os.Stat(newTemplateFile); err == nil {
		fmt.Printf("Checking for new parameters in %s\n", newTemplateFile)

		if err := addMissingParameters(newTemplateFile, newFieldsFile, existingFields); err != nil {
			fmt.Printf("Error adding missing parameters: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Warning: New template file not found at %s\n", newTemplateFile)
	}
}

func findLatestVersion(integrationPath string) (string, error) {
	var versions []string

	err := filepath.WalkDir(integrationPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && strings.HasPrefix(d.Name(), "v") && path != integrationPath {
			versions = append(versions, path)
			return filepath.SkipDir // Don't go deeper into version directories
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no previous version found in %s", integrationPath)
	}

	// Sort versions and return the latest
	sort.Slice(versions, func(i, j int) bool {
		return filepath.Base(versions[i]) < filepath.Base(versions[j])
	})

	return versions[len(versions)-1], nil
}

func addMissingParameters(templateFile, fieldsFile string, existingFields FieldsConfig) error {
	// Read template
	templateData, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("reading template: %w", err)
	}

	var template CloudFormationTemplate
	if err := yaml.Unmarshal(templateData, &template); err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	// Get existing field names
	existingFieldNames := make(map[string]bool)
	for _, field := range existingFields.Fields {
		existingFieldNames[field.Name] = true
	}

	// Find missing parameters
	var missingParams []string
	for paramName := range template.Parameters {
		if !existingFieldNames[paramName] {
			missingParams = append(missingParams, paramName)
		}
	}

	if len(missingParams) == 0 {
		fmt.Println("No new parameters found")
		return nil
	}

	fmt.Printf("Adding missing parameters: %v\n", missingParams)

	// Add missing parameters to existing fields
	for _, paramName := range missingParams {
		param := template.Parameters[paramName]

		defaultValue := ""
		if param.Default != nil {
			defaultValue = fmt.Sprintf("%v", param.Default)
		}

		displayName := camelCaseToDisplayName(paramName)

		newField := Field{
			Name:         paramName,
			Visible:      true,
			Required:     false,
			Predefined:   false,
			Type:         "text",
			DefaultValue: defaultValue,
			DisplayName:  displayName,
			TooltipText:  param.Description,
			Placeholder:  displayName,
		}

		existingFields.Fields = append(existingFields.Fields, newField)
	}

	// Write updated fields
	return writeFieldsConfig(fieldsFile, existingFields)
}

func camelCaseToDisplayName(s string) string {
	// Insert space before uppercase letters
	re := regexp.MustCompile(`([A-Z])`)
	result := re.ReplaceAllString(s, " $1")

	// Trim leading space and return
	return strings.TrimSpace(result)
}

func writeFieldsConfig(filename string, config FieldsConfig) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("marshaling fields config: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}
