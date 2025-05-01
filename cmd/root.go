package cmd

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"llm-bulk-migration/prompts"
	"os"
	"strings"
)

var rootCmd = &cli.Command{
	Name:    "llm-bulk-migration",
	Usage:   "Bulk migration tool using Large Language Models",
	Version: "0.0.1",
	Action: func(context.Context, *cli.Command) error {
		fmt.Println("Welcome to the LLM Bulk Migration tool!")
		fmt.Println("Use --help to see the available commands")
		return nil
	},
}

func Execute() {
	err := rootCmd.Run(context.Background(), os.Args)
	fmt.Println(strings.Join(prompts.MigratePrompt("./resources/__tests__/component2.test.tsx"), "\n"))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
