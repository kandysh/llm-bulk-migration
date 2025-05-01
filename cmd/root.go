package cmd

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"llm-bulk-migration/prompts"
	"os"
	"strings"
)

var validSteps = map[string]bool{
	"fix-enzyme": true,
	"fix-jest":   true,
	"fix-eslint": true,
	"fix-tsc":    true,
	"":           true,
}

var (
	matchPath      string
	stepType       string
	rtlExamplePath string
)

var rootCmd = &cli.Command{
	Name:    "llm-bulk-migration",
	Usage:   "Bulk migration tool using Large Language Models",
	Version: "0.1.0",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "match",
			Aliases:     []string{"m"},
			Usage:       "Specify file or directory path to process",
			Destination: &matchPath,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "step",
			Aliases:     []string{"s"},
			Usage:       "Migration step to perform: fix-enzyme, fix-jest, fix-eslint, fix-tsc",
			Destination: &stepType,
			Required:    false,
			Value:       "",
		},
		&cli.StringFlag{
			Name:        "rtl-examples",
			Aliases:     []string{"e"},
			Usage:       "Path to directory containing RTL examples",
			Destination: &rtlExamplePath,
			Required:    false,
			Value:       "resources/rtl-examples",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		if !validSteps[stepType] {
			return fmt.Errorf("invalid step type: %s. Must be one of: fix-enzyme, fix-jest, fix-eslint, fix-tsc", stepType)
		}

		promptText := strings.Join(prompts.MigratePrompt(matchPath, rtlExamplePath), "\n")
		fmt.Println(promptText)

		return nil
	},
}

func Execute() {
	err := rootCmd.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
