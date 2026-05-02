package organizer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lechip/switch-screenshots/internal/gameids"
)

// knownID maps to "1-2-Switch" in game_ids.json
const knownID = "2B1F1288BC05B2D89D8431910DBA2878"

// unknownID is deliberately not in game_ids.json
const unknownID = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"

func TestRun_BasicCopy(t *testing.T) {
	input := t.TempDir()
	output := t.TempDir()

	files := []string{
		"shot-" + knownID + ".jpg",
		"shot-" + unknownID + ".mp4",
		"notascreenshot.txt",
		filepath.Join("subdir", "nested-"+knownID+".jpg"),
	}
	if err := os.MkdirAll(filepath.Join(input, "subdir"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(input, f), []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	cfg := Config{Input: input, Output: output}
	res, err := Run(cfg)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	if res.Processed != 3 {
		t.Errorf("Processed = %d, want 3", res.Processed)
	}
	if res.Skipped != 1 {
		t.Errorf("Skipped = %d, want 1", res.Skipped)
	}
	if res.UnknownIDs[unknownID] != 1 {
		t.Errorf("UnknownIDs[%s] = %d, want 1", unknownID, res.UnknownIDs[unknownID])
	}

	// Known game files should land in the title-named folder
	wantDir := filepath.Join(output, "1-2-Switch")
	for _, name := range []string{"shot-" + knownID + ".jpg", "nested-" + knownID + ".jpg"} {
		if _, err := os.Stat(filepath.Join(wantDir, name)); err != nil {
			t.Errorf("expected output file %s: %v", filepath.Join(wantDir, name), err)
		}
	}

	// Unknown game file should fall back to the raw ID folder
	wantUnknown := filepath.Join(output, unknownID, "shot-"+unknownID+".mp4")
	if _, err := os.Stat(wantUnknown); err != nil {
		t.Errorf("expected fallback file %s: %v", wantUnknown, err)
	}
}

func TestRun_DryRun(t *testing.T) {
	input := t.TempDir()
	output := t.TempDir()

	name := "shot-" + knownID + ".jpg"
	if err := os.WriteFile(filepath.Join(input, name), []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := Config{Input: input, Output: output, DryRun: true}
	res, err := Run(cfg)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if res.Processed != 1 {
		t.Errorf("Processed = %d, want 1", res.Processed)
	}

	// Output directory should remain empty in dry-run mode
	entries, _ := os.ReadDir(output)
	if len(entries) != 0 {
		t.Errorf("dry-run wrote %d entries to output, want 0", len(entries))
	}
}

func TestRun_ByID(t *testing.T) {
	input := t.TempDir()
	output := t.TempDir()

	name := "shot-" + knownID + ".jpg"
	if err := os.WriteFile(filepath.Join(input, name), []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := Config{Input: input, Output: output, ByID: true}
	if _, err := Run(cfg); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	// --by-id: folder should be the raw game ID, not the title
	wantPath := filepath.Join(output, knownID, name)
	if _, err := os.Stat(wantPath); err != nil {
		t.Errorf("expected raw-ID folder file %s: %v", wantPath, err)
	}
}

func TestSanitizeTitle(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1-2-Switch", "1-2-Switch"},
		{"The Legend of Zelda: Breath of the Wild", "The Legend of Zelda Breath of the Wild"},
		{"A/B", "A B"},
		{"trailing dots...", "trailing dots"},
		{"CON", "CON_"},
		{"", "_"},
	}
	for _, tt := range tests {
		got := sanitizeTitle(tt.input)
		if got != tt.want {
			t.Errorf("sanitizeTitle(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSanitizeTitle_AllGameIDs(t *testing.T) {
	for _, title := range gameids.All() {
		got := sanitizeTitle(title)
		if got == "" {
			t.Errorf("sanitizeTitle(%q) = empty string", title)
		}
		if strings.Contains(got, "/") {
			t.Errorf("sanitizeTitle(%q) = %q contains '/'", title, got)
		}
		if strings.Contains(got, "\x00") {
			t.Errorf("sanitizeTitle(%q) = %q contains null byte", title, got)
		}
		if got == "." || got == ".." {
			t.Errorf("sanitizeTitle(%q) = %q is an invalid path component", title, got)
		}
	}
}

func TestRun_ShowSkipped(t *testing.T) {
	input := t.TempDir()
	output := t.TempDir()

	screenshot := "shot-" + knownID + ".jpg"
	ignored := "not-a-screenshot.txt"
	for _, f := range []string{screenshot, ignored} {
		if err := os.WriteFile(filepath.Join(input, f), []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("show-skipped enabled", func(t *testing.T) {
		res, err := Run(Config{Input: input, Output: output, ShowSkipped: true})
		if err != nil {
			t.Fatalf("Run() error: %v", err)
		}
		if len(res.SkippedFiles) != 1 {
			t.Fatalf("SkippedFiles = %v, want 1 entry", res.SkippedFiles)
		}
		if !strings.HasSuffix(res.SkippedFiles[0], ignored) {
			t.Errorf("SkippedFiles[0] = %q, want suffix %q", res.SkippedFiles[0], ignored)
		}
	})

	t.Run("show-skipped disabled", func(t *testing.T) {
		res, err := Run(Config{Input: input, Output: output, ShowSkipped: false})
		if err != nil {
			t.Fatalf("Run() error: %v", err)
		}
		if len(res.SkippedFiles) != 0 {
			t.Errorf("SkippedFiles = %v, want empty", res.SkippedFiles)
		}
	})
}
