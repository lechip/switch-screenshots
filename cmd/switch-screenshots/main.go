package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/lechip/switch-screenshots/internal/organizer"
)

// Injected by GoReleaser via -ldflags "-X main.version=..."
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var (
		input   string
		output  string
		dryRun  bool
		byID    bool
		verbose bool
		ver     bool
	)

	flag.StringVar(&input, "input", "", "path to the Switch SD card Album folder (required)")
	flag.StringVar(&input, "i", "", "shorthand for --input")
	flag.StringVar(&output, "output", "./screenshotsOutput", "destination folder")
	flag.StringVar(&output, "o", "./screenshotsOutput", "shorthand for --output")
	flag.BoolVar(&dryRun, "dry-run", false, "print planned operations without copying any files")
	flag.BoolVar(&byID, "by-id", false, "use the raw 32-char game ID for folder names instead of the title")
	flag.BoolVar(&verbose, "v", false, "log each file operation")
	flag.BoolVar(&ver, "version", false, "print version information and exit")
	flag.Parse()

	if ver {
		fmt.Printf("switch-screenshots %s (commit %s, built %s)\n", version, commit, date)
		return
	}

	if input == "" {
		fmt.Fprintln(os.Stderr, "error: --input / -i is required")
		flag.Usage()
		os.Exit(1)
	}

	cfg := organizer.Config{
		Input:   input,
		Output:  output,
		DryRun:  dryRun,
		ByID:    byID,
		Verbose: verbose,
	}

	fmt.Printf("Organizing screenshots from %s → %s\n", cfg.Input, cfg.Output)
	if cfg.DryRun {
		fmt.Println("(dry run — no files will be copied)")
	}

	result, err := organizer.Run(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf("Processed %d file(s)", result.Processed)
	if result.Skipped > 0 {
		fmt.Printf(", skipped %d unrecognized file(s)", result.Skipped)
	}
	fmt.Println()

	if len(result.UnknownIDs) > 0 {
		fmt.Println("\nUnknown game IDs (consider running: make update-gameids):")
		ids := make([]string, 0, len(result.UnknownIDs))
		for id := range result.UnknownIDs {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		for _, id := range ids {
			fmt.Printf("  %s  (%d file(s))\n", id, result.UnknownIDs[id])
		}
	}
}
