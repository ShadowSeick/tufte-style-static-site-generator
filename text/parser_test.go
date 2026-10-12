package text

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	var parser Parser
	// Get the fiel
	filePath := filepath.Join("data", "testdata.md")
	source, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", filePath, err)
	}

	expectedTokens := []Token{
		&Header{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# "),
						offset: len("# "),
					},
					end: pos{
						line:   0,
						column: len("# Title"),
						offset: len("# Title"),
					},
				},
			},
			Level: 1,
		},
		&LineFeed{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title"),
						offset: len("# Title"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n"),
						offset: len("# Title\n"),
					},
				},
			},
		},
		&Subheader{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader]("),
						offset: len("# Title\n?[^subheader]("),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader"),
						offset: len("# Title\n?[^subheader](Subheader"),
					},
				},
			},
		},
		&LineFeed{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader"),
						offset: len("# Title\n?[^subheader](Subheader"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)"),
						offset: len("# Title\n?[^subheader](Subheader)"),
					},
				},
			},
		},
		&Text{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)"),
						offset: len("# Title\n?[^subheader](Subheader)"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuf"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuf"),
					},
				},
			},
		},
		&LineFeed{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuf"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuf"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff"),
					},
				},
			},
		},
		&Text{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here"),
					},
				},
			},
		},
		&Link{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link]"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link]"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com"),
					},
				},
			},
			Name: []byte("link"),
		},
		&Text{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin not"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin not"),
					},
				},
			},
		},
		&MarginNote{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note]"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note]"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive"),
					},
				},
			},
			children: []Token{
				&Text{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note]("),
								offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note]("),
							},
							end: pos{
								line:   0,
								column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image "),
								offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image "),
							},
						},
					},
				},
				&Image{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption]("),
								offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption]("),
							},
							end: pos{
								line:   0,
								column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg)"),
								offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg)"),
							},
						},
					},
					Caption: []byte("catpion"),
				},
				&Text{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg)"),
								offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg)"),
							},
							end: pos{
								line:   0,
								column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive"),
								offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive"),
							},
						},
					},
				},
			},
		},
		&LineFeed{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)"),
					},
				},
			},
		},
		&Header{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n##"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n##"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second titl"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second titl"),
					},
				},
			},
			Level: 2,
		},
		&LineFeed{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second titl"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second titl"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title"),
					},
				},
			},
		},
		&Text{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here "),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here "),
					},
				},
			},
		},
		&Bold{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold*"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold*"),
					},
				},
			},
		},
		&Text{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold*"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold*"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and "),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and "),
					},
				},
			},
		},
		&Italic{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**"),
					},
				},
			},
		},
		&Text{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**."),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**."),
					},
				},
			},
		},
		&LineFeed{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**."),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**."),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**.\n"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**.\n"),
					},
				},
			},
		},
		&Text{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**.\n"),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**.\n"),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**.\nFinal note here."),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**.\nFinal note here."),
					},
				},
			},
		},
		&LineFeed{
			baseToken: baseToken{
				content: segment{
					start: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**.\nFinal note here."),
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**.\nFinal note here."),
					},
					end: pos{
						line:   0,
						column: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**.\nFinal note here.") + 1,
						offset: len("# Title\n?[^subheader](Subheader)\nThis is the main stuff\nSome text here [link](https://example.com). Here it goes a margin note?[^margin-note](with some image ![caption](articles/images/something.jpg). This is quite impressive)\n## Second title\nThis is more text here *bold* and **italic**.\nFinal note here.") + 1,
					},
				},
			},
		},
	}
	tokens, errs := parser.Parse(source)
	if len(errs) != 0 {
		t.Errorf("expected to have 0 errors but got a lot %d", len(errs))
	}

	if len(tokens) != len(expectedTokens) {
		t.Errorf("expected to have %d tokens, but got %d", len(expectedTokens), len(tokens))
	}

	for i := range len(tokens) {
		token := tokens[i]
		expectedToken := expectedTokens[i]
		if token.String() != expectedToken.String() {
			t.Errorf("expected token to be %s but got %s", expectedTokens[i].String(), tokens[i].String())
		}

		if token.Start() != expectedToken.Start() || token.End() != expectedToken.End() {
			t.Errorf("expected token %s to start in %d and end in %d but started in %d and ended in %d", token.String(), expectedToken.Start(), expectedToken.End(), token.Start(), token.End())
		}
	}
}

