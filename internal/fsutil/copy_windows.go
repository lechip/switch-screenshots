//go:build windows

package fsutil

// copyFile uses a buffered copy on Windows.
// TODO: replace with windows.CopyFileW from golang.org/x/sys/windows for
// better performance — see migration-to-go.md §4.2.
func copyFile(src, dst string) error {
	return copyBuffered(src, dst)
}
