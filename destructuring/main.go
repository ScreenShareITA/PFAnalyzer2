package destructuring

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// AnalyzeContent analizza il contenuto decompresso di un file prefetch
func AnalyzeContent(buffer []byte, extName string, extHash string, modifyTime int64, createTime int64) error {
	reader := bytes.NewReader(buffer)
	header := &Header{}
	if err := binary.Read(reader, binary.LittleEndian, header); err != nil {
		return err // non so se fare print
	}
	if header.Version < 30 {
		return fmt.Errorf("version not supported (%d)", header.Version)
	}
	/*


		_______________________________________________

		  _    _ ______          _____  ______ _____
		 | |  | |  ____|   /\   |  __ \|  ____|  __ \
		 | |__| | |__     /  \  | |  | | |__  | |__) |
		 |  __  |  __|   / /\ \ | |  | |  __| |  _  /
		 | |  | | |____ / ____ \| |__| | |____| | \ \
		 |_|  |_|______/_/    \_\_____/|______|_|  \_\


		_______________________________________________
	*/

	headerInfo := &HeaderInfo{} // not needed could be put in the processor
	headerInfo.ProcessName = decodeUTF16(header.ProcessName[:])
	headerInfo.Hash = fmt.Sprintf("%08X", header.Hash)
	if header.FMAOffset == 304 {
		headerInfo.PrefetchingFormat = "Windows 10"
	} else if header.FMAOffset == 296 {
		headerInfo.PrefetchingFormat = "Windows 10 / Windows 11"
	} else {
		headerInfo.PrefetchingFormat = "?? Windows 10 ??"
	}

	PrintTitle("Header Useful Info")
	fmt.Printf("Version: %d (%s)\n", header.Version, headerInfo.PrefetchingFormat)
	fmt.Printf("Process Name: %s\n", headerInfo.ProcessName)
	fmt.Printf("Hash: %s\n", headerInfo.Hash)
	/*


		__________________________________________


		  _____  ______ _____ _   _ ______ ____
		 |  __ \|  ____|_   _| \ | |  ____/ __ \
		 | |__) | |__    | | |  \| | |__ | |  | |
		 |  ___/|  __|   | | | . ` |  __|| |  | |
		 | |    | |     _| |_| |\  | |   | |__| |
		 |_|    |_|    |_____|_| \_|_|    \____/


		__________________________________________
	*/

	pfInfo1 := &PrefetchInfo1{}
	if err := ReadData(reader, pfInfo1); err != nil {
		return err
	}

	var Remnant2 uint64 //if version is 30.1 with also remnant2
	if header.FMAOffset == 304 {
		if err := ReadData(reader, &Remnant2); err != nil {
			return err
		}
	}
	pfInfo2 := &PrefetchInfo2{}
	if err := ReadData(reader, pfInfo2); err != nil {
		return err
	}

	if _, err := reader.Seek(int64(pfInfo2.AppNameOffset), 0); err != nil { // posiziono sul appname
		return err
	}
	appNameBytes := make([]byte, pfInfo2.AppNameSize)    // creo bytearr
	if _, err := reader.Read(appNameBytes); err != nil { // leggo appname
		return err
	}
	//reader.Seek(int64(header.FMAOffset), 0) // ritorno sulla posizione originale (NOT NEEDED)

	PrintTitle("Useful PFInfo")
	var dates uint32 = 0
	// cambiare tutta la logica con 7 cicli previousrun ecc..
	lasttime := pfInfo1.RunTimesB[0]
	for i := range pfInfo1.RunTimesB {
		RunTime := pfInfo1.RunTimesB[i]
		if RunTime != 0 {
			if lasttime < RunTime {
				fmt.Printf("ANORMAL TIMESTAMP ORDER (%d > %d)\n", RunTime, lasttime)
			}
			dates++
			fmt.Printf("RunTime%d: %s\n", i+1, uint64ToStringDate(RunTime))
			lasttime = RunTime
		}
	}
	splitpoint := len(appNameBytes) - 2
	fmt.Printf("RunCount: %d\n", pfInfo2.RunCount)
	fmt.Printf("Dates: %d/8\n", dates)
	fmt.Printf("AppName: %s\n", appNameBytes[:splitpoint])

	/*


		________________________________________________________________________


		 __      ______  _     _    _ __  __ ______   _____ _   _ ______ ____
		 \ \    / / __ \| |   | |  | |  \/  |  ____| |_   _| \ | |  ____/ __ \
		  \ \  / / |  | | |   | |  | | \  / | |__      | | |  \| | |__ | |  | |
		   \ \/ /| |  | | |   | |  | | |\/| |  __|     | | | . ` |  __|| |  | |
		    \  / | |__| | |___| |__| | |  | | |____   _| |_| |\  | |   | |__| |
		     \/   \____/|______\____/|_|  |_|______| |_____|_| \_|_|    \____/



		________________________________________________________________________
	*/
	_, err := reader.Seek(int64(header.VIOffset), 0)
	if err != nil {
		return err
	}

	VolumeInfo := &VolumeInfo{}
	if err := ReadData(reader, VolumeInfo); err != nil {
		return err
	}
	VolumeInfo.DevicePathOffset += header.VIOffset
	DevicePathSize := VolumeInfo.DevicePathNOC * 2
	DevicePathEnd := VolumeInfo.DevicePathOffset + DevicePathSize

	VolumeInfo.FileRefsOffset += header.VIOffset
	//FileRefsEnd :=  VolumeInfo.FileRefsOffset + VolumeInfo.FileRefsSize

	VolumeInfo.DirStringsOffset += header.VIOffset
	//DirStringsEnd := VolumeInfo.DevicePathOffset + VolumeInfo.DirStringsSize
	/*FileRefOffset:=header.VIOffset+VolumeInfo.FileRefsOffset
	DirStringsOffset:=header.VIOffset+VolumeInfo.DirStringsOffset*/

	PrintTitle("VolumeInfo Useful")
	fmt.Printf("Volume Creation Time: %s\n", uint64ToStringDate(VolumeInfo.VolumeCreationTime))
	fmt.Printf("Volume Serial Number: %8X\n", VolumeInfo.VolumeSerialNumber)
	fmt.Printf("Volume Device Path: %s\n", decodeUTF16(buffer[VolumeInfo.DevicePathOffset:DevicePathEnd]))

	_, err = reader.Seek(int64(VolumeInfo.DirStringsOffset), 0)
	if err != nil {
		return err
	}
	PrintTitle(fmt.Sprintf("Referenced Directories (%d)", VolumeInfo.DirStringsCount))
	for i := 0; i < int(VolumeInfo.DirStringsCount); i++ {
		var NOC uint16
		if err := ReadData(reader, &NOC); err != nil {
			return err
		}
		size := NOC*2 + 2
		stringBytes := make([]byte, size)
		if err := ReadData(reader, stringBytes); err != nil {
			return err
		}
		dirstr := decodeUTF16(stringBytes)
		fmt.Printf("%02d - %s\n", i+1, dirstr)
	}

	/*


		____________________________________________________


		  __  __  ____  _____  _    _ _      ______  _____
		 |  \/  |/ __ \|  __ \| |  | | |    |  ____|/ ____|
		 | \  / | |  | | |  | | |  | | |    | |__  | (___
		 | |\/| | |  | | |  | | |  | | |    |  __|  \___ \
		 | |  | | |__| | |__| | |__| | |____| |____ ____) |
		 |_|  |_|\____/|_____/ \____/|______|______|_____/



		____________________________________________________
	*/
	_, err = reader.Seek(int64(header.FMAOffset), 0)
	if err != nil {
		return err
	}

	PrintTitle(fmt.Sprintf("Modules Filenames (%d)", header.FMAEntries))
	var lastOffset uint32 = header.FNOffset
	for i := uint32(0); i < header.FMAEntries; i++ {
		fmaEntry := &FMA{}
		if err := ReadData(reader, fmaEntry); err != nil {
			return err
		}
		fnOffset := header.FNOffset + fmaEntry.FileNameOffset
		fnEnd := fnOffset + fmaEntry.FileNameNOC*2
		if lastOffset < fnOffset {
			fmt.Printf("Hole Found (%d):\n%s\n", fnOffset-lastOffset, buffer[lastOffset:fnOffset])
		}
		fmt.Printf("%02d - %s\n",
			i+1,
			decodeUTF16(buffer[fnOffset:fnEnd]))
		lastOffset = fnEnd + 2
	}

	return nil
}
