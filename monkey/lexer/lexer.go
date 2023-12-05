package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads sets the current ch value to the character at the index readPosition, then increments readPosition
// currently only supports ascii
// TODO: support Unicode and emojis!
func (l *Lexer) readChar() {
	// if readPosition at end of input, sets ch to ASCII NUL ("0"),
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		// set ch to the character at current readPosition
		l.ch = l.input[l.readPosition]
	}
	// increment readPosition
	l.position = l.readPosition
	l.readPosition += 1
}

// NextToken looks at the character currently at l.ch,  returns a corresponding token, and then increments to the next ch.
func (l *Lexer) NextToken() token.Token {
	var _token token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		_token = newToken(token.ASSIGN, l.ch)
	case ';':
		_token = newToken(token.SEMICOLON, l.ch)
	case '(':
		_token = newToken(token.LPAREN, l.ch)
	case ')':
		_token = newToken(token.RPAREN, l.ch)
	case ',':
		_token = newToken(token.COMMA, l.ch)
	case '+':
		_token = newToken(token.PLUS, l.ch)
	case '{':
		_token = newToken(token.LBRACE, l.ch)
	case '}':
		_token = newToken(token.RBRACE, l.ch)
	case 0:
		_token.Literal = ""
		_token.Type = token.EOF
	default:
		// in order to parse whole identifiers like "let" or words,
		// if we hit a char we use readIdentifier() to read a whole identifier and return it as a token
		// otherwise we don't know what the character is, and return illegal
		if isLetter(l.ch) {
			_token.Literal = l.readIdentifier()
			_token.Type = token.LookupIdent(_token.Literal)
			return _token
		} else if isDigit(l.ch) {
			_token.Type = token.INT
			_token.Literal = l.readNumber()
			return _token
		} else {
			_token = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return _token
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier reads an identifier and continues through input as long as it finds a letter
// when it reaches a non-letter character, it returns the slice of input corresponding all the letters in a row
// (a whole identifier).
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) { // 'for' loops as long as it reads true, and can be used as a while loop?!
		l.readChar()
	}
	return l.input[position:l.position]
}

// isLetter accepts a char as byte and determines if its value matches a-zA-Z_
// underscore allows snake case variable names
// TODO: match unicode (with Go builtin?)
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}
