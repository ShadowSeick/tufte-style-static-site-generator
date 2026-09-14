package text

type tokenType uint8

const (
	headerType tokenType = iota
	subheaderType
	marginNoteType
	sideNoteType
	linkType
	imageType
	inlineCodeType
	blockCodeType
	boldType
	italicType
	textType
	jumpType
	tokenTypeCount
)

var tokenString = [tokenTypeCount]string{
	headerType:     "header",
	subheaderType:  "subheader",
	marginNoteType: "note",
	sideNoteType:   "note",
	linkType:       "link",
	imageType:      "image",
	inlineCodeType: "code",
	blockCodeType:  "block code",
	boldType:       "font_modifier",
	italicType:     "italic",
	jumpType:       "jump",
	textType:       "text",
}

type pos struct {
	line   int
	column int
	offset int
}

type Segment struct {
	start pos
	end   pos
}

func (s *Segment) SetStart(line, column, offset int) {
	s.start = pos{line: line, column: column, offset: offset}
}

func (s *Segment) SetEnd(line, column, offset int) {
	s.end = pos{line: line, column: column, offset: offset}
}

type Token interface {
	SetSegment(s Segment)
	String() string
	Start() int
	End() int
}

type baseToken struct {
	Content Segment
}

func (b *baseToken) SetSegment(s Segment) {
	b.Content = s
}

func (b baseToken) Start() int {
	return b.Content.start.offset
}

func (b baseToken) End() int {
	return b.Content.end.offset
}

type Header struct {
	baseToken
	Level int
}

func (h Header) String() string {
	return tokenString[headerType]
}

type Link struct {
	baseToken
	Name []byte
}

func (l Link) String() string {
	return tokenString[linkType]
}

type Text struct {
	baseToken
}

func (t Text) String() string {
	return tokenString[textType]
}

type Subheader struct {
	baseToken
}

func (sub Subheader) String() string {
	return tokenString[subheaderType]
}

type Jump struct {
	baseToken
}

func (j Jump) String() string {
	return tokenString[jumpType]
}
