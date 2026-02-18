package domain

import (
	"testing"
)

func TestHeader(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldMatch bool
		expected    string
	}{
		{
			name:        "H2 header",
			input:       "## My Header {@header-id}",
			shouldMatch: true,
			expected:    `<h2 id="header-id">My Header</h2>`,
		},
		{
			name:        "H1 header",
			input:       "# Main Title {@main}",
			shouldMatch: true,
			expected:    `<h1 id="main">Main Title</h1>`,
		},
		{
			name:        "H3 header",
			input:       "### Subsection {@subsection-id}",
			shouldMatch: true,
			expected:    `<h3 id="subsection-id">Subsection</h3>`,
		},
		{
			name:        "Not a header",
			input:       "Just regular text",
			shouldMatch: false,
			expected:    "Just regular text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, matched := Header.Html(tt.input)
			
			if matched != tt.shouldMatch {
				t.Errorf("Expected match=%v, got match=%v", tt.shouldMatch, matched)
			}
			
			if tt.shouldMatch && result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestSubheader(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldMatch bool
		expected    string
	}{
		{
			name:        "Subheader in text",
			input:       "Some [^sub-header](This is a subheader) text",
			shouldMatch: true,
			expected:    `Some <p class="subtitle">This is a subheader</p> text`,
		},
		{
			name:        "Plain subheader",
			input:       "[^sub-header](My Subtitle)",
			shouldMatch: true,
			expected:    `<p class="subtitle">My Subtitle</p>`,
		},
		{
			name:        "No subheader",
			input:       "Just regular text",
			shouldMatch: false,
			expected:    "Just regular text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, matched := Subheader.Html(tt.input)
			
			if matched != tt.shouldMatch {
				t.Errorf("Expected match=%v, got match=%v", tt.shouldMatch, matched)
			}
			
			if tt.shouldMatch && result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestLink(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldMatch bool
		expected    string
	}{
		{
			name:        "Link in text",
			input:       "Check [this link](https://example.com) out",
			shouldMatch: true,
			expected:    `Check <a href="https://example.com">this link</a> out`,
		},
		{
			name:        "Multiple links",
			input:       "[first](url1.com) and [second](url2.com)",
			shouldMatch: true,
			expected:    `<a href="url1.com">first</a> and <a href="url2.com">second</a>`,
		},
		{
			name:        "No link",
			input:       "Just text",
			shouldMatch: false,
			expected:    "Just text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, matched := Link.Html(tt.input)
			
			if matched != tt.shouldMatch {
				t.Errorf("Expected match=%v, got match=%v", tt.shouldMatch, matched)
			}
			
			if tt.shouldMatch && result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestSideNote(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldMatch bool
		expected    string
	}{
		{
			name:        "Side note with ID",
			input:       "Main text [^side-note @note1](This is a side note) continues",
			shouldMatch: true,
			expected:    `Main text <label for="sn-note1" class="margin-toggle sidenote-number"></label><input type="checkbox" id="sn-note1" class="margin-toggle"><span class="sidenote">This is a side note</span> continues`,
		},
		{
			name:        "Side note with already processed link",
			input:       `Text [^side-note @note2](See <a href="url.com">this</a> for more) text`,
			shouldMatch: true,
			expected:    `Text <label for="sn-note2" class="margin-toggle sidenote-number"></label><input type="checkbox" id="sn-note2" class="margin-toggle"><span class="sidenote">See <a href="url.com">this</a> for more</span> text`,
		},
		{
			name:        "Side note with already processed bold",
			input:       "[^side-note @note3](<b>Important</b> note)",
			shouldMatch: true,
			expected:    `<label for="sn-note3" class="margin-toggle sidenote-number"></label><input type="checkbox" id="sn-note3" class="margin-toggle"><span class="sidenote"><b>Important</b> note</span>`,
		},
		{
			name:        "Side note with already processed italic",
			input:       "[^side-note @note4](<em>Emphasis</em> here)",
			shouldMatch: true,
			expected:    `<label for="sn-note4" class="margin-toggle sidenote-number"></label><input type="checkbox" id="sn-note4" class="margin-toggle"><span class="sidenote"><em>Emphasis</em> here</span>`,
		},
		{
			name:        "Side note with already processed inline code",
			input:       "[^side-note @note5](Use <code>code</code> here)",
			shouldMatch: true,
			expected:    `<label for="sn-note5" class="margin-toggle sidenote-number"></label><input type="checkbox" id="sn-note5" class="margin-toggle"><span class="sidenote">Use <code>code</code> here</span>`,
		},
		{
			name:        "Side note with multiple processed elements",
			input:       "[^side-note @note6](<b>Bold</b> and <em>italic</em> with <code>code</code>)",
			shouldMatch: true,
			expected:    `<label for="sn-note6" class="margin-toggle sidenote-number"></label><input type="checkbox" id="sn-note6" class="margin-toggle"><span class="sidenote"><b>Bold</b> and <em>italic</em> with <code>code</code></span>`,
		},
		{
			name:        "No side note",
			input:       "Just text",
			shouldMatch: false,
			expected:    "Just text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, matched := SideNote.Html(tt.input)
			
			if matched != tt.shouldMatch {
				t.Errorf("Expected match=%v, got match=%v", tt.shouldMatch, matched)
			}
			
			if tt.shouldMatch && result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestMarginNote(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldMatch bool
		expected    string
	}{
		{
			name:        "Margin note with ID",
			input:       "Main text [^margin-note @mn1](This is a margin note) continues",
			shouldMatch: true,
			expected:    `Main text <label for="mn-mn1" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-mn1" class="margin-toggle"><span class="marginnote">This is a margin note</span> continues`,
		},
		{
			name:        "Margin note with already processed link",
			input:       `[^margin-note @mn2](Check <a href="url.com">link</a>)`,
			shouldMatch: true,
			expected:    `<label for="mn-mn2" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-mn2" class="margin-toggle"><span class="marginnote">Check <a href="url.com">link</a></span>`,
		},
		{
			name:        "Margin note with already processed bold",
			input:       "[^margin-note @mn3](<b>Bold text</b> note)",
			shouldMatch: true,
			expected:    `<label for="mn-mn3" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-mn3" class="margin-toggle"><span class="marginnote"><b>Bold text</b> note</span>`,
		},
		{
			name:        "Margin note with already processed italic",
			input:       "[^margin-note @mn4](<em>Italic text</em> note)",
			shouldMatch: true,
			expected:    `<label for="mn-mn4" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-mn4" class="margin-toggle"><span class="marginnote"><em>Italic text</em> note</span>`,
		},
		{
			name:        "Margin note with already processed inline code",
			input:       "[^margin-note @mn5](Use <code>function()</code> here)",
			shouldMatch: true,
			expected:    `<label for="mn-mn5" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-mn5" class="margin-toggle"><span class="marginnote">Use <code>function()</code> here</span>`,
		},
		{
			name:        "Margin note with multiple processed elements",
			input:       `[^margin-note @mn6](<b>See</b> <a href="doc.html">docs</a> with <code>code</code>)`,
			shouldMatch: true,
			expected:    `<label for="mn-mn6" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-mn6" class="margin-toggle"><span class="marginnote"><b>See</b> <a href="doc.html">docs</a> with <code>code</code></span>`,
		},
		{
			name:        "No margin note",
			input:       "Just text",
			shouldMatch: false,
			expected:    "Just text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, matched := MarginNote.Html(tt.input)
			
			if matched != tt.shouldMatch {
				t.Errorf("Expected match=%v, got match=%v", tt.shouldMatch, matched)
			}
			
			if tt.shouldMatch && result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestCode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldMatch bool
		expected    string
	}{
		{
			name:        "Go code block",
			input:       "```go\nfunc main() {\n  text := 1\n  return text\n}\n```",
			shouldMatch: true,
			expected:    `<pre><code class="language-go">func main() {
  text := 1
  return text
}
</code></pre>`,
		},
		{
			name:        "Python code block",
			input:       "```python\ndef hello():\n    print('world')\n```",
			shouldMatch: true,
			expected:    `<pre><code class="language-python">def hello():
    print('world')
</code></pre>`,
		},
		{
			name:        "Code block no language",
			input:       "```\nsome code\n```",
			shouldMatch: true,
			expected:    `<pre><code class="language-">some code
</code></pre>`,
		},
		{
			name:        "No code block",
			input:       "Just text",
			shouldMatch: false,
			expected:    "Just text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, matched := Code.Html(tt.input)
			
			if matched != tt.shouldMatch {
				t.Errorf("Expected match=%v, got match=%v", tt.shouldMatch, matched)
			}
			
			if tt.shouldMatch && result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestInlineCode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldMatch bool
		expected    string
	}{
		{
			name:        "Inline code",
			input:       "Here we have some thing\nUse `fmt.Println()` to print",
			shouldMatch: true,
			expected:    `Here we have some thing
Use <code>fmt.Println()</code> to print`,
		},
		{
			name: "Inline real code",
			input: "Maybe `we are dealing with existing code` and we have to work around it or we need to get out that feature so we cannot focus on what’s really important, data. Most of the time, if we have a good data representation, problems and features become easier to solve and the application becomes simpler and easier to maintain. Technical debt can be done by hurry that feature or not thinking twice about the solution we think is the correct one. It also can be develop by changing the needs of the business. Solutions are a matter of time and space, what was once good could no longer be. Refactors exist for a reason and we should do them from time to time, but sometimes if we spend a little more time thinking on the big picture we could make not so expensive refactors nor making features take longer than should be.",
			shouldMatch: true,
			expected: "Maybe <code>we are dealing with existing code</code> and we have to work around it or we need to get out that feature so we cannot focus on what’s really important, data. Most of the time, if we have a good data representation, problems and features become easier to solve and the application becomes simpler and easier to maintain. Technical debt can be done by hurry that feature or not thinking twice about the solution we think is the correct one. It also can be develop by changing the needs of the business. Solutions are a matter of time and space, what was once good could no longer be. Refactors exist for a reason and we should do them from time to time, but sometimes if we spend a little more time thinking on the big picture we could make not so expensive refactors nor making features take longer than should be.",
		},
		{
			name:        "Multiple inline code",
			input:       "`var1` and `var2`",
			shouldMatch: true,
			expected:    `<code>var1</code> and <code>var2</code>`,
		},
		{
			name:        "No inline code",
			input:       "Just text",
			shouldMatch: false,
			expected:    "Just text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, matched := InlineCode.Html(tt.input)
			
			if matched != tt.shouldMatch {
				t.Errorf("Expected match=%v, got match=%v", tt.shouldMatch, matched)
			}
			
			if tt.shouldMatch && result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestImage(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldMatch bool
		expected    string
	}{
		{
			name:        "Image with alt text",
			input:       "![alt text](image.png)",
			shouldMatch: true,
			expected:    `<figure><img src="image.png" alt="alt text"/></figure>`,
		},
		{
			name:        "Image without alt text",
			input:       "![](photo.jpg)",
			shouldMatch: true,
			expected:    `<figure><img src="photo.jpg" alt=""/></figure>`,
		},
		{
			name:        "No image",
			input:       "Just text",
			shouldMatch: false,
			expected:    "Just text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, matched := Image.Html(tt.input)
			
			if matched != tt.shouldMatch {
				t.Errorf("Expected match=%v, got match=%v", tt.shouldMatch, matched)
			}
			
			if tt.shouldMatch && result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestImageWithMarginNote tests images combined with margin notes
// This happens after both Image and MarginNote processing
func TestImageWithMarginNote(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Image followed by margin note",
			input:    `<label for="mn-img1" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-img1" class="margin-toggle"><span class="marginnote">See the trend here</span>![diagram](chart.png)`,
			expected: `<figure><label for="mn-img1" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-img1" class="margin-toggle"><span class="marginnote">See the trend here</span><img src="chart.png" alt="diagram"/></figure>`,
		},
		{
			name:     "Image with margin note containing processed elements",
			input:    `<label for="mn-img2" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-img2" class="margin-toggle"><span class="marginnote"><b>Key point:</b> Notice the <em>spike</em> at <code>t=5</code></span>![alt text](assets/image.png)`,
			expected: `<figure><label for="mn-img2" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-img2" class="margin-toggle"><span class="marginnote"><b>Key point:</b> Notice the <em>spike</em> at <code>t=5</code></span><img src="assets/image.png" alt="alt text"/></figure>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First process margin notes (since images are already processed in these test inputs)
			result, _ := Image.Html(tt.input)
			
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestItalic(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldMatch bool
		expected    string
	}{
		{
			name:        "Italic text",
			input:       "This is **italic text** here",
			shouldMatch: true,
			expected:    `This is <em>italic text</em> here`,
		},
		{
			name:        "Multiple italic",
			input:       "**first** and **second**",
			shouldMatch: true,
			expected:    `<em>first</em> and <em>second</em>`,
		},
		{
			name:        "No italic",
			input:       "Just text",
			shouldMatch: false,
			expected:    "Just text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, matched := Italic.Html(tt.input)
			
			if matched != tt.shouldMatch {
				t.Errorf("Expected match=%v, got match=%v", tt.shouldMatch, matched)
			}
			
			if tt.shouldMatch && result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestBold(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldMatch bool
		expected    string
	}{
		{
			name:        "Bold text",
			input:       "This is *bold text* here",
			shouldMatch: true,
			expected:    `This is <b>bold text</b> here`,
		},
		{
			name:        "Multiple bold",
			input:       "*first* and *second*",
			shouldMatch: true,
			expected:    `<b>first</b> and <b>second</b>`,
		},
		{
			name:        "No bold",
			input:       "Just text",
			shouldMatch: false,
			expected:    "Just text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, matched := Bold.Html(tt.input)
			
			if matched != tt.shouldMatch {
				t.Errorf("Expected match=%v, got match=%v", tt.shouldMatch, matched)
			}
			
			if tt.shouldMatch && result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestParagraph(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Paragraph with already processed bold",
			input:    "This is <b>bold text</b> in a paragraph",
			expected: `<p>This is <b>bold text</b> in a paragraph</p>`,
		},
		{
			name:     "Paragraph with already processed link",
			input:    `Check <a href="https://example.com">this link</a> out`,
			expected: `<p>Check <a href="https://example.com">this link</a> out</p>`,
		},
		{
			name:     "Paragraph with already processed side note",
			input:    `Text <label for="sn-sn1" class="margin-toggle sidenote-number"></label><input type="checkbox" id="sn-sn1" class="margin-toggle"><span class="sidenote">note here</span> continues`,
			expected: `<p>Text <label for="sn-sn1" class="margin-toggle sidenote-number"></label><input type="checkbox" id="sn-sn1" class="margin-toggle"><span class="sidenote">note here</span> continues</p>`,
		},
		{
			name:     "Paragraph with already processed italic",
			input:    "This has <em>italic</em> text",
			expected: `<p>This has <em>italic</em> text</p>`,
		},
		{
			name:     "Paragraph with already processed inline code",
			input:    "Use <code>fmt.Println()</code> function",
			expected: `<p>Use <code>fmt.Println()</code> function</p>`,
		},
		{
			name:     "Paragraph with multiple processed elements",
			input:    `This is <b>bold</b> and <em>italic</em> with <code>code</code> and <a href="url.com">link</a>`,
			expected: `<p>This is <b>bold</b> and <em>italic</em> with <code>code</code> and <a href="url.com">link</a></p>`,
		},
		{
			name:     "Paragraph with bold, link and side note",
			input:    `See <b>important</b> <a href="url.com">doc</a> <label for="sn-sn2" class="margin-toggle sidenote-number"></label><input type="checkbox" id="sn-sn2" class="margin-toggle"><span class="sidenote">details</span> here`,
			expected: `<p>See <b>important</b> <a href="url.com">doc</a> <label for="sn-sn2" class="margin-toggle sidenote-number"></label><input type="checkbox" id="sn-sn2" class="margin-toggle"><span class="sidenote">details</span> here</p>`,
		},
		{
			name:     "Paragraph with margin note",
			input:    `Text with <label for="mn-mn1" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-mn1" class="margin-toggle"><span class="marginnote">margin note</span> here`,
			expected: `<p>Text with <label for="mn-mn1" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-mn1" class="margin-toggle"><span class="marginnote">margin note</span> here</p>`,
		},
		{
			name:     "Plain paragraph",
			input:    "Just plain text",
			expected: `<p>Just plain text</p>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := Paragraph.Html(tt.input)
			
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
