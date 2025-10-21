// +build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	oldImport := "github.com/onetakesolutions/onetake-corpsite-backend"
	newImport := "github.com/ChrisChantszto/futurefrontier"
	
	count := 0
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories and non-Go files
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		
		// Skip this file itself
		if strings.Contains(path, "fix_imports.go") {
			return nil
		}
		
		// Read file
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		
		// Replace imports
		newContent := strings.ReplaceAll(string(content), oldImport, newImport)
		
		// Only write if changed
		if string(content) != newContent {
			err = os.WriteFile(path, []byte(newContent), info.Mode())
			if err != nil {
				return err
			}
			fmt.Printf("✓ Fixed: %s\n", path)
			count++
		}
		
		return nil
	})
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("\n✓ Fixed %d files\n", count)
	fmt.Println("\nNow run: go mod tidy")
}