func TestGetContentFromBalanced(t *testing.T) {
	tests := []struct {
		name                 string
		input                []byte
		start                int
		expectedRes          []byte
		expectedContentStart int
		expectedContentEnd   int
		err                  error
	}{
		{
			name:                 "correctly balanced link",
			input:                []byte("some text here [valid link](https://example.com) whatever"),
			start:                len("some text here "),
			expectedRes:          []byte("valid link"),
			expectedContentStart: len("some text here [valid link]("),
			expectedContentEnd:   len("some text here [valid link](https://example.com)"),
			err:                  nil,
		},
		{
			name:                 "correctly balanced link with wrong start should be treated as not a markdown structure",
			input:                []byte("[valid link](https://example.com)"),
			start:                1,
			expectedRes:          nil,
			expectedContentStart: 0,
			expectedContentEnd:   0,
			err:                  nil,
		},
		{
			name:                 "uncorrectly balanced structure",
			input:                []byte("[not valid}(../source/image.txt)"),
			start:                0,
			expectedRes:          nil,
			expectedContentStart: 0,
			expectedContentEnd:   0,
			err:                  ErrNotValidMarkdownStructureName,
		},
		{
			name:                 "correctly with a sub structure balanced structure",
			input:                []byte("[^some-random-thing](This is something else ![image here](../../source/image.png))"),
			start:                0,
			expectedRes:          []byte("^some-random-thing"),
			expectedContentStart: len("[^some-random-thing]("),
			expectedContentEnd:   len("[^some-random-thing](This is something else ![image here](../../source/image.png))"),
			err:                  nil,
		},
		{
			name:                 "uncorrectly added",
			input:                []byte("[^some-random-thing])This is something else ![image here](../../source/image.png))"),
			start:                0,
			expectedRes:          nil,
			expectedContentStart: 0,
			expectedContentEnd:   0,
			err:                  ErrNotValidMarkdownStructureContent,
		},
		{
			name:                 "correctly with uncorrect sub structure",
			input:                []byte("[^some-random-thing](This is something else ![image here}(../../source/image.png))"),
			start:                0,
			expectedRes:          []byte("^some-random-thing"),
			expectedContentStart: len("[^some-random-thing]("),
			expectedContentEnd:   len("[^some-random-thing](This is something else ![image here](../../source/image.png))"),
			err:                  nil,
		},
		{
			name:                 "uncorrectly terminated",
			input:                []byte("[^some-random-thing](This is something \nelse)"),
			start:                0,
			expectedRes:          nil,
			expectedContentStart: 0,
			expectedContentEnd:   0,
			err:                  ErrNotValidMarkdownStructureContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, contentStart, contentEnd, err := getContentFromBalanced(tt.start, tt.input)
			if !bytes.Equal(tt.expectedRes, res) {
				t.Errorf("expeced %s but got %s", tt.expectedRes, res)
			}

			if tt.expectedContentStart != contentStart {
				t.Errorf("expected content to start in %d idx, but it starts in %d", tt.expectedContentStart, contentStart)
			}

			if tt.expectedContentEnd != contentEnd {
				t.Errorf("expected content to end i %d idx, but it ends in %d", tt.expectedContentEnd, contentEnd)
			}

			if err != nil && !errors.Is(tt.err, err) {
				t.Errorf("expected content balanced to be %s but got %s", tt.err, err)
			}
		})
	}
}

