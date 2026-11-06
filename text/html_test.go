package text

import (
	"os"
	"path/filepath"
	"testing"
)

func TestToHTML(t *testing.T) {
	filePath := filepath.Join("data", "testdata.md")
	source, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", filePath, err)
	}

	filePath = filepath.Join("data", "testdata_out.txt")
	expectedOutput, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", filePath, err)
	}

	var parser Parser
	tokens, errs := parser.Parse(source)
	if len(errs) != 0 {
		t.Fatalf("got %d errors from source file", len(errs))
	}

	transformer := Transformer{
		source: source,
	}

	output, err := transformer.Article(tokens)
	if err != nil {
		t.Errorf("expected no errors but got: %s", err)
	}

	if output != string(expectedOutput) {
		t.Errorf("expected output in testdata_out, but got:\n%s", output)
	}
}
