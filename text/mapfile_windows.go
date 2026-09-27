package text

import (
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

// mapFile maps a font file read-only into memory. The pages belong to the
// file, so Windows can drop and reload them, and they count as shareable
// rather than private memory. The mapping lasts for the life of the process,
// as the fonts do.
func mapFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := st.Size()
	if size == 0 {
		return nil, nil
	}
	h, err := windows.CreateFileMapping(windows.Handle(f.Fd()), nil, windows.PAGE_READONLY, uint32(size>>32), uint32(size), nil)
	if err != nil {
		return os.ReadFile(path)
	}
	addr, err := windows.MapViewOfFile(h, windows.FILE_MAP_READ, 0, 0, uintptr(size))
	// The view keeps the mapping alive; the handle is no longer needed.
	windows.CloseHandle(h)
	if err != nil {
		return os.ReadFile(path)
	}
	// Read the address through a pointer rather than converting the uintptr, which
	// go vet cannot tell from an unsafe conversion.
	return unsafe.Slice((*byte)(*(*unsafe.Pointer)(unsafe.Pointer(&addr))), size), nil
}
