package main

import "testing"

func TestExpressionParsing(t *testing.T) {
	tests := []struct {
		name string
		fileContents string
		expected string
	}{
		{"Booleans/true", "true", "true"},
		{"Booleans/false", "false", "false"},
		{"Nil", "nil", "nil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner(tt.fileContents)

			err := scanner.Scan()
			if err != nil {
				t.Errorf("Scanner: tokenizing error: %v", err)
			}

			parser := NewParser(scanner.tokens)

			err = parser.Parse()
			if err != nil {
				t.Errorf("Parser: Error while parsing expressions")
			}

			actual := parser.StringifyExpressions()
			if actual != tt.expected {
				t.Errorf("Expression parsing result is incorrect\nExpected:\n\n%s\n\nGot:\n\n%s", tt.expected, actual)
			}
		})
	}
}