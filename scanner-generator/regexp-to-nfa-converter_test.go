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
	}

	for _, tt := range tests {
		regexpToNfaConverter := &RegexpToNfaConverter{regexp: tt.input}
		_, err := regexpToNfaConverter.Convert()

		if err == nil {
			t.Fatalf("Expected an error to be '%s'. Got nil", tt.expected)
		}

		if err.Error() != tt.expected.Error() {
			t.Fatalf("Expected error to be '%s'. Got '%s'.", tt.expected, err)
		}
	}
}
