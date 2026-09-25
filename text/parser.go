package text

import (
	"errors"
	"fmt"
	"slices"
)

var (
	ErrTitleNotValid              = errors.New("not valid title")
	ErrCustomMarkdownNotValid     = errors.New("not valid custom markdown")
	ErrMarkdownTypeNotImplemented = errors.New("markdown type not implemented")
	ErrNotValidLink               = errors.New("not valid link")
)

// TODO: I need to implement image, code and font modifiers

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

func (p *Parser) Parse() ([]Token, []error) {
	var errs []error
	var tokens []Token
	for p.offset < len(p.source) {
		// We are at the start of the line;
		// Test for header or subheader
		if p.column == 0 {
			switch p.source[p.offset] {
			case title:
				t, err := p.parseTitle()
				if t != nil {
					tokens = append(tokens, t)
					continue
				}
				if err != nil {
					errs = append(errs, err)
				}

			case custom:
				t, err := p.parseCustom()
				if t != nil {
					tokens = append(tokens, t)
					continue
				}
				if err != nil {
					errs = append(errs, err)
				}
			}
		}

		switch p.source[p.offset] {
		case custom:
			t, err := p.parseCustom()
			if t != nil {
				tokens = append(tokens, t)
				continue
			}
			if err != nil {
				errs = append(errs, err)
			}
		case link:
			t, err := p.parseLink()
			if t != nil {
				tokens = append(tokens, t)
				continue
			}
			if err != nil {
				errs = append(errs, err)
			}
		case image:
			// Parse maybeImage
		case code:
			// Parse code
		case fontModifiers:
			// Parse font modifiers
		case breakLine:
			tokens = append(tokens, &Jump{})
			p.advance()
			p.line++
			p.column = 0
			continue
		default:
			t, err := p.parseText()
			if t != nil {
				tokens = append(tokens, t)
			}
			if err != nil {
				errs = append(errs, err)
			}
		}

		if !p.isAtEOF() {
			p.advance()
		}
	}
	return tokens, errs
}

func (p *Parser) advance() {
	p.offset++
	p.column++
}

func (p *Parser) peek(n int) byte {
	return p.source[p.offset+n]
}

func (p *Parser) parseTitle() (Token, error) {
	const maxLevel = 2
	var level int
	var start int
	for start < maxLevel {
		start++
		b := p.peek(start)
		if b == space {
			break
		} else if b == title {
			level++
		} else {
			return nil, ErrTitleNotValid
		}
	}

	for range start + 1 {
		p.advance()
	}

	token := &Header{
		Level: level,
	}
	p.setSegment(token, breakLine)

	return token, nil
}

func (p *Parser) parseCustom() (Token, error) {
	identifier, contentStart, contentEnd, err := getContentFromBalanced(p.offset, p.source)
	if err != nil {
		return p.parseText()
	}

	for range contentStart {
		p.advance()
	}

	var token Token
	id := string(identifier)
	switch id {
	case subheaderIdentifier:
		token = &Subheader{}
	case marginNoteIdentifier:
		token, err = p.parseChildren(&MarginNote{}, contentStart, contentEnd)
	case sideNoteIdentifier:
		token, err = p.parseChildren(&SideNote{}, contentStart, contentEnd)
	}

	p.setBalancedTokenSegment(token, contentStart, contentEnd)
	return token, err
}

func (p *Parser) parseChildren(token CustomToken, start, end int) (Token, error) {
	var errs []error
	var tokens []Token
	for start < end {
		switch p.source[start] {
		case link:
			t, err := p.parseLink()
			if t != nil {
				tokens = append(tokens, t)
				continue
			}
			if err != nil {
				errs = append(errs, err)
			}
		case image:
			// Parse maybeImage
		case code:
			// Parse code
		case fontModifiers:
			// Parse font modifiers
		default:
			t, err := p.parseText()
			if t != nil {
				tokens = append(tokens, t)
			}
			if err != nil {
				errs = append(errs, err)
			}
		}

		if !p.isAtEOF() {
			p.advance()
		}
	}
	var err error
	for _, e := range errs {
		err = fmt.Errorf("%w: %w", err, e)
	}

	token.SetChildren(tokens)
	return token, err
}

func (p *Parser) parseLink() (Token, error) {
	name, contentStart, contentEnd, err := getContentFromBalanced(p.offset, p.source)
	if err != nil {
		return p.parseText()
	}

	link := &Link{
		Name: name,
	}

	for range contentStart {
		p.advance()
	}

	p.setBalancedTokenSegment(link, contentStart, contentEnd)
	return link, nil
}

var (
	ErrNotValidMarkdownStructureName    = errors.New("not valid markdown balanced structure []() name")
	ErrNotValidMarkdownStructureContent = errors.New("not valid markdown balanced structure []() content")
)

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
	var nameIdx int
nameLoop:
	for content < len(source) {
		switch source[content] {
		case '\n':
			break nameLoop
		case ']':
			nameIdx = content
			break nameLoop
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

	var depth int
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

func (p *Parser) parseImage() (Token, error) {
	return nil, nil
}

func (p *Parser) parseCode() (Token, error) {
	return nil, nil
}

func (p *Parser) parseFontModifiers() (Token, error) {
	return nil, nil
}

func (p *Parser) parseText() (Token, error) {
	token := &Text{}
	var segment Segment
	segment.SetStart(p.line, p.column, p.offset)
	for !slices.Contains([]byte{title, custom, link, code, image, fontModifiers, breakLine}, p.source[p.offset]) {
		p.advance()
		if p.isAtEOF() {
			break
		}
	}

	segment.SetEnd(p.line, p.column, p.offset)
	token.SetSegment(segment)
	return token, nil
}

func (p *Parser) setSegment(token Token, terminalChar byte) {
	var segment Segment
	segment.SetStart(p.line, p.column, p.offset)
	for p.source[p.offset] != terminalChar {
		p.advance()
		if p.isAtEOF() {
			break
		}
	}

	segment.SetEnd(p.line, p.column, p.offset)
	token.SetSegment(segment)
}

func (p *Parser) isAtEOF() bool {
	return p.offset >= len(p.source)
}

func (p *Parser) setBalancedTokenSegment(token Token, start, end int) {
	var segment Segment
	// I don't know in which column is it; for now I will set it as the same of offset
	// TODO: This needs to change to show proper errors
	segment.SetStart(p.line, start, start)
	segment.SetEnd(p.line, end, end)
	// Drop ')'
	for range end - start + 1 {
		p.advance()
	}

	token.SetSegment(segment)
}
