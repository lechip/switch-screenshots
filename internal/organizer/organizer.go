package organizer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lechip/switch-screenshots/internal/fsutil"
	"github.com/lechip/switch-screenshots/internal/gameids"
)

// fileRe matches the 32-char hex game ID and extension in a Switch filename.
// Example: 2018101310443100-EBDA73B0F0F1E8C57869EB26D8F65095.jpg
var fileRe = regexp.MustCompile(`(?i)-([0-9A-Fa-f]{32})\.(jpg|mp4)$`)

// Config controls a Run invocation.
type Config struct {
	Input       string
	Output      string
	DryRun      bool
	ByID        bool
	Verbose     bool
	ShowSkipped bool
}

// Result summarizes a completed run.
type Result struct {
	Processed    int
	Skipped      int
	SkippedFiles []string       // populated when Config.ShowSkipped is true
	UnknownIDs   map[string]int // game ID -> number of files with that ID
}

// Run walks cfg.Input, copies matching media files into per-game subdirectories
// under cfg.Output, and returns a summary.
func Run(cfg Config) (Result, error) {
	res := Result{UnknownIDs: make(map[string]int)}

	if !cfg.DryRun {
		if err := os.MkdirAll(cfg.Output, 0o755); err != nil {
			return res, fmt.Errorf("create output dir: %w", err)
		}
	}

	err := filepath.WalkDir(cfg.Input, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}

		m := fileRe.FindStringSubmatch(d.Name())
		if m == nil {
			res.Skipped++
			if cfg.ShowSkipped {
				res.SkippedFiles = append(res.SkippedFiles, path)
			}
			return nil
		}

		id := strings.ToUpper(m[1])
		folder := id
		if !cfg.ByID {
			if title, ok := gameids.Lookup(id); ok {
				folder = sanitizeTitle(title)
			} else {
				res.UnknownIDs[id]++
			}
		}

		dst := filepath.Join(cfg.Output, folder, d.Name())
		if cfg.Verbose || cfg.DryRun {
			fmt.Printf("%s → %s\n", path, dst)
		}
		if !cfg.DryRun {
			if err := fsutil.Copy(path, dst); err != nil {
				return fmt.Errorf("copy %s: %w", path, err)
			}
		}
		res.Processed++
		return nil
	})
	return res, err
}

var illegalCharsReplacer = strings.NewReplacer(
	":", " ",
	"/", " ",
	"<", " ",
	">", " ",
	`"`, " ",
	`\`, " ",
	"|", " ",
	"?", " ",
	"*", " ",
)

var windowsReserved = map[string]bool{
	"CON": true, "PRN": true, "AUX": true, "NUL": true,
	"COM1": true, "COM2": true, "COM3": true, "COM4": true,
	"COM5": true, "COM6": true, "COM7": true, "COM8": true, "COM9": true,
	"LPT1": true, "LPT2": true, "LPT3": true, "LPT4": true,
	"LPT5": true, "LPT6": true, "LPT7": true, "LPT8": true, "LPT9": true,
}

// sanitizeTitle converts a game title to a safe directory name on all
// supported filesystems (APFS, ext4/btrfs, exFAT, NTFS).
func sanitizeTitle(title string) string {
	s := illegalCharsReplacer.Replace(title)
	s = strings.Join(strings.Fields(s), " ") // collapse consecutive spaces
	s = strings.TrimRight(s, ". ")
	if s == "" {
		return "_"
	}
	if windowsReserved[strings.ToUpper(s)] {
		s += "_"
	}
	return s
}
