package main

import "testing"

func TestNormalizeNextPath(t *testing.T) {
	tests := map[string]string{
		"":                       "/admin",
		"/admin":                 "/admin",
		"/sso/authorize?state=x": "/sso/authorize?state=x",
		"//evil.example":         "/admin",
		"https://evil.example":   "/admin",
		"/admin\r\nLocation: /":  "/admin",
	}

	for input, want := range tests {
		if got := normalizeNextPath(input); got != want {
			t.Errorf("normalizeNextPath(%q) = %q, want %q", input, got, want)
		}
	}
}
