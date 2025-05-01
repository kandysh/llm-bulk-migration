package prompts

import (
	"go.uber.org/zap"
	"llm-bulk-migration/internal"
	"log"
	"sync"
)

// resourceCache provides a singleton cache for static resources
// to avoid loading them repeatedly across function calls
type resourceCache struct {
	rtlExamples string
	once        sync.Once
}

var cache = &resourceCache{}

// MigratePrompt generates a prompt for migrating Enzyme tests to React Testing Library.
// It loads the necessary source code and examples, then formats them into a prompt array.
// The sourcePath parameter should point to the test file that needs migration.
func MigratePrompt(sourcePath string, rtlExamplePath string) []string {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	sugar := logger.Sugar()
	defer logger.Sync()

	sugar.Info("Generating prompt for %s", sourcePath)

	cache.once.Do(func() {
		sugar.Info("Loading example RTL files")
		cache.rtlExamples = internal.FilesFromPath(rtlExamplePath)
		sugar.Info("Example RTL files loaded")
	})

	sugar.Info("Loading sibling test files")
	siblingTests := internal.FilesFromPath(sourcePath, true)
	sugar.Info("Sibling test files loaded")

	sugar.Info("Loading test file")
	testFile := internal.FileFromPath(sourcePath)
	sugar.Info("Test file loaded")

	sugar.Info("Extracting imports")
	imports, err := internal.ExtractImportsContent(sourcePath)

	if err != nil {
		sugar.Warn("Warning: failed to extract imports from %s: %v", sourcePath, err)
		imports = "// No imports could be extracted"
	} else {
		sugar.Info("Imports extracted")
	}

	prompt := []string{
		"Convert this Enzyme test to React Testing Library:",
		"SIBLING TESTS:\n" + siblingTests,
		"RTL EXAMPLES:\n" + cache.rtlExamples,
		"IMPORTS:\n" + imports,
		"TEST TO MIGRATE:\n" + testFile,
	}
	sugar.Info("Prompt generated")
	return prompt
}
