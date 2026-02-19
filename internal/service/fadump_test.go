package service

import (
	"testing"
)

func TestParseFontAwesomeStyles(t *testing.T) {
	tests := []struct {
		data     string
		expected string
	}{
		{
			`.fa-markdown {
  --fa: "\f60f"; }

.fa-sourcetree {
  --fa: "\f7d3"; }`,
			"markdown: \\uf60f\nsourcetree: \\uf7d3\n",
		},
	}
	for _, test := range tests {
		result, err := parseFontAwesomeStyles([]byte(test.data))
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func TestFormatUnicodeLiteral(t *testing.T) {
	tests := []struct {
		unicode  string
		expected string
	}{
		{"30", "\\u0030"},
		{"f3e1", "\\uf3e1"},
	}
	for _, test := range tests {
		result := formatUnicodeLiteral(test.unicode)
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}
