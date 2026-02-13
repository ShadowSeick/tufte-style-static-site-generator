package http

type Header uint8

const (
	ContentType Header = iota
	HeaderCount
)

var headersString = [HeaderCount]string{
	ContentType: "Content-Type",
}

func (h Header) String() string {
	if h >= HeaderCount {
		panic("header not available")
	}
	return headersString[h]
}

type ContentTypes uint8

const (
	ApplicationOctetStream ContentTypes = iota
	ContentTypesCount
)

var contentTypeString = [ContentTypesCount]string {
	ApplicationOctetStream: "application/octet-stream",
}

func (ct ContentTypes) String() string {
	if ct >= ContentTypesCount {
		panic("content type not avaliable")
	}
	return contentTypeString[ct]
}
