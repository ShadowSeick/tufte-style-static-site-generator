package main

import (
	"testing"
)

func TestHeader(t *testing.T) {
	text := "## My Header {@header-id}"
	matches := Header.Match(text)
	
	if matches == nil {
		t.Fatal("Header regex did not match")
	}
	
	text1, text2 := Header.Text(matches)
	
	expectedText := "My Header"
	expectedId := "header-id"
	
	if text1 != expectedText {
		t.Errorf("Expected header text '%s', got '%s'", expectedText, text1)
	}
	
	if text2 != expectedId {
		t.Errorf("Expected header id '%s', got '%s'", expectedId, text2)
	}
}

func TestSubheader(t *testing.T) {
	text := "Some [^sub-header](This is a subheader) text"
	matches := Subheader.Match(text)
	
	if matches == nil {
		t.Fatal("Subheader regex did not match")
	}
	
	text1, text2 := Subheader.Text(matches)
	
	expectedText := "This is a subheader"
	
	if text1 != expectedText {
		t.Errorf("Expected subheader text '%s', got '%s'", expectedText, text1)
	}
	
	if text2 != "" {
		t.Errorf("Expected empty string for text2, got '%s'", text2)
	}
}

func TestLink(t *testing.T) {
	text := "Check [this link](https://example.com) out"
	matches := Link.Match(text)
	
	if matches == nil {
		t.Fatal("Link regex did not match")
	}
	
	text1, text2 := Link.Text(matches)
	
	expectedText := "this link"
	expectedUrl := "https://example.com"
	
	if text1 != expectedText {
		t.Errorf("Expected link text '%s', got '%s'", expectedText, text1)
	}
	
	if text2 != expectedUrl {
		t.Errorf("Expected link URL '%s', got '%s'", expectedUrl, text2)
	}
}

func TestSideNote(t *testing.T) {
	text := "Main text [^side-note](This is a side note) continues"
	matches := SideNote.Match(text)
	
	if matches == nil {
		t.Fatal("SideNote regex did not match")
	}
	
	text1, text2 := SideNote.Text(matches)
	
	expectedText := "This is a side note"
	
	if text1 != expectedText {
		t.Errorf("Expected side note text '%s', got '%s'", expectedText, text1)
	}
	
	if text2 != "" {
		t.Errorf("Expected empty string for text2, got '%s'", text2)
	}
}

func TestMarginNote(t *testing.T) {
	text := "Main text [^margin-note](This is a margin note) continues"
	matches := MarginNote.Match(text)
	
	if matches == nil {
		t.Fatal("MarginNote regex did not match")
	}
	
	text1, text2 := MarginNote.Text(matches)
	
	expectedText := "This is a margin note"
	
	if text1 != expectedText {
		t.Errorf("Expected margin note text '%s', got '%s'", expectedText, text1)
	}
	
	if text2 != "" {
		t.Errorf("Expected empty string for text2, got '%s'", text2)
	}
}

func TestCode(t *testing.T) {
	text := "```go\nfunc main() {\n  text := 1\n  return text\n}\n```"
	matches := Code.Match(text)
	
	if matches == nil {
		t.Fatal("Code regex did not match")
	}
	
	language, code := Code.Text(matches)
	
	// Note: Based on your regex "```(.*)```", it captures everything between backticks
	// The actual capture will depend on your regex implementation
	expectedLanguage := "go"
	expectedCode := "func main() {\n  text := 1\n  return text\n}\n"
	
	if language != expectedLanguage {
		t.Errorf("Expected language '%s', got '%s'", expectedLanguage, language)
	}
	
	if code != expectedCode {
		t.Errorf("Expected code content: '%s', got '%s'", expectedCode, code)
	}
}

func TestInlineCode(t *testing.T) {
	text := "Use `fmt.Println()` to print"
	matches := InlineCode.Match(text)
	
	if matches == nil {
		t.Fatal("InlineCode regex did not match")
	}
	
	text1, text2 := InlineCode.Text(matches)
	
	expectedCode := "fmt.Println()"
	
	if text1 != expectedCode {
		t.Errorf("Expected inline code '%s', got '%s'", expectedCode, text1)
	}
	
	if text2 != "" {
		t.Errorf("Expected empty string for text2, got '%s'", text2)
	}
}

func TestImage(t *testing.T) {
	text := "![alt text](image.png)"
	matches := Image.Match(text)
	
	if matches == nil {
		t.Fatal("Image regex did not match")
	}
	
	text1, text2 := Image.Text(matches)
	
	expectedAlt := "alt text"
	expectedUrl := "image.png"
	
	if text1 != expectedAlt {
		t.Errorf("Expected image alt text '%s', got '%s'", expectedAlt, text1)
	}
	
	if text2 != expectedUrl {
		t.Errorf("Expected image URL '%s', got '%s'", expectedUrl, text2)
	}
}

func TestItalic(t *testing.T) {
	text := "This is **italic text** here"
	matches := Italic.Match(text)
	
	if matches == nil {
		t.Fatal("Italic regex did not match")
	}
	
	text1, text2 := Italic.Text(matches)
	
	expectedText := "italic text"
	
	if text1 != expectedText {
		t.Errorf("Expected italic text '%s', got '%s'", expectedText, text1)
	}
	
	if text2 != "" {
		t.Errorf("Expected empty string for text2, got '%s'", text2)
	}
}

func TestBold(t *testing.T) {
	text := "This is *bold text* here"
	matches := Bold.Match(text)
	
	if matches == nil {
		t.Fatal("Bold regex did not match")
	}
	
	text1, text2 := Bold.Text(matches)
	
	expectedText := "bold text"
	
	if text1 != expectedText {
		t.Errorf("Expected bold text '%s', got '%s'", expectedText, text1)
	}
	
	if text2 != "" {
		t.Errorf("Expected empty string for text2, got '%s'", text2)
	}
}
