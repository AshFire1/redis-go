package main

import (
	"bufio"
	"errors"
	"net"
)

type Parser struct {
	conn net.Conn
	r    *bufio.Reader
	line []byte
	pos  int
}

func NewParser(conn net.Conn) *Parser {
	return &Parser{
		conn: conn,
		r:    bufio.NewReader(conn),
		line: make([]byte, 0),
		pos:  0,
	}
}
func (p *Parser) current() byte {
	if p.atEnd() {
		return '\r'
	}
	return p.line[p.pos]
}
func (p *Parser) advance() {
	p.pos++
}
func (p *Parser) atEnd() bool {
	return p.pos >= len(p.line)
}
func (p *Parser) readLine() ([]byte, error) {
	line, err := p.r.ReadBytes('\r')
	if err != nil {
		return nil, err
	}
	if _, err := p.r.ReadByte(); err != nil {
		return nil, err
	}
	return line[:len(line)-1], nil
}
func (p *Parser) consumeString() (s []byte, err error) {
	for p.current() != '"' && !p.atEnd() {
		cur := p.current()
		p.advance()
		next := p.current()
		if cur == '\\' && next == '"' {
			s = append(s, '"')
			p.advance()
		} else {
			s = append(s, cur)
		}
	}
	if p.current() != '"' {
		return nil, errors.New("Unbalanced quotes in Request")
	}
	p.advance()
	return
}
