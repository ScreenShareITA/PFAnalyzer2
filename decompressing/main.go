package decompressing

import (
	"fmt"
	"os"
)

// ProcessFile elabora un file prefetch e mostra le informazioni estratte
func ProcessFile(file *os.File, info os.FileInfo) (string, string, int64, int64, []byte, error) {
	fmt.Printf("\033[91mExternal File Info:\033[0m\n")
	procName, Hash := AnalyzePath(info.Name())
	modifyTime, createTime := AnalyzeFile(info)
	byteArr, err := DecompressAnalysis(file, info.Size())

	// returns decompression error directly
	return procName, Hash, modifyTime, createTime, byteArr, err
}
