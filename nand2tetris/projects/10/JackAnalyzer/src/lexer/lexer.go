package lexer

import (
	"os"

	"github.com/ishwar00/JackAnalyzer/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
	file         string
	onColumn     int
	onLine       int
}

func LexString(input string) *Lexer {
    l := &Lexer{
        input: input,
        onColumn: -1,
    }
    l.readChar()
    return l
}

func LexFile(fileName string) (*Lexer, error) {
	input, err := os.ReadFile(fileName)
	if err != nil {
		return &Lexer{}, err
	}

	l := &Lexer{
		input:    string(input),
		file:     fileName,
        onColumn: -1,
	}

	l.readChar() // initializing with first character
	return l, nil
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

    l.consumeComment()
	l.consumeWhiteSpace()

	switch l.char {
	case '{':
		tok = l.newToken(token.LBRACE, string(l.char))
	case '}':
		tok = l.newToken(token.RBRACE, string(l.char))
	case '(':
		tok = l.newToken(token.LPAREN, string(l.char))
	case ')':
		tok = l.newToken(token.RPAREN, string(l.char))
	case '[':
		tok = l.newToken(token.LBRACK, string(l.char))
	case ']':
		tok = l.newToken(token.RBRACK, string(l.char))
	case '.':
		tok = l.newToken(token.PERIOD, string(l.char))
	case ';':
		tok = l.newToken(token.SEMICO, string(l.char))
	case '-':
		tok = l.newToken(token.MINUS,  string(l.char))
	case '=':
		tok = l.newToken(token.EQ,     string(l.char))
    case '*':
		tok = l.newToken(token.ASTERI, string(l.char))
	case '/':
        tok = l.newToken(token.SLASH,  string(l.char))
	case '&':
		tok = l.newToken(token.AMPERS, string(l.char))
	case '|':
		tok = l.newToken(token.PIPE,   string(l.char))
	case '<':
		tok = l.newToken(token.LT,     string(l.char))
	case '>':
		tok = l.newToken(token.GT,     string(l.char))
	case '+':
		tok = l.newToken(token.PLUS,   string(l.char))
    case ',':
        tok = l.newToken(token.COMMA,  string(l.char))
	case '~':
		tok = l.newToken(token.TILDE,  string(l.char))
	case '"':
        l.readChar() // consuming opening double quote
        str := l.readString()
        if l.char != '"' { // there is no closing double quote
            tok = l.newToken(token.ILLEGAL, str)
            return tok
        } else { // there is closing double quote
            tok = l.newToken(token.STR_CONST, str)
        }
    case 0:
        tok = l.newToken(token.EOF, "")
        return tok
    default:
        if isLetter(l.char) && !isDigit(l.char) { // either identifier or keyword
            identifier := l.readIdentifier()
            tokenType := token.LookupIdent(identifier)
            tok = l.newToken(tokenType, identifier)
            return tok
        } else if isDigit(l.char) { // integer constant
            number := l.readNumber()
            tok = l.newToken(token.INT_CONST, number)
            return tok
        } else {
            tok = l.newToken(token.ILLEGAL, string(l.char))
        }
	}
	l.readChar() // advance to next character
	return tok
}

func (l *Lexer) readString() string {
    position := l.position
    for l.char != '"' && l.char != '\n' && l.char != 0 {
        l.readChar()
    }
    return l.input[position: l.position]
}

func (l *Lexer) readIdentifier() string {
    position := l.position
    for isLetter(l.char) {
        l.readChar()
    }
    return l.input[position: l.position]
}

func (l *Lexer) readNumber() string {
    position := l.position
    for isDigit(l.char) {
        l.readChar()
    }
    return l.input[position: l.position]
}

func (l *Lexer) consumeWhiteSpace() {
	for l.char == ' ' || l.char == '\n' || l.char == '\t' || l.char == '\r' {
		if l.char == '\n' {
			l.onColumn = -1 
			l.onLine++
		}
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
    if l.readPosition >= len(l.input) {
        return 0
    } 
    return l.input[l.readPosition]
}

func (l *Lexer) consumeComment() {
    l.consumeWhiteSpace()
    if l.char == '/' && l.peekChar() == '*' {
        l.consumeMultiLineComment()
        l.readChar() // consume *
        l.readChar() // consume /
        l.consumeComment() // going for next immediate comment
    } else if l.char == '/' && l.peekChar() == '/' {
        l.consumeSingleLineComment()
        l.consumeComment() // going for next immediate comment
        // not consuming current \n character, we will leave that
        // job to consumeWhiteSpace routine
    }
}

func (l *Lexer) consumeMultiLineComment() {
    for l.char != 0 && !(l.char == '*' && l.peekChar() == '/') {
        if l.char == '\n' {
            l.onLine++
            l.onColumn = -1
        }
        l.readChar()
    }
}

func (l *Lexer) consumeSingleLineComment() {
    for l.char != '\n' && l.char != 0 {
        l.readChar()
    }
}

func (l *Lexer) newToken(tokenType token.TokenType, tok string) token.Token {
    var onColumn int
    tokenLength := len(tok)
    if tokenLength == 1 && !isLetter(tok[0]){ // because l.position will be on symbol(eg: +, -), not after it.
        onColumn = l.onColumn
    } else {
        onColumn = l.onColumn - tokenLength // l.position is pointing just after token
    }

	return token.Token{
		Literal:  tok,
		Type:     tokenType,
		OnLine:   l.onLine,
		OnColumn: onColumn, 
		InFile:   l.file,
	}
}

func (l *Lexer) readChar() {
    if l.readPosition > len(l.input) {
        return
    }

	if l.readPosition == len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}
    l.position = l.readPosition
    l.readPosition++
    l.onColumn++
}

func isLetter(char byte) bool {
    return 'a' <= char && char <= 'z' || 
           'A' <= char && char <= 'Z' ||
           char == '_'                || 
           isDigit(char)
}

func isDigit(char byte) bool {
    return '0' <= char && char <= '9'
}


