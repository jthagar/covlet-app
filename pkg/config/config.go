package config

import (
    "errors"
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

// Resume holds all the information from the YAML file.
type Resume struct {
	Name             string       `yaml:"name"`
	Email            string       `yaml:"email"`
	Phone            string       `yaml:"phone"`
	Address          string       `yaml:"address"`
	Website          string       `yaml:"website"`
	Github           string       `yaml:"github"`
	Education        []Education  `yaml:"education"`
	Experience       []Experience `yaml:"experience"`
	Skills           []string     `yaml:"skills"`
	Projects         []Project    `yaml:"projects"`
	CompanyToApplyTo string       `yaml:"company_to_apply_to"`
	RoleToApplyTo    string       `yaml:"role_to_apply_to"`
}

// Education represents a single educational entry.
type Education struct {
	Institution string `yaml:"institution"`
	Degree      string `yaml:"degree"`
	StartDate   string `yaml:"start_date"`
	EndDate     string `yaml:"end_date"`
	GPA         string `yaml:"gpa"`
}

// Experience represents a single work experience entry.
type Experience struct {
	Company          string   `yaml:"company"`
	Position         string   `yaml:"position"`
	StartDate        string   `yaml:"start_date"`
	EndDate          string   `yaml:"end_date"`
	Responsibilities []string `yaml:"responsibilities"`
}

// Project represents a single project entry.
type Project struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	URL         string `yaml:"url"`
}

// Config handles loading of the resume configuration
type Config struct {
	Resume Resume
}

// application main directory (not persisted to YAML by default)
var mainDir string

// InitMainDir initializes the app main directory using environment or current working dir.
// Environment variable: COVLET_HOME
func InitMainDir() {
    if md := os.Getenv("COVLET_HOME"); md != "" {
        if abs, err := filepath.Abs(md); err == nil {
            if info, err2 := os.Stat(abs); err2 == nil && info.IsDir() {
                mainDir = abs
                return
            }
        }
    }
    // Default to a directory in the user's home: ~/.local/share/covlet
    if home, err := os.UserHomeDir(); err == nil {
        def := filepath.Join(home, ".local", "share", "covlet")
        _ = os.MkdirAll(def, 0o755)
        mainDir = def
        return
    }
    // Fallback to current working directory
    if cwd, err := os.Getwd(); err == nil {
        mainDir = cwd
    }
}

// SetMainDir sets the application's main directory; must exist and be a directory.
func SetMainDir(dir string) error {
	if dir == "" {
		return errors.New("main directory cannot be empty")
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return fmt.Errorf("path not accessible: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", abs)
	}
	mainDir = abs
	return nil
}

// GetMainDir returns the current main directory, initializing it if needed.
func GetMainDir() string {
	if mainDir == "" {
		InitMainDir()
	}
	return mainDir
}

// TemplatesDir returns the path to the templates root inside the main dir.
func TemplatesDir() string {
    return filepath.Join(GetMainDir(), "templates")
}

// EnsureTemplatesDir makes sure the templates directory exists.
func EnsureTemplatesDir() (string, error) {
    dir := TemplatesDir()
	if fi, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
				return dir, mkErr
			}
			return dir, nil
		}
		return dir, err
	} else if !fi.IsDir() {
		return dir, fmt.Errorf("templates path exists but is not a directory: %s", dir)
	}
	return dir, nil
}

// ValuesDir returns the directory where values (YAML) can be stored.
func ValuesDir() string {
    return filepath.Join(GetMainDir(), "values")
}

// EnsureValuesDir creates the values directory if needed.
func EnsureValuesDir() (string, error) {
    dir := ValuesDir()
    if fi, err := os.Stat(dir); err != nil {
        if os.IsNotExist(err) {
            if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
                return dir, mkErr
            }
            return dir, nil
        }
        return dir, err
    } else if !fi.IsDir() {
        return dir, fmt.Errorf("values path exists but is not a directory: %s", dir)
    }
    return dir, nil
}

// EnsureDownloadsCovletDir returns the default output directory for exported PDFs.
// On Linux it will be: ~/Downloads/covlet
func EnsureDownloadsCovletDir() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    dir := filepath.Join(home, "Downloads", "covlet")
    if err := os.MkdirAll(dir, 0o755); err != nil {
        return "", err
    }
    return dir, nil
}

// TODO: figure out why URLs don't work for the template generator

// LoadConfig reads and parses the YAML configuration file
func LoadConfig(filename string) (*Config, error) {
	yamlFile, err := os.ReadFile(filename)

	// if file does not exist, run through a setup and create file
	// ask for user input to create the directory and file
	if err != nil {
		return nil, err
	}

	var resume Resume
	if err = yaml.Unmarshal(yamlFile, &resume); err != nil {
		return nil, err
	}

	return &Config{Resume: resume}, nil
}

func (c *Config) SaveConfig(filename string) error {
	return nil
}

func (c *Config) CreateConfig() error {
	return nil
}
