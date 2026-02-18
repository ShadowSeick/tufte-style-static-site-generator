package http

type Status uint8

const (
	StatusUnauthorized Status = iota
	StatusNotFound
	StatusCount
)

var statusCode = [StatusCount]int{
	StatusUnauthorized: 401,
	StatusNotFound: 404,
}

func (s Status) Code() int {
	if s >= StatusCount {
		panic("status not handled")
	}
	return statusCode[s]
}
