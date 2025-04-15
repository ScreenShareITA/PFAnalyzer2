package decompressing

import (
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
	"syscall"

	"pfanalyzer2/utils"
)

func AnalyzePath(baseName string) (string, string) {
	re := regexp.MustCompile(`^(.{1,29})-([A-F0-9]{8})`)
	matches := re.FindStringSubmatch(baseName)

	if len(matches) == 3 { // start with filename analysis
		fmt.Printf("Process Name: %s\nHash: %s\n", matches[1], matches[2])
		return matches[1], matches[2]
	}
	fmt.Println("Process Name: INVALID FILENAME!\nHash: INVALID FILENAME!")
	return "", ""
}

func AnalyzeFile(info os.FileInfo) (int64, int64) {
	sys := info.Sys().(*syscall.Win32FileAttributeData)
	//needed for time confront
	accessTime := sys.LastAccessTime.Nanoseconds()
	modifyTime := sys.LastWriteTime.Nanoseconds()
	createTime := sys.CreationTime.Nanoseconds()
	fmt.Printf(
		"ReadOnly: %v\nHidden: %v\nCreation Time: %v\nModify Time: %v\nAccess Time: %v\nModify-Access Time Match: %v\n",
		(sys.FileAttributes&syscall.FILE_ATTRIBUTE_READONLY) != 0,
		(sys.FileAttributes&syscall.FILE_ATTRIBUTE_HIDDEN) != 0,
		utils.IntUnixToStringDate(createTime),
		utils.IntUnixToStringDate(modifyTime),
		utils.IntUnixToStringDate(accessTime),
		accessTime == modifyTime,
	)
	return modifyTime, createTime
}

func DecompressAnalysis(file *os.File, fileSize int64) ([]byte, error) {
	uncompressedSize, err := checkPrefetchHeader(file)
	if err != nil {
		return nil, err
	}
	var fileData []byte

	if uncompressedSize > 0 { // decompress
		compressedData := make([]byte, fileSize-8)
		if _, err = file.Read(compressedData); err != nil {
			return nil, fmt.Errorf("error reading compressed data: %v", err)
		}

		if fileData, err = DecompressBuffer(compressedData, uncompressedSize); err != nil {
			return nil, fmt.Errorf("error during decompression: %v", err)
		}
		fmt.Printf("Compressed: true\nUncompressed Size: %d\n", uncompressedSize)
		return fileData, nil
	}

	// file is not compressed
	fileData = make([]byte, fileSize)
	if _, err = file.Read(fileData); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	fmt.Println("Compressed: false")
	return fileData, nil
}

func checkPrefetchHeader(file *os.File) (uint32, error) {
	header := make([]byte, 8)
	_, err := file.Read(header)
	if err != nil {
		return 0, fmt.Errorf("error while reading file header: %v", err)
	}

	if header[0] == 'M' && header[1] == 'A' && header[2] == 'M' && header[3] == 0x04 { // file is compressed
		return binary.LittleEndian.Uint32(header[4:]), nil
	}
	if header[4] == 'S' && header[5] == 'C' && header[6] == 'C' && header[7] == 'A' { // file is not compressed
		_, err = file.Seek(0, 0) // pointer back to 0
		if err != nil {
			return 0, fmt.Errorf("error seeking file: %v", err)
		}
		return 0, nil
	}

	return 0, fmt.Errorf("the specified file is not a prefetch file or is corrupted")
}
