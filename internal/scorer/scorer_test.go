package scorer

import (
	"testing"
)

func TestEvaluate_PerfectScore(t *testing.T) {
	env := map[string]string{
		"APP_NAME":    "myapp",
		"API_SECRET":  "supersecretvalue123",
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "strongpassword99",
	}
	s := Evaluate(env)
	if s.Total != 100 {
		t.Errorf("expected 100, got %d (penalties: %v)", s.Total, s.Penalties)
	}
	if len(s.Penalties) != 0 {
		t.Errorf("expected no penalties, got %v", s.Penalties)
	}
}

func TestEvaluate_LowercaseKeyPenalty(t *testing.T) {
	env := map[string]string{
		"app_name": "myapp",
		"DB_HOST":  "localhost",
	}
	s := Evaluate(env)
	if s.Breakdown["key_format"] >= 25 {
		t.Errorf("expected key_format penalty, got %d", s.Breakdown["key_format"])
	}
	if s.Total >= 100 {
		t.Errorf("expected total < 100, got %d", s.Total)
	}
}

func TestEvaluate_EmptyValuePenalty(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "",
		"DB_HOST":  "localhost",
	}
	s := Evaluate(env)
	if s.Breakdown["value_complete"] >= 25 {
		t.Errorf("expected value_complete penalty, got %d", s.Breakdown["value_complete"])
	}
	found := false
	for _, p := range s.Penalties {
		if p == "empty value for key: APP_NAME" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected penalty for APP_NAME empty value, got %v", s.Penalties)
	}
}

func TestEvaluate_WeakSecretPenalty(t *testing.T) {
	env := map[string]string{
		"API_SECRET": "abc",
		"APP_NAME":   "myapp",
	}
	s := Evaluate(env)
	if s.Breakdown["secret_hygiene"] >= 25 {
		t.Errorf("expected secret_hygiene penalty, got %d", s.Breakdown["secret_hygiene"])
	}
}

func TestEvaluate_EmptyMap(t *testing.T) {
	s := Evaluate(map[string]string{})
	if s.Total != 0 {
		t.Errorf("expected 0 for empty map, got %d", s.Total)
	}
	if len(s.Penalties) == 0 {
		t.Error("expected at least one penalty for empty map")
	}
}

func TestEvaluate_BreakdownSumsToTotal(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"DB_HOST":  "",
		"api_key":  "short",
	}
	s := Evaluate(env)
	sum := 0
	for _, v := range s.Breakdown {
		sum += v
	}
	if sum != s.Total {
		t.Errorf("breakdown sum %d does not match total %d", sum, s.Total)
	}
}

func TestIsSensitive_Patterns(t *testing.T) {
	cases := []struct {
		key      string
		want     bool
	}{
		{"DB_PASSWORD", true},
		{"API_TOKEN", true},
		{"APP_NAME", false},
		{"PRIVATE_KEY", true},
		{"HOST", false},
	}
	for _, tc := range cases {
		got := isSensitive(tc.key)
		if got != tc.want {
			t.Errorf("isSensitive(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}
