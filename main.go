package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pfanalyzer2/decompressing"
	"pfanalyzer2/destructuring"
)

func main() {
	flag.Parse()

	// Verifica che sia stato fornito il percorso del file
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: PFAnalyzer2 <filepath>")
		os.Exit(1)
	}

	filePath := args[0]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error: Path '%s' does not exist\n", filePath)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		fmt.Printf("Error while gathering info for %s \n", filePath)
		os.Exit(1)
	}

	// Process single file or directory
	if info.IsDir() {
		// Process all .pf files in directory
		fmt.Printf("\033[92mProcessing all .pf files in %s:\033[0m\n\n\n", filePath)
		err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Analysis error for %s: %v\n", path, err)
				return nil
			}
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".pf") { //change the condition to match MAM signature or SCCA
				// Apri un nuovo file handle per ogni file .pf
				currentFile, err := os.Open(path)
				if err != nil {
					fmt.Printf("Error opening file '%s': %v\n", path, err)
					return nil
				}
				defer currentFile.Close()

				processFile(path, currentFile, info)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error walking directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Process single file
		processFile(filePath, file, info)
	}
	fmt.Print("\033[33mPress ENTER to exit...\033[0m")
	fmt.Scanln()
}

func processFile(filePath string, file *os.File, info os.FileInfo) {
	fmt.Printf("\033[94mProcessing: %s\033[0m\n\n", filePath)
	extName, extHash, modifyTime, createTime, bytearr, err := decompressing.ProcessFile(file, info)
	if err != nil {
		fmt.Printf("Error processing file '%s': %v\n", info.Name(), err)
		return
	}
	err = destructuring.AnalyzeContent(bytearr, extName, extHash, modifyTime, createTime)
	if err != nil {
		fmt.Printf("Error analyzing content of file '%s': %v\n", info.Name(), err)
		return
	}
	fmt.Print("\n\n\n")
}
