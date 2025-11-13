package config

import (
	"os"

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

// Config represents the full application configuration file
// home_dir is stored at the top level, while Resume fields are inlined at the root as well.
type Config struct {
	HomeDir string `yaml:"home_dir"`
	Resume  Resume `yaml:",inline"`
}

// LoadConfig reads and parses the YAML configuration file
func LoadConfig(filename string) (*Config, error) {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig writes the configuration back to the given YAML file.
func (c *Config) SaveConfig(filename string) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0644)
}

func (c *Config) CreateConfig() error {
	return nil
}
