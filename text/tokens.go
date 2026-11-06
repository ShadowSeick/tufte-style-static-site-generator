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

type headerLevel int

const (
	_ headerLevel = iota
	header1
	header2
	header3
	headerCount
)

type Header struct {
	baseToken
	Level headerLevel
}

func (h Header) String() string {
	return tokenString[headerType]
}

type Link struct {
	baseToken
	Name []byte
}

func (Link) String() string {
	return tokenString[linkType]
}

func (l Link) End() int {
	return l.baseToken.End() - 1
}

type Image struct {
	baseToken
	Caption []byte
}

func (Image) String() string {
	return tokenString[imageType]
}

func (i Image) End() int {
	return i.baseToken.End() - 1
}

type Bold struct {
	baseToken
}

func (Bold) String() string {
	return tokenString[boldType]
}

func (b Bold) End() int {
	return b.baseToken.End() - 1
}

type Italic struct {
	baseToken
}

func (Italic) String() string {
	return tokenString[italicType]
}

func (i Italic) End() int {
	return i.baseToken.End() - 2
}

type LineFeed struct {
	baseToken
}

func (LineFeed) String() string {
	return tokenString[jumpType]
}

type Text struct {
	baseToken
}

func (Text) String() string {
	return tokenString[textType]
}

type Subheader struct {
	baseToken
}

func (sub Subheader) String() string {
	return tokenString[subheaderType]
}

func (sub Subheader) End() int {
	return sub.baseToken.End() - 1
}

type MarginNote struct {
	baseToken
	children []Token
}

func (MarginNote) String() string {
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

func (SideNote) String() string {
	return tokenString[sideNoteType]
}

func (sn *SideNote) SetChildren(children []Token) {
	sn.children = append(sn.children, children...)
}

func (sn SideNote) GetChildren() []Token {
	return sn.children
}

type InlineCode struct {
	baseToken
}

func (_ InlineCode) String() string {
	return tokenString[inlineCodeType]
}

type BlockCode struct {
	baseToken
	Language []byte
}

func (_ BlockCode) String() string {
	return tokenString[blockCodeType]
}
