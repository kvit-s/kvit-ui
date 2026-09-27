//go:build !(linux || darwin || freebsd || netbsd || openbsd || windows)

package text

import "os"

// mapFile reads a font file into memory where mapping is not available.
func mapFile(path string) ([]byte, error) { return os.ReadFile(path) }
