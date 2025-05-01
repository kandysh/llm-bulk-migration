// Package internal provides utility functions for file handling and processing.
// It includes functionality for reading single files and recursively gathering
// content from directories with proper error handling and timeout mechanisms.
package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	maxRecursionDepth = 100              // Limits directory traversal depth to prevent stack overflows
	initialBufferSize = 4096             // 4KB initial allocation for strings.Builder
	recursionTimeout  = 30 * time.Second // Maximum duration for recursive file operations
)

// FileFromPath reads a single file and returns its contents as a string
// with the filename as a header. If the path doesn't exist or points to a directory,
// an error message is returned as a string.
func FileFromPath(filePath string) string {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Sprintf("Error accessing file: %v", err)
	}

	if fileInfo.IsDir() {
		return fmt.Sprintf("Error: %s is a directory, not a file", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("Error reading file %s: %v", filePath, err)
	}

	var result strings.Builder
	result.Grow(len(content) + len(filePath) + 5)
	WriteFileToBuilder(&result, string(content), filePath)

	return result.String()
}

// FilesFromPath reads all files from the given path and returns their contents
// as a single formatted string. It handles both single files and directories.
//
// For directories, it recursively traverses subdirectories (up to maxRecursionDepth)
// and combines all file contents with appropriate headers. If ignoreSourceFileOpt
// is true and sourcePath is a file, that file will be excluded from directory contents.
func FilesFromPath(sourcePath string, ignoreSourceFileOpt ...bool) string {
	ignoreSourceFile := false
	if len(ignoreSourceFileOpt) > 0 {
		ignoreSourceFile = ignoreSourceFileOpt[0]
	}

	ctx, cancel := context.WithTimeout(context.Background(), recursionTimeout)
	defer cancel()

	result, err := filesFromPathWithDepth(ctx, sourcePath, ignoreSourceFile, 0)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return result
}

// filesFromPathWithDepth is an internal recursive helper function that processes
// files at the given path with tracking for recursion depth. It checks for context
// timeout, respects maximum recursion depth, and handles both directories and files.
func filesFromPathWithDepth(ctx context.Context, sourcePath string, ignoreSourceFile bool, depth int) (string, error) {
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("operation timed out or was canceled")
	default:
	}

	if depth > maxRecursionDepth {
		return "", fmt.Errorf("maximum recursion depth reached at: %s", sourcePath)
	}

	fileInfo, err := os.Stat(sourcePath)
	if err != nil {
		return "", fmt.Errorf("error accessing path %s: %w", sourcePath, err)
	}

	var dir string
	var sourceFileName string

	if fileInfo.IsDir() {
		dir = sourcePath
	} else {
		dir = filepath.Dir(sourcePath)
		sourceFileName = filepath.Base(sourcePath)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("error reading directory %s: %w", dir, err)
	}

	if len(files) == 0 {
		return "No files found", nil
	}

	allInfo := strings.Builder{}
	allInfo.Grow(initialBufferSize)
	filesProcessed := 0

	for _, fileEntry := range files {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("operation timed out or was canceled")
		default:
		}

		currentFileName := fileEntry.Name()
		filePath := filepath.Join(dir, currentFileName)

		if fileEntry.IsDir() {
			subDirContents, err := filesFromPathWithDepth(ctx, filePath, ignoreSourceFile, depth+1)
			if err != nil {
				fmt.Fprintf(&allInfo, "\n// Error processing directory %s: %s\n",
					currentFileName, err.Error())
				continue
			}

			if subDirContents != "No files found" && subDirContents != "No file contents found" {
				fmt.Fprintf(&allInfo, "\n// Directory: %s\n%s", currentFileName, subDirContents)
				filesProcessed++
			}
			continue
		}

		if ignoreSourceFile && sourceFileName != "" && currentFileName == sourceFileName {
			continue
		}

		fileContent := FileFromPath(filePath)
		allInfo.WriteString(fileContent)
		allInfo.WriteString("\n")
		filesProcessed++
	}

	if filesProcessed == 0 {
		return "No file contents found", nil
	}

	return allInfo.String(), nil
}

// WriteFileToBuilder writes file content with a header to the provided strings.Builder.
func WriteFileToBuilder(builder *strings.Builder, fileContent string, fileName string) {
	if strings.HasPrefix(fileContent, "Error:") || strings.HasPrefix(fileContent, "Error reading file") {
		return
	}

	builder.WriteString("// ")
	builder.WriteString(fileName)
	builder.WriteString("\n")
	builder.WriteString(fileContent)
	builder.WriteString("\n")
}

// ExtractImportsContent extracts and retrieves content from all relative imports
// in a TSX test file. It parses import statements, resolves relative paths, and
// returns the content of all imported files combined.
func ExtractImportsContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	baseDir := filepath.Dir(filePath)
	importRegex := regexp.MustCompile(`import\s+(?:.*\s+from\s+)?['"](\.[^'"]+)['"]`)
	matches := importRegex.FindAllSubmatch(content, -1)

	importContents := strings.Builder{}
	importContents.Grow(initialBufferSize)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		relativePath := string(match[1])
		fullPath := filepath.Join(baseDir, relativePath)

		if !strings.HasSuffix(relativePath, ".tsx") && !strings.HasSuffix(relativePath, ".ts") {
			fullPath = resolveImportPath(fullPath)
		}

		if fullPath == "" {
			continue
		}

		fileContent, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		importContents.Grow(len(fileContent) + len(filePath) + 5)
		WriteFileToBuilder(&importContents, string(fileContent), fullPath)
	}

	return importContents.String(), nil
}

// resolveImportPath attempts to resolve an import path by trying common file extensions
// and checking for index files in directories. Returns an empty string if resolution fails.
func resolveImportPath(basePath string) string {
	for _, ext := range []string{".tsx", ".ts", ".js", ".jsx"} {
		if _, err := os.Stat(basePath + ext); err == nil {
			return basePath + ext
		}

		indexPath := filepath.Join(basePath, "index"+ext)
		if _, err := os.Stat(indexPath); err == nil {
			return indexPath
		}
	}

	return ""
}
