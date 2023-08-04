package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type ParserState int

const (
	ParsingHeader ParserState = iota
	ParsingBody
)

type Parser struct {
	state       ParserState
	contentLen  int
	headerBuf   bytes.Buffer
	bodyBuf     bytes.Buffer
}

func NewParser() *Parser {
	return &Parser{
		state: ParsingHeader,
	}
}

func (p *Parser) Parse(r io.Reader) (Request[json.RawMessage], error) {
	reader := bufio.NewReader(r)

	for {
		switch p.state {
		case ParsingHeader:
			lineBytes, _, err := reader.ReadLine()
			if err != nil {
				return Request[json.RawMessage]{}, err
			}

			line := string(lineBytes)
			p.headerBuf.WriteString(line)
			p.headerBuf.WriteString("\r\n")

			if line == "" {
				contentLenHeader := "Content-Length: "
				if strings.HasPrefix(p.headerBuf.String(), contentLenHeader) {
					p.contentLen, _ = strconv.Atoi(strings.TrimPrefix(strings.Split(p.headerBuf.String(), "\r\n")[0], contentLenHeader))
				} else {
					return Request[json.RawMessage]{}, fmt.Errorf("invalid header: %s", p.headerBuf.String())
				}

				p.state = ParsingBody
				p.headerBuf.Reset()
			}

		case ParsingBody:
			if p.contentLen <= 0 {
				return Request[json.RawMessage]{}, fmt.Errorf("invalid content length")
			}

			bodyBytes := make([]byte, p.contentLen)
			_, err := io.ReadFull(reader, bodyBytes)
			if err != nil {
				return Request[json.RawMessage]{}, err
			}

			p.bodyBuf.Write(bodyBytes)
			if p.bodyBuf.Len() == p.contentLen {
				p.state = ParsingHeader

        var message Request[json.RawMessage]
				err := json.NewDecoder(&p.bodyBuf).Decode(&message);

				p.bodyBuf.Reset()
				p.contentLen = 0
				return message, err
			}
		}
	}
}

func ParseLSPRequest(request Request[json.RawMessage]) (LSPParams, error) {
  switch request.Method {
  case "initialize":
    var params InitializeParams
    err := json.Unmarshal(request.Params, &params)
    return params, err
  case "textDocument/didOpen":
    var params DidOpenTextDocumentParams
    err := json.Unmarshal(request.Params, &params)
    return params, err
  }
  return nil, fmt.Errorf("unknown method: %s", request.Method)
}
