package text

import (
	"errors"
	"fmt"
	"slices"
)

var (
	ErrHeaderNotValid                   = errors.New("not valid header")
	ErrFontModifiersNotValid            = errors.New("not valid font modifiers")
	ErrBoldNotValid                     = errors.New("not valid bold")
	ErrCodeNotValid                     = errors.New("not valid code")
	ErrInlineCodeNotValid               = errors.New("not valid inline code")
	ErrImageNotValid                    = errors.New("not valid image")
	ErrLinkNotValid                     = errors.New("not valid link")
	ErrNotValidMarkdown                 = errors.New("not valid markdown")
	ErrCustomMarkdownNotValid           = errors.New("not valid custom markdown")
	ErrMarkdownTypeNotImplemented       = errors.New("markdown type not implemented")
	ErrNotValidMarkdownStructureName    = errors.New("not valid markdown balanced structure []() name")
	ErrNotValidMarkdownStructureContent = errors.New("not valid markdown balanced structure []() content")
)

const (
	header        = '#'
	custom        = '?'
	link          = '['
	code          = '`'
	image         = '!'
	fontModifiers = '*'
	lineFeed      = '\n'

	space = ' '

	subheaderIdentifier  = "^subheader"
	marginNoteIdentifier = "^margin-note"
	sideNoteIdentifier   = "^side-note"

	inlineCodeTerminalChars = 1
	blockCodeTerminalChars  = 3
	boldTerminalChars       = 1
	italicTerminalChars     = 2
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

// Parse remade plenty of times. Finaly inlined to be simpler.
// For more information why http://number-none.com/blow/blog/programming/2014/09/26/carmack-on-inlined-code.html
func (p *Parser) Parse(source []byte) ([]Token, []error) {
	var errs []error
	var tokens []Token
	var content int
	var token Token
	var err error
	var childParser Parser
	for p.offset < len(source) {
	out:
		switch source[p.offset] {
		case header:
			if p.column > 1 {
				break
			}

			const maxLevel = 2
			level := 1
			for content < maxLevel {
				if p.peek(source, level) == header {
					level++
				} else {
					break
				}
			}

			if p.peek(source, level) != space {
				err = ErrHeaderNotValid
				content = level // We set the content to level, otherwise it will be an infinite loop
				break
			}

			for range level + 1 { // Drop space
				p.advance()
			}

			token = &Header{
				Level: level,
			}

			for p.peek(source, content) != lineFeed {
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
			for range start - p.offset {
				p.advance()
			}

			content = end - start
			if customToken != nil {
				// We use the childParser to not affect the current one.
				// Surely there is a better way, but for now it will be enough
				childParser.Reset()
				childParser.offset = p.offset
				childParser.column = p.column
				childParser.line = p.line
				children, childErrs := childParser.Parse(source[:p.offset+content-1])
				if len(childErrs) != 0 {
					for _, childErr := range childErrs {
						if err == nil {
							errs = append(errs, childErr)
						} else {
							errs = append(errs, fmt.Errorf("%w: %w", err, childErr))
						}
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

			for range start - p.offset {
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

			for range start - p.offset {
				p.advance()
			}

			content = end - start
		case code:
			terminalChars := 1
			for terminalChars < blockCodeTerminalChars {
				if !isAtEOF(source, p.offset+terminalChars) && p.peek(source, terminalChars) == code {
					terminalChars++
				} else {
					break
				}
			}

			content = terminalChars
			if terminalChars != inlineCodeTerminalChars && terminalChars != blockCodeTerminalChars {
				err = ErrCodeNotValid
				break
			}

			var language []byte
			for i := 0; i < terminalChars; {
				content++

				if isAtEOF(source, p.offset+content) {
					err = ErrCodeNotValid
					break
				}

				switch terminalChars {
				case inlineCodeTerminalChars:
					if p.peek(source, content) == '\n' {
						err = ErrInlineCodeNotValid
						break out
					}
				case blockCodeTerminalChars:
					if len(language) == 0 && p.peek(source, content) == '\n' {
						language = source[p.offset+terminalChars : content]
					}
				}

				if p.peek(source, content) == code {
					i++
				}
			}

			for range terminalChars + len(language) {
				p.advance()
			}

			switch terminalChars {
			case inlineCodeTerminalChars:
				token = &InlineCode{}
			case blockCodeTerminalChars:
				token = &BlockCode{
					Language: language,
				}
			}

			content = content - len(language)
		case fontModifiers:
			terminalChars := 1
			for terminalChars < italicTerminalChars {
				if !isAtEOF(source, p.offset+terminalChars) && p.peek(source, terminalChars) == fontModifiers {
					terminalChars++
				} else {
					break
				}
			}

			content = terminalChars
			for i := 0; i < terminalChars; {
				content++

				if isAtEOF(source, p.offset+content) {
					err = ErrFontModifiersNotValid
					break out
				}

				if p.peek(source, content) == '\n' {
					err = ErrFontModifiersNotValid
					break out
				}

				if p.peek(source, content) == fontModifiers {
					i++
				}
			}

			switch terminalChars {
			case boldTerminalChars:
				token = &Bold{}
			case italicTerminalChars:
				token = &Italic{}
				content--
			}

			for range terminalChars {
				p.advance()
			}
		case lineFeed:
			token = &LineFeed{}
			content = 1
			p.line++
			p.column = 0
		}
		// This is used as a goto statement.
		// It simplifies the logic and makes the program not need to check isAtEOF at first of each markdown type
		if token == nil {
			token = &Text{}
			for !isAtEOF(source, p.offset+content) && !slices.Contains([]byte{header, custom, link, code, image, fontModifiers, lineFeed}, p.peek(source, content)) {
				content++
			}
		}

		if err != nil {
			errs = append(errs, err)
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

	content := start
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

	return source[start+1 : nameIdx], nameIdx + 2, contentEnd + 1, nil
}