func TestParseImage(t *testing.T) {
	tests := []struct {
		name          string
		source        []byte
		expectedToken Token
		expectedErr   error
	}{
		{
			name:   "correct image",
			source: []byte("![caption here](articles/img/image.png)"),
			expectedToken: &Image{
				baseToken: baseToken{
					content: segment{
						start: pos{
							line:   0,
							column: len("![caption here]("),
							offset: len("![caption here]("),
						},
						end: pos{
							line:   0,
							column: len("![caption here](articles/img/image.png)"),
							offset: len("![caption here](articles/img/image.png)"),
						},
					},
				},
				Caption: []byte("caption here"),
			},
		},
		{
			name:   "caption is always treat as text",
			source: []byte("![[a link to someplace](https://jeje.com)](articles/img/article1.png)"),
			expectedToken: &Image{
				baseToken: baseToken{
					content: segment{
						start: pos{
							line:   0,
							column: len("![[a link to someplace](https://jeje.com)]("),
							offset: len("![[a link to someplace](https://jeje.com)]("),
						},
						end: pos{
							line:   0,
							column: len("![[a link to someplace](https://jeje.com)](articles/img/article1.png)"),
							offset: len("![[a link to someplace](https://jeje.com)](articles/img/article1.png)"),
						},
					},
				},
				Caption: []byte("[a link to someplace](https://jeje.com)"),
			},
		},
	}

	var parser Parser
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.Reset()
			tokens, errs := parser.Parse(tt.source)
			if len(errs) != 0 && !errors.Is(tt.expectedErr, errs[0]) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, errs[0])
			}

			if len(tokens) != 1 {
				t.Errorf("expected to have one token, but got %d instead", len(tokens))
			}

			switch v := tokens[0].(type) {
			case *Image:
				if v.Start() != tt.expectedToken.Start() && v.End() != tt.expectedToken.End() {
					t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", tt.expectedToken.Start(), tt.expectedToken.End(), v.Start(), v.End())
				}
			default:
				t.Fatalf("expected to be an Image but got %s", v.String())
			}
		})
	}
}

func TestParseLink(t *testing.T) {
	tests := []struct {
		name          string
		source        []byte
		expectedToken Token
		expectedErr   error
	}{
		{
			name:   "correct link",
			source: []byte("[correct link](https://example.com)"),
			expectedToken: &Link{
				baseToken: baseToken{
					content: segment{
						start: pos{
							line:   0,
							column: len("[correct link]("),
							offset: len("[correct link]("),
						},
						end: pos{
							line:   0,
							column: len("[correct link](https://example.com)"),
							offset: len("[correct link](https://example.com)"),
						},
					},
				},
				Name: []byte("correct link"),
			},
			expectedErr: nil,
		},
		{
			name:   "name is always treat as text",
			source: []byte("[*something bold*](https://example.com)"),
			expectedToken: &Link{
				baseToken: baseToken{
					content: segment{
						start: pos{
							line:   0,
							column: len("[*something bold*]("),
							offset: len("[*something bold*]("),
						},
						end: pos{
							line:   0,
							column: len("[*something bold*](https://example.com)"),
							offset: len("[*something bold*](https://example.com)"),
						},
					},
				},
				Name: []byte("*something bold*"),
			},
			expectedErr: nil,
		},
	}

	var parser Parser
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.Reset()
			tokens, errs := parser.Parse(tt.source)
			if len(errs) != 0 && !errors.Is(tt.expectedErr, errs[0]) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, errs[0])
			}

			if len(tokens) != 1 {
				t.Errorf("expected to have one token, but got %d instead", len(tokens))
			}

			switch v := tokens[0].(type) {
			case *Link:
				if v.Start() != tt.expectedToken.Start() && v.End() != tt.expectedToken.End() {
					t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", tt.expectedToken.Start(), tt.expectedToken.End(), v.Start(), v.End())
				}
			default:
				t.Fatalf("expected to be an Image but got %s", v.String())
			}
		})
	}
}

