package main

import (
	"fmt"
	"unicode"
)

type Lexer struct {
	buf []rune
	loc Location
}

var (
	reserved = []rune{
		'[', ']', ':', '@', '\n', '!', '(', ')', '-',
	}
)

func isReserved(r rune) bool {
	for i := 0; i < len(reserved); i++ {
		if reserved[i] == r {
			return true
		}
	}

	return false
}

func (l *Lexer) peekc() rune {
	if l.loc.index >= len(l.buf) {
		return rune(0)
	}

	return l.buf[l.loc.index]
}

func (l *Lexer) nextc() rune {
	c := l.peekc()
	l.loc.index++
	if c == '\n' {
		l.loc.line++
		l.loc.column = 0
	} else {
		l.loc.column++
	}

	return c
}

func (l *Lexer) next() (Token, error) {
	for unicode.IsSpace(l.peekc()) && l.peekc() != '\n' {
		l.nextc()
	}

	tok := Token{loc: l.loc}
	switch l.peekc() {
	case '[':
		l.nextc()
		tok.value = []rune{'['}
		tok.kind = kindArrayStart
	case ']':
		l.nextc()
		tok.value = []rune{']'}
		tok.kind = kindArrayEnd
	case ':':
		l.nextc()
		tok.value = []rune{':'}
		tok.kind = kindArraySep
	case '@':
		l.nextc()
		tok.value = []rune{'@'}
		tok.kind = kindArgSubst
	case '!':
		l.nextc()
		tok.value = []rune{'!'}
		tok.kind = kindRef
	case '\n':
		l.nextc()
		tok.value = []rune("LF")
		tok.kind = kindLF
	case '(':
		l.nextc()
		tok.value = []rune("(")
		tok.kind = kindParOpen
	case ')':
		l.nextc()
		tok.value = []rune(")")
		tok.kind = kindParClose
	case '-':
		l.nextc()
		if l.peekc() == '>' {
			tok.kind = kindChainOperator
			tok.value = []rune("->")
			l.nextc()
		} else {
			return tok, fmt.Errorf("- expected.")
		}
	case 0:
		tok.kind = kindEOF
	default:
		start := l.loc.index
		if unicode.IsNumber(l.peekc()) {
			for l.peekc() != 0 && unicode.IsNumber(l.peekc()) {
				l.nextc()
			}

			tok.kind = kindNumber
		} else {
			for l.peekc() != 0 && l.peekc() != ' ' && !isReserved(l.peekc()) {
				l.nextc()
			}

			tok.kind = kindIdent
		}

		tok.value = l.buf[start:l.loc.index]
	}

	return tok, nil
}

func (l *Lexer) peek() (Token, error) {
	oldLoc := l.loc
	t, err := l.next()
	l.loc = oldLoc
	return t, err
}
