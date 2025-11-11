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

// Config handles loading of the resume configuration
type Config struct {
	Resume Resume
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
