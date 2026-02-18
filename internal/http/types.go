package http

type Request struct {
	URL string
	Headers map[string]string
	Body []byte
}