func TestParseFontModifiers(t *testing.T) {
	tests := []struct {
		name          string
		source        []byte
		expectedToken []Token
		expectedErr   error
	}{
		{
			name:   "correct bold",
			source: []byte("*bold*"),
			expectedToken: []Token{
				&Bold{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: len("*"),
								offset: len("*"),
							},
							end: pos{
								line:   0,
								column: len("*bold*"),
								offset: len("*bold*"),
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:   "incorrect bold",
			source: []byte("*bold"),
			expectedToken: []Token{
				&Text{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: 0,
								offset: 0,
							},
							end: pos{
								line:   0,
								column: len("*bold"),
								offset: len("*bold"),
							},
						},
					},
				},
			},
			expectedErr: ErrFontModifiersNotValid,
		},
		{
			name:   "correct italic",
			source: []byte("**italic**"),
			expectedToken: []Token{
				&Italic{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: len("**"),
								offset: len("**"),
							},
							end: pos{
								line:   0,
								column: len("**italic**"),
								offset: len("**italic**"),
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:   "incorrect italic",
			source: []byte("**italic*"),
			expectedToken: []Token{
				&Text{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: 0,
								offset: 0,
							},
							end: pos{
								line:   0,
								column: len("**italic*"),
								offset: len("**italic*"),
							},
						},
					},
				},
			},
			expectedErr: ErrFontModifiersNotValid,
		},
		{
			name:   "incorrect italic by line feed",
			source: []byte("**italic\n"),
			expectedToken: []Token{
				&Text{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: 0,
								offset: 0,
							},
							end: pos{
								line:   0,
								column: len("**italic"),
								offset: len("**italic"),
							},
						},
					},
				},
				&LineFeed{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: len("**italic"),
								offset: len("**italic"),
							},
							end: pos{
								line:   0,
								column: len("**italic\n"),
								offset: len("**italic\n"),
							},
						},
					},
				},
			},
			expectedErr: ErrFontModifiersNotValid,
		},
	}

	var parser Parser
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.Reset()
			tokens, errs := parser.Parse(tt.source)
			if len(errs) != 0 && !errors.Is(tt.expectedErr, errs[0]) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, errs[0])
			}

			if len(tokens) != len(tt.expectedToken) {
				t.Errorf("expected to have %d tokens, but got %d instead", len(tt.expectedToken), len(tokens))
			}

			for i := range len(tokens) {
				if tokens[i].String() != tt.expectedToken[i].String() {
					t.Errorf("expecting token to be %s but got %s", tt.expectedToken[i].String(), tokens[0].String())
				}

				if tokens[i].Start() != tt.expectedToken[i].Start() || tokens[i].End() != tt.expectedToken[i].End() {
					t.Errorf("expected token %s to start in %d and end in %d but got start in %d and end in %d", tokens[i].String(), tt.expectedToken[i].Start(), tt.expectedToken[i].End(), tokens[i].Start(), tokens[i].End())
				}

				switch v := tokens[i].(type) {
				case *BlockCode:
					blockCode := tt.expectedToken[i].(*BlockCode)
					if !bytes.Equal(blockCode.Language, v.Language) {
						t.Errorf("expected block code to have language %s but got %s", blockCode.Language, v.Language)
					}
				}
			}
		})
	}
}

func TestParseHeader(t *testing.T) {
	tests := []struct {
		name           string
		source         []byte
		expectedTokens []Token
		expectedErr    error
	}{
		{
			name:   "correct header h1",
			source: []byte("# Header\n"),
			expectedTokens: []Token{
				&Header{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: len("# "),
								offset: len("# "),
							},
							end: pos{
								line:   0,
								column: len("# Header"),
								offset: len("# Header"),
							},
						},
					},
					Level: 1,
				},
				&LineFeed{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: len("# Header"),
								offset: len("# Header"),
							},
							end: pos{
								line:   0,
								column: len("# Header\n"),
								offset: len("# Header\n"),
							},
						},
					},
				},
			},
		},
		{
			name:   "correct header h2",
			source: []byte("## Header\n"),
			expectedTokens: []Token{
				&Header{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: len("## "),
								offset: len("## "),
							},
							end: pos{
								line:   0,
								column: len("## Header"),
								offset: len("## Header"),
							},
						},
					},
					Level: 2,
				},
				&Text{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: len("## Header"),
								offset: len("## Header"),
							},
							end: pos{
								line:   0,
								column: len("## Header\n"),
								offset: len("## Header\n"),
							},
						},
					},
				},
			},
		},
		{
			name:   "incorrect header",
			source: []byte("#Header"),
			expectedTokens: []Token{
				&Text{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: 0,
								offset: 0,
							},
							end: pos{
								line:   0,
								column: len("#Header"),
								offset: len("#Header"),
							},
						},
					},
				},
			},
			expectedErr: ErrHeaderNotValid,
		},
	}

	var parser Parser
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.Reset()
			tokens, errs := parser.Parse(tt.source)
			if len(errs) != 0 && !errors.Is(tt.expectedErr, errs[0]) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, errs[0])
			}

			if len(tokens) != len(tt.expectedTokens) {
				t.Errorf("expected to have %d tokens, but got %d instead", len(tt.expectedTokens), len(tokens))
			}

			for i := range len(tokens) {
				token := tokens[i]
				expectedToken := tt.expectedTokens[i]

				if token.Start() != expectedToken.Start() || token.End() != expectedToken.End() {
					t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", expectedToken.Start(), expectedToken.End(), token.Start(), token.End())
				}

				switch v := token.(type) {
				case *Header:
					expectedHeader := expectedToken.(*Header)
					if v.Level != expectedHeader.Level {
						t.Errorf("expected header to have level %d but got %d", v.Level, expectedHeader.Level)
					}
				}

			}
		})
	}
}

