package logic

import "testing"

func TestNormalizeRepoSubPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
		valid bool
	}{
		{"team/app", "team/app", true},
		{" team\\app/ ", "team/app", true},
		{"", "", true},
		{"../other", "", false},
		{"a/../../other", "", false},
		{"/absolute", "", false},
		{"C:\\repo", "", false},
	}
	for _, tt := range tests {
		got, err := normalizeRepoSubPath(tt.input)
		if (err == nil) != tt.valid || got != tt.want {
			t.Errorf("normalizeRepoSubPath(%q) = %q, %v; want %q valid=%v", tt.input, got, err, tt.want, tt.valid)
		}
	}
}
