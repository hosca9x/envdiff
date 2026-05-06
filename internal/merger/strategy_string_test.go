package merger_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/merger"
)

func TestStrategyString(t *testing.T) {
	tests := []struct {
		strategy merger.Strategy
		want     string
	}{
		{merger.FirstWins, "FirstWins"},
		{merger.LastWins, "LastWins"},
		{merger.ErrorOnConflict, "ErrorOnConflict"},
		{merger.Strategy(99), "Strategy(99)"},
	}
	for _, tt := range tests {
		if got := tt.strategy.String(); got != tt.want {
			t.Errorf("Strategy(%d).String() = %q, want %q", tt.strategy, got, tt.want)
		}
	}
}

func TestParseStrategy_Valid(t *testing.T) {
	tests := []struct {
		input string
		want  merger.Strategy
	}{
		{"first", merger.FirstWins},
		{"FirstWins", merger.FirstWins},
		{"last", merger.LastWins},
		{"LastWins", merger.LastWins},
		{"error", merger.ErrorOnConflict},
		{"ErrorOnConflict", merger.ErrorOnConflict},
	}
	for _, tt := range tests {
		got, err := merger.ParseStrategy(tt.input)
		if err != nil {
			t.Errorf("ParseStrategy(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseStrategy(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseStrategy_Invalid(t *testing.T) {
	_, err := merger.ParseStrategy("unknown")
	if err == nil {
		t.Error("expected error for unknown strategy, got nil")
	}
}