func TestParseSubheader(t *testing.T) {
	tests := []struct {
		name          string
		source        []byte
		expectedToken Token
		expectedErr   error
	}{
		{
			name:   "correct subheader",
			source: []byte("?[^subheader](Subheader)"),
			expectedToken: &Subheader{
				baseToken: baseToken{
					content: segment{
						start: pos{
							line:   0,
							column: len("?[^subheader]("),
							offset: len("?[^subheader]("),
						},
						end: pos{
							line:   0,
							column: len("?[^subheader](Subheader)"),
							offset: len("?[^subheader](Subheader)"),
						},
					},
				},
			},
		},
	}

	var parser Parser
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.Reset()
			tokens, errs := parser.Parse(tt.source)
			if len(errs) != 0 && !errors.Is(tt.expectedErr, errs[0]) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, errs[0])
			}

			if len(tokens) != 1 {
				t.Errorf("expected to have one token, but got %d instead", len(tokens))
			}

			switch v := tokens[0].(type) {
			case *Subheader:
				if v.Start() != tt.expectedToken.Start() && v.End() != tt.expectedToken.End() {
					t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", tt.expectedToken.Start(), tt.expectedToken.End(), v.Start(), v.End())
				}
			case *Text:
				if v.Start() != tt.expectedToken.Start() && v.End() != tt.expectedToken.End() {
					t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", tt.expectedToken.Start(), tt.expectedToken.End(), v.Start(), v.End())
				}
			default:
				t.Errorf("expecting token to be %s but got %s", tt.expectedToken.String(), v.String())
			}
		})
	}
}

func TestParseText(t *testing.T) {
	tests := []struct {
		name          string
		source        []byte
		expectedToken Token
		expectedErr   error
	}{
		{
			name:   "simple text",
			source: []byte("this is just text"),
			expectedToken: &Text{
				baseToken: baseToken{
					content: segment{
						start: pos{
							line:   0,
							column: 0,
							offset: 0,
						},
						end: pos{
							line:   0,
							column: len("this is just text"),
							offset: len("this is just text"),
						},
					},
				},
			},
			expectedErr: ErrNotValidMarkdown,
		},
	}

	var parser Parser
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.Reset()
			tokens, errs := parser.Parse(tt.source)
			if len(errs) != 0 && !errors.Is(tt.expectedErr, errs[0]) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, errs[0])
			}

			if len(tokens) != 1 {
				t.Errorf("expected to have one token, but got %d instead", len(tokens))
			}

			text, ok := tokens[0].(*Text)
			if !ok {
				t.Errorf("expecting token to be %s but got %s", tt.expectedToken.String(), text.String())
			}

			if text.Start() != tt.expectedToken.Start() && text.End() != tt.expectedToken.End() {
				t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", tt.expectedToken.Start(), tt.expectedToken.End(), text.Start(), text.End())
			}
		})
	}
}

