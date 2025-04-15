package destructuring

// Header represents the prefetch file header
type Header struct {
	// Effective Header
	Version           uint32   // Format version (17, 23, 26, 30, 31)		WRITE
	Signature         [4]byte  // Signature string (SCCA) 					USELESS
	Unknown1          uint32   // Unknown value 						USELESS
	FileSize          uint32   // File size 							USELESS
	ProcessName       [60]byte // Executable name with EOS 					WRITE
	Hash              uint32   // Prefetch hash							WRITE
	X_PrefetchingType uint32   // X_PrefetchingType 					USELESS
	// Offset e entries
	FMAOffset  uint32 // File Metrics Array offset 						USEFUL
	FMAEntries uint32 // File Metrics Array entries 					USEFUL
	TCAOffset  uint32 // Trace Chains Array offset 						USEFUL
	TCAEntries uint32 // Trace Chains Array entries 					USEFUL
	FNOffset   uint32 // File Names offset 								USEFUL
	FNSize     uint32 // File Names size								USEFUL
	VIOffset   uint32 // Volume Information offset 						USEFUL
	NVolumes   uint32 // Number of volumes 							USEFUL
	VISize     uint32 // Volume Information size 						USEFUL
}

type HeaderInfo struct { //eliminare questo struct o commentarlo
	// Calculated fields
	PrefetchingFormat string // Format based on version 					WRITE
	ProcessName       string // Executable name without EOS 					WRITE
	Hash              string // Prefetch hash in hexadecimal 					WRITE
	FMAEntrySize      uint32 // Size of each FMA entry 							USELESS
	TCAEntrySize      uint32 // Size of each TCA entry 							USELESS

	FMASize uint32 // File Metrics Array size 							USELESS
	FMAEnd  uint32 // File Metrics Array end 							USELESS
	TCASize uint32 // Trace Chains Array size							USELESS
	TCAEnd  uint32 // Trace Chains Array end 							USELESS
	FNEnd   uint32 // File Names end 									USELESS
	VIEnd   uint32 // Volume Information end 							USELESS
	//RightSignature    bool   // Signature verification = SCCA 			USELESS

	// Verification fields
	//HashMismatch           bool // External hash verification 			USELESS
	//ExecutableNameMismatch bool // Comparison with process name			USELESS
	//ValidOffsets           bool // Verify correct distances 				USELESS
	//ValidEntries           bool // Verify FNEntries = FMAEntries 			USELESS
}

type PrefetchInfo1 struct {
	Unknown1_1 uint32    // 4 byte di dati sconosciuti
	Unknown1_2 uint32    // 4 byte di dati sconosciuti
	RunTimesB  [8]uint64 // 8x8=64 byte timestamp degli avvii
	Remnant1   uint64    // 8 byte di dati sconosciuti
}

//remnant 2 uint64
type PrefetchInfo2 struct {
	RunCount      uint32   // 4 byte contatore di esecuzioni
	Unknown2_1    uint32   // 4 byte di dati sconosciuti
	Unknown2_2    uint32   // 4 byte di dati sconosciuti
	AppNameOffset uint32   // 4 byte offset della stringa hash
	AppNameSize   uint32   // 4 byte dimensione della stringa hash
	EmptyValues   [76]byte // 76 byte di dati vuoti
}

type FMA struct { // version 23,26,30,31
	PFStartTime    uint32
	PFDuration     uint32
	PFAvgDuration  uint32
	FileNameOffset uint32
	FileNameNOC    uint32 // number of characters
	UnknownFlag    uint32
	MFTReference   MftRef
}

type TCA struct { // 30 31
	TotalBlockLoadCount uint32
	Unknown1            uint8
	Unknown2            uint8
	Unknown3            uint16
}

type VolumeInfo struct {
	DevicePathOffset   uint32   // Offset from start of volume information
	DevicePathNOC      uint32   // Number of characters in device path
	VolumeCreationTime uint64   // FILETIME timestamp
	VolumeSerialNumber uint32   // Volume serial number
	FileRefsOffset     uint32   // File references offset
	FileRefsSize       uint32   // File references data size
	DirStringsOffset   uint32   // Directory strings offset
	DirStringsCount    uint32   // Number of directory strings
	Unknown1           uint32   // Unknown (possible relation to file references remnant)
	EmptyValues1       [24]byte // Unknown empty values
	Unknown2           uint32   // Unknown (possible copy of directory strings count)
	EmptyValues2       [24]byte // Unknown empty values
	AlignmentPadding   uint32   // Unknown (possible alignment padding, may contain remnant data)
}
type fileRefs1 struct {
	Unknown1      uint32 // maybe version
	FileRefsCount uint32
	Unknown2      uint32
	Unknown3      uint32
}
type MftRef struct {
	MFTEntryIndex  [6]byte
	SequenceNumber uint16
}
