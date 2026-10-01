package text

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

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
			expectedContentEnd:   len("some text here [valid link](https://example.com"),
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
			expectedContentEnd:   len("[^some-random-thing](This is something else ![image here](../../source/image.png)"),
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
			expectedContentEnd:   len("[^some-random-thing](This is something else ![image here](../../source/image.png)"),
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
			source: []byte("[caption here](articles/img/image.png)"),
			expectedToken: &Image{
				baseToken: baseToken{
					Content: Segment{
						start: pos{
							line:   0,
							column: len("[caption here]("),
							offset: len("[caption here]("),
						},
						end: pos{
							line:   0,
							column: len("articles/img/image.png)"),
							offset: len("articles/img/image.png)"),
						},
					},
				},
				Caption: []byte("caption here"),
			},
			expectedErr: nil,
		},
		{
			name:   "caption is always treat as text",
			source: []byte("[[a link to someplace](https://jeje.com)](articles/img/article1.png)"),
			expectedToken: &Image{
				baseToken: baseToken{
					Content: Segment{
						start: pos{
							line:   0,
							column: len("[[a link to someplace](https://jeje.com)]("),
							offset: len("[[a link to someplace](https://jeje.com)]("),
						},
						end: pos{
							line:   0,
							column: len("articles/img/article1.png)"),
							offset: len("articles/img/article1.png)"),
						},
					},
				},
				Caption: []byte("[a link to someplace](https://jeje.com)"),
			},
			expectedErr: nil,
		},
	}

	var parser Parser
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.Reset(tt.source)
			token, err := parser.parseImage()
			if err != nil && !errors.Is(tt.expectedErr, err) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, err)
			}

			if token.Start() != tt.expectedToken.Start() && token.End() != tt.expectedToken.End() {
				t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", tt.expectedToken.Start(), tt.expectedToken.End(), token.Start(), token.End())
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
					Content: Segment{
						start: pos{
							line:   0,
							column: len("[correct link]("),
							offset: len("[correct link]("),
						},
						end: pos{
							line:   0,
							column: len("https://example.com)"),
							offset: len("https://example.com)"),
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
					Content: Segment{
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
			parser.Reset(tt.source)
			token, err := parser.parseLink()
			if err != nil && !errors.Is(tt.expectedErr, err) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, err)
			}

			if token.Start() != tt.expectedToken.Start() && token.End() != tt.expectedToken.End() {
				t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", tt.expectedToken.Start(), tt.expectedToken.End(), token.Start(), token.End())
			}
		})
	}
}

