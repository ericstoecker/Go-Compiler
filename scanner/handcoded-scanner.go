package scanner

import "compiler/token"

type HandcodedScanner struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewHandcodedScanner(input string) *HandcodedScanner {
	l := &HandcodedScanner{input: input}
	l.readChar()
	return l
}

func (s *HandcodedScanner) readChar() {
	if s.readPosition >= len(s.input) {
		s.ch = 0
	} else {
		s.ch = s.input[s.readPosition]
	}
	s.position = s.readPosition
	s.readPosition += 1
}

func (s *HandcodedScanner) NextToken() token.Token {
	var tok token.Token

	s.skipWhitespace()

	switch s.ch {
	case '=':
		tok = s.readTwoCharToken(tok, '=', token.EQUALS, token.ASSIGN)
	case '!':
		tok = s.readTwoCharToken(tok, '=', token.NOT_EQUALS, token.BANG)
	case ';':
		tok = newToken(token.SEMICOLON, s.ch)
	case ',':
		tok = newToken(token.COMMA, s.ch)
	case '(':
		tok = newToken(token.LPAREN, s.ch)
	case ')':
		tok = newToken(token.RPAREN, s.ch)
	case '{':
		tok = newToken(token.LBRACE, s.ch)
	case '}':
		tok = newToken(token.RBRACE, s.ch)
	case '[':
		tok = newToken(token.LBRACKET, s.ch)
	case ']':
		tok = newToken(token.RBRACKET, s.ch)
	case '+':
		tok = newToken(token.PLUS, s.ch)
	case '-':
		tok = newToken(token.MINUS, s.ch)
	case '/':
		tok = newToken(token.SLASH, s.ch)
	case '*':
		tok = newToken(token.ASTERIK, s.ch)
	case ':':
		tok = newToken(token.COLON, s.ch)
	case '<':
		tok = s.readTwoCharToken(tok, '=', token.LESS_EQUAL, token.LT)
	case '>':
		tok = s.readTwoCharToken(tok, '=', token.GREATER_EQUAL, token.GT)
	case '&':
		tok = s.readTwoCharToken(tok, '&', token.AND, token.ILLEGAL)
	case '|':
		tok = s.readTwoCharToken(tok, '|', token.OR, token.ILLEGAL)
	case '"':
		tok.Literal = s.readString()
		tok.Type = token.STRING
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(s.ch) {
			tok.Literal = s.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(s.ch) {
			tok.Literal = s.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, s.ch)
		}
	}
	s.readChar()
	return tok
}

func (s *HandcodedScanner) peek() byte {
	if s.readPosition >= len(s.input) {
		return 0
	}
	return s.input[s.readPosition]
}

func (s *HandcodedScanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\r' || s.ch == '\n' || s.ch == '\t' {
		s.readChar()
	}
}

func (s *HandcodedScanner) readIdentifier() string {
	position := s.position
	for isLetter(s.ch) {
		s.readChar()
	}
	return s.input[position:s.position]
}

func (s *HandcodedScanner) readString() string {
	s.readChar()
	position := s.position
	for s.ch != '"' && s.ch != 0 {
		s.readChar()
	}
	result := s.input[position:s.position]
	s.readChar()
	return result
}

func (s *HandcodedScanner) readNumber() string {
	position := s.position
	for isDigit(s.ch) {
		s.readChar()
	}
	return s.input[position:s.position]
}

func (s *HandcodedScanner) readTwoCharToken(tok token.Token, expectedNext byte, nextType token.TokenType, alternative token.TokenType) token.Token {
	nextChar := s.peek()
	if nextChar == expectedNext {
		tok.Type = nextType
		ch := s.ch
		s.readChar()
		tok.Literal = string(ch) + string(s.ch)
	} else {
		tok = newToken(alternative, s.ch)
	}
	return tok
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
