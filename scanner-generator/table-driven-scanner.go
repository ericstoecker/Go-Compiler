package scannergenerator

import (
	"compiler/token"
	"slices"
)

type TableDrivenScanner struct {
	input        string
	position     int
	readPosition int
	ch           byte

	dfa *Dfa
}

func New(input string, dfa *Dfa) *TableDrivenScanner {
	s := &TableDrivenScanner{input: input, dfa: dfa}
	s.readChar()
	return s
}

func (s *TableDrivenScanner) NextToken() token.Token {
	s.skipWhitespace()

	if s.position >= len(s.input) {
		return token.Token{Type: token.EOF, Literal: ""}
	}

	state := s.dfa.InitialState
	lexeme := ""
	stack := []int{}

	for state != -1 {
		stack = append(stack, state)
		lexeme += string(s.ch)

		if stateAfterTransition, ok := s.dfa.Transitions[string(s.ch)][state]; ok {
			state = stateAfterTransition
			s.readChar()
		} else {
			state = -1
		}
	}

	for !s.isAcceptingState(state) && len(stack) > 1 {
		state = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		lexeme = lexeme[:len(lexeme)-1]
		s.rollback()
	}
	s.readChar()

	if s.isAcceptingState(state) {
		tokenType, ok := s.dfa.TypeTable[state]
		if !ok {
			panic("In an accepting state, but no token type found")
		}
		return token.Token{Type: tokenType, Literal: lexeme}
	}
	return token.Token{Type: token.ILLEGAL, Literal: lexeme}
}

func (s *TableDrivenScanner) isAcceptingState(state int) bool {
	return slices.Contains(s.dfa.AcceptingStates, state)
}

func (s *TableDrivenScanner) readChar() {
	if s.readPosition >= len(s.input) {
		s.ch = 0
	} else {
		s.ch = s.input[s.readPosition]
	}

	s.position = s.readPosition
	s.readPosition += 1
}

func (s *TableDrivenScanner) rollback() {
	if s.readPosition <= 1 {
		return
	}

	s.readPosition = s.position
	s.position -= 1
	s.ch = s.input[s.position]
}

func (s *TableDrivenScanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\r' || s.ch == '\n' || s.ch == '\t' {
		s.readChar()
	}
}
