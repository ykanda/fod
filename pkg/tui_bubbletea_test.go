package fod

import (
	"testing"

	"github.com/mattn/go-runewidth"
)

func TestTruncateLine_DisplayWidth(t *testing.T) {
	text := "あああ"

	if got := truncateLine(text, 6); got != text {
		t.Fatalf("truncateLine(%q, 6) = %q, want %q", text, got, text)
	}

	if got := truncateLine(text, 5); got != "ああ" {
		t.Fatalf("truncateLine(%q, 5) = %q, want %q", text, got, "ああ")
	}

	if got := truncateLine(text, 5); runewidth.StringWidth(got) > 5 {
		t.Fatalf("truncateLine(%q, 5) width = %d, want <= 5", text, runewidth.StringWidth(got))
	}
}

func TestTruncateRunes_DisplayWidth(t *testing.T) {
	text := "あい"

	if got := truncateRunes(text, 4); got != text {
		t.Fatalf("truncateRunes(%q, 4) = %q, want %q", text, got, text)
	}

	if got := truncateRunes(text, 3); got != "あ" {
		t.Fatalf("truncateRunes(%q, 3) = %q, want %q", text, got, "あ")
	}

	if got := truncateRunes(text, 1); got != "" {
		t.Fatalf("truncateRunes(%q, 1) = %q, want %q", text, got, "")
	}
}
