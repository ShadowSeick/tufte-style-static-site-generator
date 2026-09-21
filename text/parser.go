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

	leftParenthesis         = '('
	rightParenthesis        = ')'
	closeSquaredParenthesis = ']'
	space                   = ' '

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
	identifier, contentStart, contentEnd, isBalanced := getContentFromBalanced(p.offset, p.source)
	if !isBalanced {
		return p.parseText()
	}

	for range contentStart {
		p.advance()
	}

	var err error
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
	name, contentStart, contentEnd, isNotLink := getContentFromBalanced(p.offset, p.source)

	if isNotLink {
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
	openParenthesis  = []byte{'[', '('}
	closeParenthesis = []byte{']', ')'}
	contraryChars    = map[byte]byte{
		']': '[',
		')': '(',
	}
)

// What if I make 2 slices
// 1 -> []int -> maintaining the idx
// 2 -> []byte -> maintaining the byte
// I will need to identify afterwards the depth, but this is not the problem.
// When do I stop?
//
// I am thinking this wrongly. I just need to split it in two!
// First [] then (), this way it's much easier to make

func getContentFromBalanced(start int, source []byte) ([]byte, int, int, bool) {
	if source[start] != '[' {
		return nil, 0, 0, false
	}

	content := start + 1
	var curr byte
	balanced := []byte{source[start]}
	depth := 1
	var nameIdx int
	var isNotValid bool
	for {
		if content >= len(source) {
			isNotValid = true
			break
		}

		curr = source[content]
		// Balanced markdown structures must live in the same line
		if curr == '\n' {
			isNotValid = true
			break
		}

		if slices.Contains(openParenthesis, curr) {
			depth++
			balanced = append(balanced, curr)
		}

		if depth == 0 && slices.Contains(closeParenthesis, curr) {
			isNotValid = true
			break
		}

		var prev byte
		if len(balanced) > 0 {
			prev = balanced[len(balanced)-1]
		}

		contrary, ok := contraryChars[curr]
		if ok && prev != contrary {
			isNotValid = true
			break
		}

		switch curr {
		case ']':
			depth--
			balanced = balanced[0 : len(balanced)-1]
			if depth == 0 {
				nameIdx = content
			}
			break
		case ')':
			depth--
			balanced = balanced[0 : len(balanced)-1]
		}

		if len(balanced) == 0 && curr == ')' {
			break
		}
		content++
	}

	var res []byte
	var contentStart int
	var contentEnd int
	if !isNotValid {
		res = source[start+1 : nameIdx]
		// Discard ']('
		contentStart = nameIdx + 2
		contentEnd = content
	}
	return res, contentStart, contentEnd, !isNotValid
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
