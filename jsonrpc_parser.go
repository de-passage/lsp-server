package main


import (
	"bufio"
	"bytes"
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

func (p *Parser) Parse(r io.Reader) (string, error) {
	reader := bufio.NewReader(r)

	for {
		switch p.state {
		case ParsingHeader:
			lineBytes, _, err := reader.ReadLine()
			if err != nil {
				return "", err
			}

			line := string(lineBytes)
			p.headerBuf.WriteString(line)
			p.headerBuf.WriteString("\r\n")

			if line == "" {
				contentLenHeader := "Content-Length: "
				if strings.HasPrefix(p.headerBuf.String(), contentLenHeader) {
					p.contentLen, _ = strconv.Atoi(strings.TrimPrefix(strings.Split(p.headerBuf.String(), "\r\n")[0], contentLenHeader))
				} else {
					return "", fmt.Errorf("invalid header: %s", p.headerBuf.String())
				}

				p.state = ParsingBody
				p.headerBuf.Reset()
			}

		case ParsingBody:
			if p.contentLen <= 0 {
				return "", fmt.Errorf("invalid content length")
			}

			bodyBytes := make([]byte, p.contentLen)
			_, err := io.ReadFull(reader, bodyBytes)
			if err != nil {
				return "", err
			}

			p.bodyBuf.Write(bodyBytes)
			if p.bodyBuf.Len() == p.contentLen {
				p.state = ParsingHeader
				message := p.bodyBuf.String()
				p.bodyBuf.Reset()
				p.contentLen = 0
				return message, nil
			}
		}
	}
}

