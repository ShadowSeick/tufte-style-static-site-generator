package text

import (
	"errors"
	"fmt"
	"slices"
)

var (
	ErrHeaderNotValid                   = errors.New("not valid header")
	ErrFontModifiersNotValid            = errors.New("not valid font modifiers")
	ErrImageNotValid                    = errors.New("not valid image")
	ErrLinkNotValid                     = errors.New("not valid link")
	ErrNotValidMarkdown                 = errors.New("not valid markdown")
	ErrCustomMarkdownNotValid           = errors.New("not valid custom markdown")
	ErrMarkdownTypeNotImplemented       = errors.New("markdown type not implemented")
	ErrNotValidMarkdownStructureName    = errors.New("not valid markdown balanced structure []() name")
	ErrNotValidMarkdownStructureContent = errors.New("not valid markdown balanced structure []() content")
)

const (
	title         = '#'
	custom        = '?'
	link          = '['
	code          = '`'
	image         = '!'
	fontModifiers = '*'
	breakLine     = '\n'

	space = ' '

	subheaderIdentifier  = "^subheader"
	marginNoteIdentifier = "^margin-note"
	sideNoteIdentifier   = "^side-note"
)

type Parser struct {
	source []byte
	offset int
	column int
	line   int
}

func (p *Parser) Reset() {
	p.offset = 0
	p.column = 0
	p.line = 0
}

// Parse can be much simpler. For more information why http://number-none.com/blow/blog/programming/2014/09/26/carmack-on-inlined-code.html
// - I am duplicating a lot of variables and logic
// - It needs a lot of implicit knowledge
// - Not really straightforward
// - A lot of jumping
//
// One example is that I am using parseChildren to get the tokens from it, it is not necessary, I can use a parse function that accepts a source in bytes
// Once I have this, I can recursively call it to parse the children for a specific part of the source
// There are many like that, but mostly is the tendency of overusing parseText and adding the same logic many times
//
// After the tests have been done, I will refactor heavily to make it smaller, simpler and more maintainable
func (p *Parser) Parse(source []byte) ([]Token, []error) {
	var errs []error
	var tokens []Token
	var content int
	var token Token
	var err error
	for p.offset < len(source) {
	out:
		switch source[p.offset] {
		case title:
			if p.column != 0 {
				break
			}
			const maxLevel = 2
			var level int
			for content < maxLevel {
				if p.peek(source, content) == title {
					level++
				} else {
					break
				}
				content++
			}

			if p.peek(source, content) != space {
				err = ErrHeaderNotValid
				break
			}

			content++ // Content starts after space
			for range content {
				p.advance()
			}

			token = &Header{
				Level: level,
			}

			for source[content] != breakLine {
				content++
			}
		case custom:
			identifier, start, end, balancedErr := getContentFromBalanced(p.offset+1, source)
			if balancedErr != nil {
				err = fmt.Errorf("%w: %w", ErrCustomMarkdownNotValid, balancedErr)
				break
			}

			var customToken CustomToken
			id := string(identifier)
			switch id {
			case subheaderIdentifier:
				token = &Subheader{}
			case marginNoteIdentifier:
				customToken = &MarginNote{}
			case sideNoteIdentifier:
				customToken = &SideNote{}
			default:
				err = ErrCustomMarkdownNotValid
				break out
			}

			// Go to the start of the content
			for range start {
				p.advance()
			}

			content = end - start
			if customToken != nil {
				children, childErrs := p.Parse(source[p.offset : p.offset+content])
				if len(childErrs) != 0 {
					for _, childErr := range childErrs {
						errs = append(errs, fmt.Errorf("%w: %w", err, childErr))
					}
				}
				customToken.SetChildren(children)
				token = customToken
			}
		case link:
			name, start, end, balancedErr := getContentFromBalanced(p.offset, source)
			if balancedErr != nil {
				err = fmt.Errorf("%w: %w", ErrLinkNotValid, balancedErr)
				break
			}

			token = &Link{
				Name: name,
			}

			for range start + 1 {
				p.advance()
			}

			content = end - start
		case image:
			caption, start, end, balancedErr := getContentFromBalanced(p.offset+1, source)
			if balancedErr != nil {
				err = fmt.Errorf("%w: %w", ErrImageNotValid, balancedErr)
				break
			}

			token = &Image{
				Caption: caption,
			}

			for range start + 1 {
				p.advance()
			}

			content = end - start
		case code:
			// Parse code
		case fontModifiers:
			terminalChars := 1
			if p.peek(source, 1) != fontModifiers {
				token = &Bold{}
			} else {
				token = &Italic{}
				terminalChars++
			}
			content = terminalChars

			for range content {
				p.advance()
			}

			for i := 0; i < terminalChars; {
				content++

				if isAtEOF(source, p.offset+content) {
					break
				}

				if p.peek(source, content) == '\n' {
					err = ErrFontModifiersNotValid
					break out
				}

				if p.peek(source, content) == fontModifiers {
					i++
				}
			}

		case breakLine:
			token = &Jump{}
			p.line++
			p.column = 0
		default:
			err = ErrNotValidMarkdown
		}

		if err != nil {
			errs = append(errs, err)
			token = &Text{}
			for !isAtEOF(source, p.offset+content) && !slices.Contains([]byte{title, custom, link, code, image, fontModifiers, breakLine}, p.peek(source, content)) {
				content++
			}
		}

		token.SetContent(pos{
			line:   p.line,
			column: p.column,
			offset: p.offset,
		}, pos{
			line:   p.line,
			column: p.column + content,
			offset: p.offset + content,
		})

		tokens = append(tokens, token)

		// Advance the content to the last byte from the token content
		for range content {
			p.advance()
		}

		// Go to the next byte
		if !isAtEOF(source, p.offset) {
			p.advance()
		}

		content = 0
		token = nil
		err = nil
	}
	return tokens, errs
}

func (p *Parser) advance() {
	p.offset++
	p.column++
}

func (p *Parser) peek(source []byte, n int) byte {
	return source[p.offset+n]
}

func isAtEOF(source []byte, count int) bool {
	return count >= len(source)
}

// getContentFromBalanced expects the start idx of the markdown structure and the markdown source in bytes array
// It will return the name in a byte array and the start content idx and end idx of the markdown structure.
// A valid markdown structure must be balanced []() and not have any jumps between it; []('\n') this is not allowed
// In case of not being valid, it will return an error
func getContentFromBalanced(start int, source []byte) ([]byte, int, int, error) {
	// It is not a markdown structure
	if source[start] != '[' {
		return nil, 0, 0, nil
	}

	var content int
	var depth int
	var nameIdx int
nameLoop:
	for content < len(source) {
		switch source[content] {
		case '[':
			depth++
		case '\n':
			break nameLoop
		case ']':
			depth--
			if depth == 0 {
				nameIdx = content
				break nameLoop
			}
		}
		content++
	}

	// There is no valid name inside the []
	if nameIdx < content {
		return nil, 0, 0, ErrNotValidMarkdownStructureName
	}

	content++
	// Next character needs to be ( or otherwise it's not valid
	if source[content] != '(' {
		return nil, 0, 0, ErrNotValidMarkdownStructureContent
	}

	var contentEnd int
contentLoop:
	for content < len(source) {
		switch source[content] {
		case '\n':
			break contentLoop
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				contentEnd = content
				break contentLoop
			}
		}
		content++
	}

	// We have reached the end of the line or the source with an invalid structure
	if contentEnd < nameIdx {
		return nil, 0, 0, ErrNotValidMarkdownStructureContent
	}

	return source[start+1 : nameIdx], nameIdx + 2, contentEnd, nil
}
