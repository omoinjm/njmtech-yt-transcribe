package main

import (
	"testing"
)

// TestSanitizeFilename tests the sanitizeFilename helper function.
func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"My Video Title.mp3", "My_Video_Title.mp3"},
		{"Video/Title:With?Special*Chars", "Video_Title_With_Special_Chars"},
		{"Another\\\\Video|Title<>", "Another__Video_Title__"},
		{"NoSpecialCharsHere", "NoSpecialCharsHere"},
		{"", ""},
		{"   leading and trailing spaces   ", "___leading_and_trailing_spaces___"},
		{"file with spaces", "file_with_spaces"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			act := sanitizeFilename(test.input)
			if act != test.expected {
				t.Errorf("For input '%s', expected '%s', got '%s'", test.input, test.expected, act)
			}
		})
	}
}

