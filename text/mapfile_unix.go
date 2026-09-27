//go:build linux || darwin || freebsd || netbsd || openbsd

package text

import (
	"os"

	"golang.org/x/sys/unix"
)

// mapFile maps a font file read-only into memory. The pages belong to the
// file, so the operating system can drop and reload them, and they do not
// count as the process's private memory the way a copy read into the heap
// does. The mapping lasts for the life of the process, as the fonts do.
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
	if st.Size() == 0 {
		return nil, nil
	}
	data, err := unix.Mmap(int(f.Fd()), 0, int(st.Size()), unix.PROT_READ, unix.MAP_SHARED)
	if err != nil {
		return os.ReadFile(path)
	}
	return data, nil
}
