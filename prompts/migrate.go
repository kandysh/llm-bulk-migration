package prompts

import (
	"llm-bulk-migration/internal"
)

var siblingTestFilesSourceCode string
var reactTestingLibraryExamples string
var nearestImportSourceCode string
var componentFileSourceCode = ""
var testFileSourceCode = ""

func MigratePrompt(sourcePath string) []string {
	siblingTestFilesSourceCode = internal.FilesFromPath(sourcePath, true)
	reactTestingLibraryExamples = internal.FilesFromPath("./resources/rtl-examples")
	nearestImportSourceCode, _ = internal.ExtractImportsContent(sourcePath)
	//var prompt = []string{
	//	"Convert this Enzyme test to React Testing Library:",
	//	"SIBLING TESTS:\n" + siblingTestFilesSourceCode,
	//	"RTL EXAMPLES:\n" + reactTestingLibraryExamples,
	//	"IMPORTS:\n" + nearestImportSourceCode,
	//	"COMPONENT SOURCE:\n" + componentFileSourceCode,
	//	"TEST TO MIGRATE:\n" + testFileSourceCode,
	//}
	return []string{nearestImportSourceCode}
}
