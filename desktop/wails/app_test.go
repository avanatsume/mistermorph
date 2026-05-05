//go:build wailsdesktop

package main

import "testing"

func TestNormalizeExternalBrowserURL(t *testing.T) {
	cases := []struct {
		name    string
		rawURL  string
		want    string
		wantErr bool
	}{
		{
			name:   "https",
			rawURL: " https://example.com/path?q=a~b(1)! ",
			want:   "https://example.com/path?q=a~b(1)!",
		},
		{
			name:   "http",
			rawURL: "http://example.com",
			want:   "http://example.com",
		},
		{
			name:    "missing host",
			rawURL:  "https:///path",
			wantErr: true,
		},
		{
			name:    "unsupported scheme",
			rawURL:  "file:///tmp/example",
			wantErr: true,
		},
		{
			name:    "control character",
			rawURL:  "https://example.com/\npath",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := normalizeExternalBrowserURL(tc.rawURL)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("normalizeExternalBrowserURL() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("normalizeExternalBrowserURL() error = %v", err)
			}
			if got != tc.want {
				t.Fatalf("normalizeExternalBrowserURL() = %q, want %q", got, tc.want)
			}
		})
	}
}