func TestFontModifiers(t *testing.T) {
	tests := []struct {
		name          string
		source        []byte
		expectedToken Token
		expectedErr   error
	}{
		{
			name:   "correct bold",
			source: []byte("*bold*"),
			expectedToken: &Bold{
				baseToken: baseToken{
					Content: Segment{
						start: pos{
							line:   0,
							column: len("*"),
							offset: len("*"),
						},
						end: pos{
							line:   0,
							column: len("*bold"),
							offset: len("*bold"),
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:   "correct italic",
			source: []byte("**italic**"),
			expectedToken: &Italic{
				baseToken: baseToken{
					Content: Segment{
						start: pos{
							line:   0,
							column: len("**"),
							offset: len("**"),
						},
						end: pos{
							line:   0,
							column: len("**italic"),
							offset: len("**italic"),
						},
					},
				},
			},
			expectedErr: nil,
		},
	}

	var parser Parser
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.Reset(tt.source)
			token, err := parser.parseFontModifiers()
			if err != nil && !errors.Is(tt.expectedErr, err) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, err)
			}

			if token.Start() != tt.expectedToken.Start() && token.End() != tt.expectedToken.End() {
				t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", tt.expectedToken.Start(), tt.expectedToken.End(), token.Start(), token.End())
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
			source: []byte("[^margin-note](*bold* but **italic** with a [link](https://example.com))"),
			expectedToken: &MarginNote{
				baseToken: baseToken{
					Content: Segment{
						start: pos{
							line:   0,
							column: len("[^margin-note]("),
							offset: len("[^margin-note]("),
						},
						end: pos{
							line:   0,
							column: len(""),
							offset: len("https://example.com)"),
						},
					},
				},
				children: []Token{
					&Bold{
						baseToken: baseToken{
							Content: Segment{
								start: pos{
									line:   0,
									column: len("[^margin-note](*"),
									offset: len("[^margin-note](*"),
								},
								end: pos{
									line:   0,
									column: len("[^margin-note](*bold"),
									offset: len("[^margin-note](*bold"),
								},
							},
						},
					},
					&Text{
						baseToken: baseToken{
							Content: Segment{
								start: pos{
									line:   0,
									column: len("[^margin-note](*bold*"),
									offset: len("[^margin-note](*bold*"),
								},
								end: pos{
									line:   0,
									column: len("[^margin-note](*bold* but"),
									offset: len("[^margin-note](*bold* but"),
								},
							},
						},
					},
					&Italic{
						baseToken: baseToken{
							Content: Segment{
								start: pos{
									line:   0,
									column: len("[^margin-note](*bold* but **"),
									offset: len("[^margin-note](*bold* but **"),
								},
								end: pos{
									line:   0,
									column: len("[^margin-note](*bold* but **italic"),
									offset: len("[^margin-note](*bold* but **italic"),
								},
							},
						},
					},
					&Text{
						baseToken: baseToken{
							Content: Segment{
								start: pos{
									line:   0,
									column: len("[^margin-note](*bold* but **italic**"),
									offset: len("[^margin-note](*bold* but **italic**"),
								},
								end: pos{
									line:   0,
									column: len("[^margin-note](*bold* but **italic** with a"),
									offset: len("[^margin-note](*bold* but **italic** with a"),
								},
							},
						},
					},
					&Link{
						baseToken: baseToken{
							Content: Segment{
								start: pos{
									line:   0,
									column: len("[^margin-note](*bold* but **italic** with a [link]("),
									offset: len("[^margin-note](*bold* but **italic** with a [link]("),
								},
								end: pos{
									line:   0,
									column: len("[^margin-note](*bold* but **italic** with a [link](https://example.com"),
									offset: len("[^margin-note](*bold* but **italic** with a [link](https://example.com"),
								},
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
	}

	var parser Parser
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.Reset(tt.source)
			token, err := parser.parseCustom()
			if err != nil && !errors.Is(tt.expectedErr, err) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, err)
			}

			if token.Start() != tt.expectedToken.Start() && token.End() != tt.expectedToken.End() {
				t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", tt.expectedToken.Start(), tt.expectedToken.End(), token.Start(), token.End())
			}

			customToken, ok := token.(*MarginNote)
			if !ok {
				t.Fatalf("it should be a Margin note: %s", err.Error())
			}

			expectedChildren := tt.expectedToken.GetChildren()
			children := customToken.GetChildren()
			if len(children) != len(expectedChildren) {
				t.Errorf("expected having %d children, instead got %d", len(expectedChildren), len(customToken.children))
			}

			for i := 0; i < len(children); i++ {
				child := children[i]
				expectedChild := expectedChildren[i]
				if child.Start() != expectedChild.Start() && child.End() != expectedChild.End() {
					t.Errorf(
						"expected token %s and got %s;\n Expected to start in %d but started in %d; Expected to end in %d but ended in %d",
						expectedChild.String(), child.String(), expectedChild.Start(), child.Start(), expectedChild.End(), child.End())
				}
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
			source: []byte("[^side-note](*bold* but **italic** with a [link](https://example.com))"),
			expectedToken: &MarginNote{
				baseToken: baseToken{
					Content: Segment{
						start: pos{
							line:   0,
							column: len("[^side-note]("),
							offset: len("[^side-note]("),
						},
						end: pos{
							line:   0,
							column: len(""),
							offset: len("https://example.com)"),
						},
					},
				},
				children: []Token{
					&Bold{
						baseToken: baseToken{
							Content: Segment{
								start: pos{
									line:   0,
									column: len("[^side-note](*"),
									offset: len("[^side-note](*"),
								},
								end: pos{
									line:   0,
									column: len("[^side-note](*bold"),
									offset: len("[^side-note](*bold"),
								},
							},
						},
					},
					&Text{
						baseToken: baseToken{
							Content: Segment{
								start: pos{
									line:   0,
									column: len("[^side-note](*bold*"),
									offset: len("[^side-note](*bold*"),
								},
								end: pos{
									line:   0,
									column: len("[^side-note](*bold* but"),
									offset: len("[^side-note](*bold* but"),
								},
							},
						},
					},
					&Italic{
						baseToken: baseToken{
							Content: Segment{
								start: pos{
									line:   0,
									column: len("[^side-note](*bold* but **"),
									offset: len("[^side-note](*bold* but **"),
								},
								end: pos{
									line:   0,
									column: len("[^side-note](*bold* but **italic"),
									offset: len("[^side-note](*bold* but **italic"),
								},
							},
						},
					},
					&Text{
						baseToken: baseToken{
							Content: Segment{
								start: pos{
									line:   0,
									column: len("[^side-note](*bold* but **italic**"),
									offset: len("[^side-note](*bold* but **italic**"),
								},
								end: pos{
									line:   0,
									column: len("[^side-note](*bold* but **italic** with a"),
									offset: len("[^side-note](*bold* but **italic** with a"),
								},
							},
						},
					},
					&Link{
						baseToken: baseToken{
							Content: Segment{
								start: pos{
									line:   0,
									column: len("[^side-note](*bold* but **italic** with a [link]("),
									offset: len("[^side-note](*bold* but **italic** with a [link]("),
								},
								end: pos{
									line:   0,
									column: len("[^side-note](*bold* but **italic** with a [link](https://example.com"),
									offset: len("[^side-note](*bold* but **italic** with a [link](https://example.com"),
								},
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
	}

	var parser Parser
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.Reset(tt.source)
			token, err := parser.parseCustom()
			if err != nil && !errors.Is(tt.expectedErr, err) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, err)
			}

			if token.Start() != tt.expectedToken.Start() && token.End() != tt.expectedToken.End() {
				t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", tt.expectedToken.Start(), tt.expectedToken.End(), token.Start(), token.End())
			}

			customToken, ok := token.(*SideNote)
			if !ok {
				t.Fatalf("it should be a Side note: %s", err.Error())
			}

			expectedChildren := tt.expectedToken.GetChildren()
			children := customToken.GetChildren()
			if len(children) != len(expectedChildren) {
				t.Errorf("expected having %d children, instead got %d", len(expectedChildren), len(customToken.children))
			}

			for i := 0; i < len(children); i++ {
				child := children[i]
				expectedChild := expectedChildren[i]
				if child.Start() != expectedChild.Start() && child.End() != expectedChild.End() {
					t.Errorf(
						"expected token %s and got %s;\n Expected to start in %d but started in %d; Expected to end in %d but ended in %d",
						expectedChild.String(), child.String(), expectedChild.Start(), child.Start(), expectedChild.End(), child.End())
				}
			}
		})
	}
}

func TestParseTitle(t *testing.T) {
	tests := []struct {
		name          string
		source        []byte
		expectedToken Token
		expectedErr   error
	}{
		{
			name:   "correct title h1",
			source: []byte("# Title"),
			expectedToken: &Header{
				baseToken: baseToken{
					Content: Segment{
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
		},
		{
			name:   "correct title h2",
			source: []byte("## Title"),
			expectedToken: &Header{
				baseToken: baseToken{
					Content: Segment{
						start: pos{
							line:   0,
							column: len("## "),
							offset: len("## "),
						},
						end: pos{
							line:   0,
							column: len("## Title"),
							offset: len("## Title"),
						},
					},
				},
				Level: 2,
			},
		},
		{
			name:   "incorrect title",
			source: []byte("#Title\n"),
			expectedToken: &Text{
				baseToken: baseToken{
					Content: Segment{
						start: pos{
							line:   0,
							column: len(""),
							offset: len(""),
						},
						end: pos{
							line:   0,
							column: len("#Title"),
							offset: len("#Title"),
						},
					},
				},
			},
		},
	}

	var parser Parser
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.Reset(tt.source)
			token, err := parser.parseTitle()
			fmt.Printf(
				"expected token to start in %d and end in %d but got start in %d and end in %d\n", tt.expectedToken.Start(), tt.expectedToken.End(), token.Start(), token.End(),
			)
			if err != nil && !errors.Is(tt.expectedErr, err) {
				t.Errorf("expected error to be %s but got %s", tt.expectedErr, err)
			}

			if token.String() != tt.expectedToken.String() {
				t.Errorf("expected token %s, but got %s", tt.expectedToken.String(), token.String())
			}

			if token.Start() != tt.expectedToken.Start() && token.End() != tt.expectedToken.End() {
				t.Errorf("expected token to start in %d and end in %d but got start in %d and end in %d", tt.expectedToken.Start(), tt.expectedToken.End(), token.Start(), token.End())
			}

			switch v := token.(type) {
			case *Header:
				expected, ok := tt.expectedToken.(*Header)
				if !ok {
					t.Fatalf("token not expected")
				}
				if v.Level != expected.Level {
					t.Errorf("expected token level to be %d but got %d", expected.Level, v.Level)
				}
			}
		})
	}
}
