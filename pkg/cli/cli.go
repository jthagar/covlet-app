package cli

import (
	"bufio"
	"cover-letter-templates/pkg/config"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"
)

func Run() error {
	app := &cli.App{
		Name:  "cover-letter",
		Usage: "Generate a cover letter from template",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "company",
				Aliases:  []string{"c"},
				Usage:    "Company to apply to",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "manager",
				Aliases:  []string{"m"},
				Usage:    "Hiring manager's name",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "position",
				Aliases:  []string{"p"},
				Usage:    "Position to apply for",
				Required: false,
			},
		},
		Action: func(cCtx *cli.Context) error {
			const cfgPath = "config.yml"

			// Try to load config; handle first-run initialization.
			cfg, err := config.LoadConfig(cfgPath)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("error loading config: %v", err)
				}
				// Config doesn't exist yet
				cfg = &config.Config{}
			}

			// If home directory not set or missing on disk, prompt user to initialize.
			if cfg.HomeDir == "" || !dirExists(cfg.HomeDir) {
				reader := bufio.NewReader(os.Stdin)
				userHome, _ := os.UserHomeDir()
				defaultHome := filepath.Join(userHome, ".cover-letter-templates")
				fmt.Printf("Enter a home directory for the app (default: %s): ", defaultHome)
				in, _ := reader.ReadString('\n')
				in = strings.TrimSpace(in)
				if in == "" {
					in = defaultHome
				}
				// Expand leading ~ if provided
				if strings.HasPrefix(in, "~") {
					in = filepath.Join(userHome, strings.TrimPrefix(in, "~"))
				}
				in = filepath.Clean(in)

				if err := os.MkdirAll(in, 0755); err != nil {
					return fmt.Errorf("failed to create home directory: %v", err)
				}
				cfg.HomeDir = in
				if err := cfg.SaveConfig(cfgPath); err != nil {
					return fmt.Errorf("failed to save config: %v", err)
				}
				fmt.Printf("Initialized app. Config saved to %s with home_dir=%s\n", cfgPath, cfg.HomeDir)
			}

			// Allow flags to override certain resume fields for generation time.
			if company := cCtx.String("company"); company != "" {
				cfg.Resume.CompanyToApplyTo = company
			}
			if position := cCtx.String("position"); position != "" {
				cfg.Resume.RoleToApplyTo = position
			}

			templateFile, err := os.ReadFile("templates/base/cover_letter.tpl")
			if err != nil {
				return fmt.Errorf("error reading template file: %v", err)
			}

			t, err := template.New("cover_letter").Parse(string(templateFile))
			if err != nil {
				return fmt.Errorf("error parsing template: %v", err)
			}

			fmt.Println("--- Generated Cover Letter ---")
			if err = t.Execute(os.Stdout, cfg.Resume); err != nil {
				return fmt.Errorf("error executing template: %v", err)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	return nil
}

func dirExists(path string) bool {
	if path == "" {
		return false
	}
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}
