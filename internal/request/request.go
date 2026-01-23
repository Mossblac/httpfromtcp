package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const buffersize = 8

const (
	Initialized int = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	ParserState int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	if r.ParserState == Initialized {
		requestL, BytesConsumed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if BytesConsumed == 0 {
			return 0, nil
		}
		if BytesConsumed > 0 {
			r.RequestLine = *requestL
			r.ParserState = Done
			return BytesConsumed, nil
		}
	}

	if r.ParserState == Done {
		return 0, fmt.Errorf("error: trying to read data in done state")
	}

	return 0, fmt.Errorf("error: unknown state")
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, buffersize, buffersize)

	readToIndex := 0

	R := Request{
		ParserState: Initialized,
	}

	for R.ParserState != Done {
		if readToIndex >= len(buf) {
			newBuffer := make([]byte, len(buf)*2)
			copy(newBuffer, buf)
			buf = newBuffer
		}
		NumBytes, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			if err == io.EOF && NumBytes == 0 {
				break
			}
		}

		readToIndex += NumBytes
		IntParsed, err := R.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[IntParsed:readToIndex])
		readToIndex -= IntParsed

	}
	if R.ParserState != Done {
		return nil, fmt.Errorf("incomplete request: unexpected EOF")
	}

	return &R, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}
