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
	marginNoteType: "margin note",
	sideNoteType:   "side note",
	linkType:       "link",
	imageType:      "image",
	inlineCodeType: "code",
	blockCodeType:  "block code",
	boldType:       "bold",
	italicType:     "italic",
	jumpType:       "jump",
	textType:       "text",
}

type pos struct {
	line   int
	column int
	offset int
}

type segment struct {
	start pos
	end   pos
}

type Token interface {
	SetContent(start pos, end pos)
	String() string
	Start() int
	End() int
}

type CustomToken interface {
	SetChildren(ts []Token)
	GetChildren() []Token
	Token
}

type baseToken struct {
	content segment
}

func (b *baseToken) SetContent(start pos, end pos) {
	b.content = segment{
		start: start,
		end:   end,
	}
}

func (b baseToken) Start() int {
	return b.content.start.offset
}

func (b baseToken) End() int {
	return b.content.end.offset
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

type Image struct {
	baseToken
	Caption []byte
}

func (i Image) String() string {
	return tokenString[imageType]
}

type Bold struct {
	baseToken
}

func (b Bold) String() string {
	return tokenString[boldType]
}

type Italic struct {
	baseToken
}

func (i Italic) String() string {
	return tokenString[italicType]
}

type Jump struct {
	baseToken
}

func (j Jump) String() string {
	return tokenString[jumpType]
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

type MarginNote struct {
	baseToken
	children []Token
}

func (mn MarginNote) String() string {
	return tokenString[marginNoteType]
}

func (mn *MarginNote) SetChildren(children []Token) {
	mn.children = append(mn.children, children...)
}

func (mn MarginNote) GetChildren() []Token {
	return mn.children
}

type SideNote struct {
	baseToken
	children []Token
}

func (sn SideNote) String() string {
	return tokenString[sideNoteType]
}

func (sn *SideNote) SetChildren(children []Token) {
	sn.children = append(sn.children, children...)
}

func (sn SideNote) GetChildren() []Token {
	return sn.children
}