func TestParseCode(t *testing.T) {
	tests := []struct {
		name          string
		source        []byte
		expectedToken []Token
		expectedErr   error
	}{
		{
			name:   "correct inline code",
			source: []byte("`this is inlined code`"),
			expectedToken: []Token{
				&InlineCode{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: len("`"),
								offset: len("`"),
							},
							end: pos{
								line:   0,
								column: len("`this is inlined code`"),
								offset: len("`this is inlined code`"),
							},
						},
					},
				},
			},
		},
		{
			name:   "correct block code",
			source: []byte("```go\nfunc() string {\n  return \"I am indeed\"\n}```"),
			expectedToken: []Token{
				&BlockCode{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: len("```go"),
								offset: len("```go"),
							},
							end: pos{
								line:   0,
								column: len("```go\nfunc() string {\n  return \"I am indeed\"\n}```"),
								offset: len("```go\nfunc() string {\n  return \"I am indeed\"\n}```"),
							},
						},
					},
					Language: []byte("go"),
				},
			},
		},
		{
			name:   "incorrect inline code",
			source: []byte("``this is inlined code``"),
			expectedToken: []Token{
				&Text{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: 0,
								offset: 0,
							},
							end: pos{
								line:   0,
								column: len("``this is inlined code"),
								offset: len("``this is inlined code"),
							},
						},
					},
				},
				&Text{
					baseToken: baseToken{
						content: segment{
							start: pos{
								line:   0,
								column: len("``this is inlined code"),
								offset: len("``this is inlined code"),
							},
							end: pos{
								line:   0,
								column: len("``this is inlined code``"),
								offset: len("``this is inlined code``"),
							},
						},
					},
				},
			},
			expectedErr: ErrCodeNotValid,
		},
	}

	var parser Parser
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.Reset()
			tokens, errs := parser.Parse(tt.source)
			if len(errs) != 0 && !errors.Is(tt.expectedErr, errs[0]) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, errs[0])
			}

			if len(tokens) != len(tt.expectedToken) {
				t.Errorf("expected to have %d tokens, but got %d instead", len(tt.expectedToken), len(tokens))
			}

			for i := range len(tokens) {
				if tokens[i].String() != tt.expectedToken[i].String() {
					t.Errorf("expecting token to be %s but got %s", tt.expectedToken[i].String(), tokens[0].String())
				}

				if tokens[i].Start() != tt.expectedToken[i].Start() && tokens[i].End() != tt.expectedToken[i].End() {
					t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", tt.expectedToken[i].Start(), tt.expectedToken[i].End(), tokens[i].Start(), tokens[i].End())
				}

				switch v := tokens[i].(type) {
				case *BlockCode:
					blockCode := tt.expectedToken[i].(*BlockCode)
					if !bytes.Equal(blockCode.Language, v.Language) {
						t.Errorf("expected block code to have language %s but got %s", blockCode.Language, v.Language)
					}
				}
			}
		})
	}
}

