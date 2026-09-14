package parse

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/ShadowSeick/tufte-style-static-site-generator/ast"
)

var ErrTitleNotValid = errors.New("not valid title")

const (
	title         = '#'
	custom        = '?'
	link          = '['
	code          = '`'
	image         = '!'
	fontModifiers = '*'

	subheaderIdentifier  = "[^subheader"
	marginNoteIdentifier = "[^margin-note"
	sideNoteIdentifier   = "[^side-note"
)

const defaultBufferSize = 1024

type Parser struct {
	rd *bufio.Reader
	wt *bufio.Writer
}

func (p *Parser) ReadFromFile(rd io.Reader) {
	if rd == nil {
		p.rd = bufio.NewReaderSize(rd, defaultBufferSize)
	} else {
		p.rd.Reset(rd)
	}
}

func (p *Parser) WriteToFile(wt io.Writer) {
	if wt == nil {
		p.wt = bufio.NewWriterSize(wt, defaultBufferSize)
	} else {
		p.wt.Reset(wt)
	}
}

func (p *Parser) ConvertToHTML() error {
	for {
		if curr.Complete() {
			fileAST = append(fileAST, curr)
		}

		bytes, err := p.parseLine()
		if err != nil && !errors.Is(io.EOF, err) {
			return err
		}

		// Process it to convert to AST
		// Keep track of the last one (len(fileAST -1))
		// We need to see if it has finished or not; if it has finished
		// create a new one, if it hasn't continue processing the current one
		// To know if it has finished we just check the last byte of the content is NUL (0 byte character)

		// Error is EOF
		if err != nil {
			break
		}
	}

	var content strings.Builder
	for _, as := range fileAST {
		content.WriteString(as.HTML())
	}

	_, err := p.wt.WriteString(content.String())
	if err != nil {
		return err
	}

	return nil
}

const (
	space     = ' '
	breakLine = '\n'
)

func (p *Parser) properParseLine() ([]ast.AST, error) {
	var fileAST []ast.AST
	var curr ast.AST
	var parseErr error
	for {
		b, err := p.rd.ReadByte()
		if err != nil {
			return fileAST, err
		}

		var bytes []byte

		switch b {
		case title:
			header, err := parseTitle()
			if errors.Is(io.EOF, err) {
				return fileAST, err
			}
			if errors.Is(ErrInvalidTitle, err) {
				continue
			}
			if err != nil {
				return nil, err
			}
			fileAST = append(fileAST, header)
			continue
		case custom:

		}

	}
}

var ErrInvalidTitle = errors.New("invalid title")

func (p *Parser) parseTitle() (ast.AST, error) {
	const maxLevel = 2
	i := 1
	for i <= maxLevel {
		b, err := p.rd.ReadByte()
		if err != nil {
			return nil, err
		}

		if b == space {
			break
		} else if b == title {
			i++
		} else {
			return nil, ErrInvalidTitle
		}
	}

	content, err := p.rd.ReadSlice(breakLine)
	if err != nil && !errors.Is(io.EOF, err) {
		return nil, err
	}

	return ast.Header{
		Level:   i,
		Content: content,
	}, err
}

func (p *Parser) parseSubtitle() (ast.AST, error) {
	var 
	for {
		b, err := p.rd.ReadByte()
		if err != nil {
			return nil, err
		}


	}
}

func (p *Parser) parseLine() ([]ast.AST, error) {
	slice, err := p.rd.ReadSlice(breakLine)
	if err != nil && !errors.Is(io.EOF, err) {
		return nil, err
	}

	if (len(slice)) == 0 {
		return nil, nil
	}

	// check first character of line
	switch slice[0] {
	case title:
		return parseTitle(slice), nil
	case custom:
		// parse subtitle
		return nil, nil
	default:
		// Do nothing
	}

	for i := 0; i < len(slice); i++ {
		switch slice[i] {
		case custom:
			// parse notes
		case link:
			// parse link
		case code:
			// parse code
		case fontModifiers:
			// parse font modifiers
		}
	}
	return slice, err
}

const supportedHeader = 3

func parseSubtitle(data []byte) []byte {
	return nil
}
