package scannergenerator

import (
	"fmt"
	"testing"
)

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected error
	}{
		{
			"a(",
			fmt.Errorf("expected closing ')'"),
		},
		{
			"a)",
			fmt.Errorf("expected opening ')'"),
		},
		{
			"[2-1]",
			fmt.Errorf("lower bound greater or equal to upper bound '[2-1]'"),
		},
		{
			"a|",
			fmt.Errorf("expected right side of |"),
		},
	}

	for _, tt := range tests {
		regexpToNfaConverter := NewRegexpToNfaConverter(tt.input)
		_, err := regexpToNfaConverter.Convert()

		if err == nil {
			t.Fatalf("Expected an error to be '%s'. Got nil", tt.expected)
		}

		if err.Error() != tt.expected.Error() {
			t.Fatalf("Expected error to be '%s'. Got '%s'.", tt.expected, err)
		}
	}
}
