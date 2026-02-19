package display

import (
	"sway-icon-to-go/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameFormatter_Format(t *testing.T) {
	testCases := []struct {
		name            string
		format          *config.Format
		workspaceNumber int64
		appIcons        []string
		expected        string
	}{
		{
			name:            "default",
			format:          config.DefaultFormat(),
			workspaceNumber: 1,
			appIcons:        []string{"app1", "app2", "app3"},
			expected:        "1: app1|app2|app3",
		},
		{
			name:            "uniq",
			format:          &config.Format{Length: 10, Delimiter: "|", Uniq: true},
			workspaceNumber: 100,
			appIcons:        []string{"app1", "app1", "app2"},
			expected:        "100: app1|app2",
		},
		{
			name:            "length",
			format:          &config.Format{Length: 10, Delimiter: "+", Uniq: true},
			workspaceNumber: 1000,
			appIcons:        []string{"1234567890123", "app2", "app3"},
			expected:        "1000: 1234567890+app2+app3",
		},
		{
			name:            "delimiter",
			format:          &config.Format{Length: 10, Delimiter: " ", Uniq: true},
			workspaceNumber: 10000,
			appIcons:        []string{"1234567890123", "app2", "app3"},
			expected:        "10000: 1234567890 app2 app3",
		},
	}
	for _, testCase := range testCases {
		formatter := NewNameFormatter(testCase.format)
		formatted := formatter.Format(testCase.workspaceNumber, testCase.appIcons)
		assert.Equal(t, testCase.expected, formatted)
	}
}
