package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// configFile represents the expected structure of a configuration file.
type configFile struct {
	Name     string   `json:"name" yaml:"name"`
	Image    string   `json:"image" yaml:"image"`
	Port     int      `json:"port" yaml:"port"`
	Env      []string `json:"env" yaml:"env"`
	Database struct {
		Host string `json:"host" yaml:"host"`
		Port int    `json:"port" yaml:"port"`
		Name string `json:"name" yaml:"name"`
	} `json:"database" yaml:"database"`
}

// validate checks the parsed config for required fields and valid value ranges.
// It returns all validation errors at once rather than stopping at the first.
func validate(config configFile) error {
	var errors []string

	if config.Name == "" {
		errors = append(errors, "  ✗ 'name' is required but missing or empty")
	}

	if config.Image == "" {
		errors = append(errors, "  ✗ 'image' is required but missing or empty")
	}

	if config.Port < 1 || config.Port > 65535 {
		errors = append(errors, fmt.Sprintf("  ✗ 'port' must be between 1 and 65535, got %d", config.Port))
	}

	if config.Database.Host == "" {
		errors = append(errors, "  ✗ 'database.host' is required but missing or empty")
	}

	if config.Database.Port < 1 || config.Database.Port > 65535 {
		errors = append(errors, fmt.Sprintf("  ✗ 'database.port' must be between 1 and 65535, got %d", config.Database.Port))
	}

	if config.Database.Name == "" {
		errors = append(errors, "  ✗ 'database.name' is required but missing or empty")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

// readFile reads a config file, parses it based on extension, and validates its contents.
func readFile(filename string) error {
	ext := filepath.Ext(filename)
	if ext != ".json" && ext != ".yaml" && ext != ".yml" {
		return fmt.Errorf("unsupported file extension: %s (expected .json, .yaml, or .yml)", ext)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	var config configFile

	// Parse based on file extension
	switch ext {
	case ".json":
		err = json.Unmarshal(data, &config)
		if err != nil {
			return fmt.Errorf("invalid JSON syntax: %w", err)
		}
		fmt.Println("  ✓ Valid JSON syntax")

	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			return fmt.Errorf("invalid YAML syntax: %w", err)
		}
		fmt.Println("  ✓ Valid YAML syntax")
	}

	// Run custom validation on parsed config
	err = validate(config)
	if err != nil {
		return err
	}

	fmt.Println("  ✓ All required fields present")
	fmt.Println("  ✓ All values within valid ranges")

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: config-file-validator <filepath>")
		os.Exit(1)
	}

	fileName := os.Args[1]

	fmt.Println()
	fmt.Printf("Validating %s...\n\n", fileName)

	err := readFile(fileName)
	if err != nil {
		fmt.Printf("\n✗ %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ Configuration is valid!\n")
}