func TestParseMarginNote(t *testing.T) {
	tests := []struct {
		name          string
		source        []byte
		expectedToken CustomToken
		expectedErr   error
	}{
		{
			name:   "correct margin note",
			source: []byte("?[^margin-note](*bold* but **italic** with a [link](https://example.com))"),
			expectedToken: &MarginNote{
				baseToken: baseToken{
					content: segment{
						start: pos{
							line:   0,
							column: len("?[^margin-note]("),
							offset: len("?[^margin-note]("),
						},
						end: pos{
							line:   0,
							column: len("?[^margin-note](*bold* but **italic** with a [link](https://example.com))"),
							offset: len("?[^margin-note](*bold* but **italic** with a [link](https://example.com))"),
						},
					},
				},
				children: []Token{
					&Bold{
						baseToken: baseToken{
							content: segment{
								start: pos{
									line:   0,
									column: len("?[^margin-note](*"),
									offset: len("?[^margin-note](*"),
								},
								end: pos{
									line:   0,
									column: len("?[^margin-note](*bold*"),
									offset: len("?[^margin-note](*bold*"),
								},
							},
						},
					},
					&Text{
						baseToken: baseToken{
							content: segment{
								start: pos{
									line:   0,
									column: len("?[^margin-note](*bold*"),
									offset: len("?[^margin-note](*bold*"),
								},
								end: pos{
									line:   0,
									column: len("?[^margin-note](*bold* but "),
									offset: len("?[^margin-note](*bold* but "),
								},
							},
						},
					},
					&Italic{
						baseToken: baseToken{
							content: segment{
								start: pos{
									line:   0,
									column: len("?[^margin-note](*bold* but **"),
									offset: len("?[^margin-note](*bold* but **"),
								},
								end: pos{
									line:   0,
									column: len("?[^margin-note](*bold* but **italic**"),
									offset: len("?[^margin-note](*bold* but **italic**"),
								},
							},
						},
					},
					&Text{
						baseToken: baseToken{
							content: segment{
								start: pos{
									line:   0,
									column: len("?[^margin-note](*bold* but **italic**"),
									offset: len("?[^margin-note](*bold* but **italic**"),
								},
								end: pos{
									line:   0,
									column: len("?[^margin-note](*bold* but **italic** with a "),
									offset: len("?[^margin-note](*bold* but **italic** with a "),
								},
							},
						},
					},
					&Link{
						baseToken: baseToken{
							content: segment{
								start: pos{
									line:   0,
									column: len("?[^margin-note](*bold* but **italic** with a [link]("),
									offset: len("?[^margin-note](*bold* but **italic** with a [link]("),
								},
								end: pos{
									line:   0,
									column: len("?[^margin-note](*bold* but **italic** with a [link](https://example.com)"),
									offset: len("?[^margin-note](*bold* but **italic** with a [link](https://example.com)"),
								},
							},
						},
					},
				},
			},
			expectedErr: ErrNotValidMarkdown,
		},
	}

	var parser Parser
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.Reset()
			tokens, errs := parser.Parse(tt.source)
			if len(errs) != 0 && !errors.Is(tt.expectedErr, errs[0]) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, errs[0])
			}

			if len(tokens) != 1 {
				t.Errorf("expected to have one token, but got %d instead", len(tokens))
			}

			marginNote, ok := tokens[0].(*MarginNote)
			if !ok {
				t.Errorf("expected token to be %s but got %s", tt.expectedToken.String(), tokens[0].String())
			}

			if marginNote.String() != tt.expectedToken.String() {
				t.Errorf("expected token to be %s but got %s", tt.expectedToken.String(), marginNote.String())
			}

			expectedChildren := tt.expectedToken.GetChildren()
			children := marginNote.GetChildren()
			if len(children) != len(expectedChildren) {
				t.Errorf("expected having %d children, instead got %d", len(expectedChildren), len(marginNote.children))
			}

			for i := range len(children) {
				child := children[i]
				expectedChild := expectedChildren[i]
				if child.String() != expectedChild.String() {
					t.Errorf("expected token to be %s, but got %s", expectedChild.String(), child.String())
				}

				if child.Start() != expectedChild.Start() || child.End() != expectedChild.End() {
					t.Errorf(
						"expected token %s and got %s;\n Expected to start in %d but started in %d; Expected to end in %d but ended in %d",
						expectedChild.String(), child.String(), expectedChild.Start(), child.Start(), expectedChild.End(), child.End())
				}
			}

			if marginNote.Start() != tt.expectedToken.Start() && marginNote.End() != tt.expectedToken.End() {
				t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", tt.expectedToken.Start(), tt.expectedToken.End(), marginNote.Start(), marginNote.End())
			}
		})
	}
}

