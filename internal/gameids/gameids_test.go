package gameids

import "testing"

func TestLookup(t *testing.T) {
	tests := []struct {
		id     string
		want   string
		wantOK bool
	}{
		// Known game — "1-2-Switch" from game_ids.json
		{"2B1F1288BC05B2D89D8431910DBA2878", "1-2-Switch", true},
		// Same ID in lowercase — lookup must be case-insensitive
		{"2b1f1288bc05b2d89d8431910dba2878", "1-2-Switch", true},
		// Unknown ID
		{"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", "", false},
	}
	for _, tt := range tests {
		got, ok := Lookup(tt.id)
		if ok != tt.wantOK {
			t.Errorf("Lookup(%q) ok=%v, want %v", tt.id, ok, tt.wantOK)
			continue
		}
		if ok && got != tt.want {
			t.Errorf("Lookup(%q) = %q, want %q", tt.id, got, tt.want)
		}
	}
}

func TestAllEntriesParsed(t *testing.T) {
	if len(ids) == 0 {
		t.Fatal("game_ids.json parsed to empty map")
	}
}
