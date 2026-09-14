package text

import (
	"errors"
	"slices"
)

var (
	ErrTitleNotValid              = errors.New("not valid title")
	ErrCustomMarkdownNotValid     = errors.New("not valid custom markdown")
	ErrMarkdownTypeNotImplemented = errors.New("markdown type not implemented")
	ErrNotValidLink               = errors.New("not valid link")
)

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

	subheaderIdentifier  = "[^subheader]"
	marginNoteIdentifier = "[^margin-note]"
	sideNoteIdentifier   = "[^side-note]"
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
		// Test for header of subheader
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
				t, err := p.parseCustom(subheaderIdentifier)
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
			// Parse maybeCustom
			// Here what I need to do is to parse it recursively
			// it can contain as many tokens possible
			// We should just parse the contents in it as well.
			//
			// The main idea is to parse only the content inside it,
			// treating it as another source, returning the actualy tokens here
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
	// We would get the identifier directly and depending on which one is it parse it one way or another
	_, contentStart, isBalanced := getContentFromBalanced(p.offset, p.source)
	if !isBalanced {
		return p.parseText()
	}

	for range contentStart {
		p.advance()
	}

	token := &Subheader{}
	p.setSegment(token, rightParenthesis)
	return token, nil
}

func (p *Parser) parseLink() (Token, error) {
	name, contentStart, isNotLink := getContentFromBalanced(p.offset, p.source)

	if isNotLink {
		return p.parseText()
	}

	link := &Link{
		Name: name,
	}

	for range contentStart {
		p.advance()
	}

	p.setSegment(link, rightParenthesis)
	return link, nil
}

var openParenthesis = []byte{'[', '('}

func getContentFromBalanced(start int, source []byte) ([]byte, int, bool) {
	if source[start] != '[' {
		return nil, 0, false
	}

	content := start + 1
	var next byte
	balanced := []byte{source[start]}
	var nameIdx int
	var isNotValid bool
outer:
	for {
		next = source[content]
		if next == breakLine {
			break
		}

		if slices.Contains(openParenthesis, next) {
			balanced = append(balanced, next)
		}
		prev := balanced[len(balanced)-1]
		switch next {
		case ']':
			if prev != '[' {
				isNotValid = true
				break outer
			}
			balanced = balanced[0 : len(balanced)-1]
			nameIdx = content
		case ')':
			if prev != '(' {
				isNotValid = true
				break outer
			}

			balanced = balanced[0 : len(balanced)-1]
		case '}':
			isNotValid = true
			break outer
		}

		if len(balanced) == 0 && next == ')' {
			break
		}
		content++
	}

	var res []byte
	var contentStart int
	if !isNotValid {
		res = source[start+1 : nameIdx]
		// Discard ']('
		contentStart = nameIdx + 2
	}
	return res, contentStart, !isNotValid
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
