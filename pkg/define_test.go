package fod

import "testing"

func TestModeString(t *testing.T) {
	cases := []struct {
		mode Mode
		want string
	}{
		{ModeInvalid, "invalid"},
		{ModeFile, "file"},
		{ModeDirectory, "directory"},
		{Mode(99), "invalid"},
	}
	for _, tc := range cases {
		if got := tc.mode.String(); got != tc.want {
			t.Fatalf("Mode(%d).String() = %q, want %q", tc.mode, got, tc.want)
		}
	}
}

func TestStringToMode(t *testing.T) {
	cases := []struct {
		in   string
		want Mode
		ok   bool
	}{
		{"file", ModeFile, true},
		{"f", ModeFile, true},
		{"directory", ModeDirectory, true},
		{"dir", ModeDirectory, true},
		{"d", ModeDirectory, true},
		{"", ModeInvalid, false},
		{"unknown", ModeInvalid, false},
	}
	for _, tc := range cases {
		got, err := StringToMode(tc.in)
		if tc.ok {
			if err != nil {
				t.Fatalf("StringToMode(%q) unexpected error: %v", tc.in, err)
			}
			if got != tc.want {
				t.Fatalf("StringToMode(%q) = %v, want %v", tc.in, got, tc.want)
			}
			continue
		}
		if err == nil {
			t.Fatalf("StringToMode(%q) expected error", tc.in)
		}
		if got != ModeInvalid {
			t.Fatalf("StringToMode(%q) = %v, want %v", tc.in, got, ModeInvalid)
		}
	}
}
