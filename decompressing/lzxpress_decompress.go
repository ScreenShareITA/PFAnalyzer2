package decompressing

import (
	"fmt"

	"unsafe"

	"golang.org/x/sys/windows"
)

// Costanti per la decompressione
var (
	ntoskrnl                       = windows.NewLazySystemDLL("ntdll.dll")
	RtlDecompressBufferEx          = ntoskrnl.NewProc("RtlDecompressBufferEx")
	COMPRESSION_FORMAT_XPRESS_HUFF = uint16(0x0004)
)

// DecompressBuffer decomprime i dati utilizzando prima RtlDecompressBufferEx e, se fallisce,
// utilizza l'implementazione manuale
func DecompressBuffer(compressedBuffer []byte, uncompressedSize uint32) ([]byte, error) {
	result := make([]byte, uncompressedSize)
	finalSize := uncompressedSize
	workspace := make([]byte, uncompressedSize*2)
	ret, _, _ := RtlDecompressBufferEx.Call(uintptr(COMPRESSION_FORMAT_XPRESS_HUFF),
		uintptr(unsafe.Pointer(&result[0])),
		uintptr(uncompressedSize),
		uintptr(unsafe.Pointer(&compressedBuffer[0])),
		uintptr(len(compressedBuffer)),
		uintptr(unsafe.Pointer(&finalSize)),
		uintptr(unsafe.Pointer(&workspace[0])))

	if ret == 0 {
		return result[:finalSize], nil
	}
	// IF IT DOESNT WORK WITH WINDOWS API. MANUAL DECOMPRESSION
	return nil, fmt.Errorf("RtlDecompressBufferEx fallita")
}
