package main

import (
	"log"
	"net"
)

type Command struct {
	args []string
	conn net.Conn
}

func (p *Parser) command() (Command, error) {
	b, err := p.r.ReadByte()
	if err != nil {
		return Command{}, err
	}
	if b == '*' {
		log.Println("resp Array")
		return p.respArray()
	} else {
		line, err := p.readLine()
		if err != nil {
			return Command{}, err
		}
		p.pos = 0
		p.line = append([]byte{}, b)
		p.line = append(p.line, line...)
		return p.inline()

	}
}
func (p *Parser) inline() (Command, error) {
	for p.current() == ' ' {
		p.advance()
	}
	cmd := Command{conn: p.conn}
	for !p.atEnd() {
		arg, err := p.consumeArg()
		if err != nil {
			return cmd, err
		}
		if arg != "" {
			cmd.args = append(cmd.args, arg)
		}
	}
	return cmd, nil
}

func (p *Parser) consumeArg() (s string, err error) {
	for p.current() == ' ' {
		p.advance()
	}
	if p.current() == '"' {
		p.advance()
		buf, err := p.consumeString()
		return string(buf), err
	}
	for !p.atEnd() && p.current() != ' ' && p.current() != '\r' {
		s += string(p.current())
		p.advance()
	}
	return
}
