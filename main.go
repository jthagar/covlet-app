package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"text/template"
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

func main() {
	// Read the resume data from the YAML file.
	yamlFile, err := os.ReadFile("resume.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Parse the YAML data into our Resume struct.
	var resume Resume
	err = yaml.Unmarshal(yamlFile, &resume)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	// Read the cover letter template.
	templateFile, err := os.ReadFile("cover_letter.tmpl")
	if err != nil {
		log.Fatalf("Error reading template file: %v", err)
	}

	// Create a new template and parse the template file content.
	t, err := template.New("cover_letter").Parse(string(templateFile))
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// Execute the template with the resume data and print the output to the console.
	fmt.Println("--- Generated Cover Letter ---")
	err = t.Execute(os.Stdout, resume)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}
