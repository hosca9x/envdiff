package sorter

import "testing"

func TestStrategyString(t *testing.T) {
	cases := []struct {
		strategy Strategy
		want     string
	}{
		{Alphabetical, "alphabetical"},
		{Reverse, "reverse"},
		{GroupedByPrefix, "grouped"},
		{Strategy(99), "strategy(99)"},
	}
	for _, tc := range cases {
		if got := tc.strategy.String(); got != tc.want {
			t.Errorf("Strategy(%d).String() = %q, want %q", int(tc.strategy), got, tc.want)
		}
	}
}

func TestParseStrategy_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  Strategy
	}{
		{"alphabetical", Alphabetical},
		{"alpha", Alphabetical},
		{"asc", Alphabetical},
		{"reverse", Reverse},
		{"desc", Reverse},
		{"grouped", GroupedByPrefix},
		{"prefix", GroupedByPrefix},
		{"GROUPED", GroupedByPrefix},
		{"  Reverse  ", Reverse},
	}
	for _, tc := range cases {
		got, err := ParseStrategy(tc.input)
		if err != nil {
			t.Errorf("ParseStrategy(%q) unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseStrategy(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParseStrategy_Invalid(t *testing.T) {
	_, err := ParseStrategy("random")
	if err == nil {
		t.Error("expected error for unknown strategy, got nil")
	}
}
