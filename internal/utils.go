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
	// maxRecursionDepth prevents stack overflows by limiting directory traversal depth
	maxRecursionDepth = 100
	// initialBufferSize sets a 4KB initial allocation for strings.Builder
	initialBufferSize = 4096
	// recursionTimeout sets the maximum duration for recursive file operations
	recursionTimeout = 30 * time.Second
)

// FileFromPath reads a single file and returns its contents as a string
// with the filename as a header.
//
// If the path doesn't exist or points to a directory, an error message
// is returned as a string. The returned string format is:
//
//	// filename.ext
//	<file contents>
//
// Parameters:
//   - filePath: the full path to the file to be read
//
// Returns:
//   - a string containing the file contents with header, or an error message
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
	WriteFileToBuilder(&result, string(content), filePath, fileInfo.Size())

	return result.String()
}

// FilesFromPath reads all files from the given path and returns their contents
// as a single formatted string. It handles both single files and directories.
//
// For directories, it recursively traverses subdirectories (up to maxRecursionDepth)
// and combines all file contents with appropriate headers.
//
// Parameters:
//   - sourcePath: path to a file or directory
//   - ignoreSourceFileOpt: optional boolean flag; if true and sourcePath are a file,
//     that file will be excluded from directory contents. When sourcePath is a
//     specific file and ignoreSourceFile is true, returns "No files found".
//
// Returns:
//   - a string containing all file contents with headers, or an error message
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
// files at the given path with tracking for recursion depth.
//
// It checks for context timeout, respects maximum recursion depth, and handles
// both directories and files by either recursing into subdirectories or reading
// and formatting file contents.
//
// Parameters:
//   - ctx: context for cancellation and timeout control
//   - sourcePath: path to process (file or directory)
//   - ignoreSourceFile: if true, skips the original source file
//   - depth: current recursion depth to track and limit traversal
//
// Returns:
//   - a string containing all processed file contents
//   - error if any occurred during processing
func filesFromPathWithDepth(ctx context.Context, sourcePath string, ignoreSourceFile bool, depth int) (string, error) {
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("operation timed out or was canceled")
	default:
	}

	if depth > maxRecursionDepth {
		return "", fmt.Errorf("maximum recursion depth reached at: %s", sourcePath)
	}

	allInfo := strings.Builder{}
	allInfo.Grow(initialBufferSize)

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
				allInfo.WriteString("\n// Error processing directory ")
				allInfo.WriteString(currentFileName)
				allInfo.WriteString(": ")
				allInfo.WriteString(err.Error())
				allInfo.WriteString("\n")
				continue
			}

			if subDirContents != "No files found" &&
				subDirContents != "No file contents found" {
				allInfo.WriteString("\n// Directory: ")
				allInfo.WriteString(currentFileName)
				allInfo.WriteString("\n")
				allInfo.WriteString(subDirContents)
				filesProcessed++
			}
			continue
		}

		if ignoreSourceFile && sourceFileName != "" && currentFileName == sourceFileName {
			continue
		}

		fileContent := FileFromPath(filePath)
		tempInfo, _ := os.Stat(filePath)
		WriteFileToBuilder(&allInfo, fileContent, filePath, tempInfo.Size())
		filesProcessed++
	}

	if filesProcessed == 0 {
		return "No file contents found", nil
	}

	return allInfo.String(), nil
}

func WriteFileToBuilder(builder *strings.Builder, fileContent string, fileName string, fileSize int64) {
	headerSize := len(fileName) + 4

	builder.Grow(int(fileSize) + headerSize)
	if strings.HasPrefix(fileContent, "Error:") || strings.HasPrefix(fileContent, "Error reading file") {

		return
	}
	builder.WriteString("\n")
	builder.WriteString("// " + fileName + "\n")
	builder.WriteString(fileContent)
	builder.WriteString("\n")

}

// ExtractImportsContent extracts and retrieves content from all relative imports in a TSX test file
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

		// Handle imports that don't specify file extension
		fullPath := filepath.Join(baseDir, relativePath)
		if !strings.HasSuffix(relativePath, ".tsx") && !strings.HasSuffix(relativePath, ".ts") {
			// Try common extensions
			for _, ext := range []string{".tsx", ".ts", ".js", ".jsx"} {
				if _, err := os.Stat(fullPath + ext); err == nil {
					fullPath = fullPath + ext
					break
				}

				// Check for index files in directories
				indexPath := filepath.Join(fullPath, "index"+ext)
				if _, err := os.Stat(indexPath); err == nil {
					fullPath = indexPath
					break
				}
			}
		}

		// Read the file content
		fileContent, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}
		tempInfo, _ := os.Stat(fullPath)
		WriteFileToBuilder(&importContents, string(fileContent), fullPath, tempInfo.Size())
	}

	return importContents.String(), nil
}