func TestParseSideNote(t *testing.T) {
	tests := []struct {
		name          string
		source        []byte
		expectedToken CustomToken
		expectedErr   error
	}{
		{
			name:   "correct side note",
			source: []byte("?[^side-note](*bold* but **italic** with a [link](https://example.com))"),
			expectedToken: &SideNote{
				baseToken: baseToken{
					content: segment{
						start: pos{
							line:   0,
							column: len("?[^side-note]("),
							offset: len("?[^side-note]("),
						},
						end: pos{
							line:   0,
							column: len("?[^side-note](https://example.com))"),
							offset: len("?[^side-note](https://example.com))"),
						},
					},
				},
				children: []Token{
					&Bold{
						baseToken: baseToken{
							content: segment{
								start: pos{
									line:   0,
									column: len("?[^side-note](*"),
									offset: len("?[^side-note](*"),
								},
								end: pos{
									line:   0,
									column: len("?[^side-note](*bold*"),
									offset: len("?[^side-note](*bold*"),
								},
							},
						},
					},
					&Text{
						baseToken: baseToken{
							content: segment{
								start: pos{
									line:   0,
									column: len("?[^side-note](*bold*"),
									offset: len("?[^side-note](*bold*"),
								},
								end: pos{
									line:   0,
									column: len("?[^side-note](*bold* but"),
									offset: len("?[^side-note](*bold* but"),
								},
							},
						},
					},
					&Italic{
						baseToken: baseToken{
							content: segment{
								start: pos{
									line:   0,
									column: len("?[^side-note](*bold* but **"),
									offset: len("?[^side-note](*bold* but **"),
								},
								end: pos{
									line:   0,
									column: len("?[^side-note](*bold* but **italic**"),
									offset: len("?[^side-note](*bold* but **italic**"),
								},
							},
						},
					},
					&Text{
						baseToken: baseToken{
							content: segment{
								start: pos{
									line:   0,
									column: len("?[^side-note](*bold* but **italic**"),
									offset: len("?[^side-note](*bold* but **italic**"),
								},
								end: pos{
									line:   0,
									column: len("?[^side-note](*bold* but **italic** with a"),
									offset: len("?[^side-note](*bold* but **italic** with a"),
								},
							},
						},
					},
					&Link{
						baseToken: baseToken{
							content: segment{
								start: pos{
									line:   0,
									column: len("?[^side-note](*bold* but **italic** with a [link]("),
									offset: len("?[^side-note](*bold* but **italic** with a [link]("),
								},
								end: pos{
									line:   0,
									column: len("?[^side-note](*bold* but **italic** with a [link](https://example.com)"),
									offset: len("?[^side-note](*bold* but **italic** with a [link](https://example.com)"),
								},
							},
						},
					},
				},
			},
			expectedErr: ErrNotValidMarkdown,
		},
	}

	var parser Parser
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.Reset()
			tokens, errs := parser.Parse(tt.source)
			if len(errs) != 0 && !errors.Is(tt.expectedErr, errs[0]) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, errs[0])
			}

			if len(tokens) != 1 {
				t.Errorf("expected to have one token, but got %d instead", len(tokens))
			}

			sideNote, ok := tokens[0].(*SideNote)
			if !ok {
				t.Errorf("expected token to be %s but got %s", tt.expectedToken.String(), tokens[0].String())
			}

			if sideNote.String() != tt.expectedToken.String() {
				t.Errorf("expected token to be %s but got %s", tt.expectedToken.String(), sideNote.String())
			}

			expectedChildren := tt.expectedToken.GetChildren()
			children := sideNote.GetChildren()
			if len(children) != len(expectedChildren) {
				t.Errorf("expected having %d children, instead got %d", len(expectedChildren), len(sideNote.children))
			}

			for i := range len(tokens) {
				child := children[i]
				expectedChild := expectedChildren[i]
				if child.String() != expectedChild.String() {
					t.Errorf("expected token to be %s, but got %s", expectedChild.String(), child.String())
				}
				if child.Start() != expectedChild.Start() || child.End() != expectedChild.End() {
					t.Errorf(
						"expected token %s and got %s;\n Expected to start in %d but started in %d; Expected to end in %d but ended in %d",
						expectedChild.String(), child.String(), expectedChild.Start(), child.Start(), expectedChild.End(), child.End())
				}
			}

			if sideNote.Start() != tt.expectedToken.Start() && sideNote.End() != tt.expectedToken.End() {
				t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", tt.expectedToken.Start(), tt.expectedToken.End(), sideNote.Start(), sideNote.End())
			}
		})
	}
}
