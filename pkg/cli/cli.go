package cli

import (
	"cover-letter-templates/pkg/config"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"text/template"
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
			configFile, err := config.LoadConfig("config.yml")
			if err != nil {
				return fmt.Errorf("error loading config: %v", err)
			}

			if company := cCtx.String("company"); company != "" {
				configFile.Resume.CompanyToApplyTo = company
			}
			if position := cCtx.String("position"); position != "" {
				configFile.Resume.RoleToApplyTo = position
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
			err = t.Execute(os.Stdout, configFile.Resume)
			if err != nil {
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
