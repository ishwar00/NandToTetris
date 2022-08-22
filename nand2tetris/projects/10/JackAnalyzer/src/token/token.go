package token

type TokenType int

type Token struct {
	Literal  string
	Type     TokenType
	OnLine   int
	OnColumn int
	InFile   string
}

const (
	// keywords tokens, their value lies between [0, 20]
	CLASS = iota // first token value 0
	FUNCTION
	METHOD
	CONSTRUCTOR
	FIELD
	STATIC
	VAR
	CHAR
	BOOLEAN
	INT
	TRUE
	FALSE
	NULL
	THIS
	LET
	IF
	ELSE
	RETURN
	DO
	VOID
	WHILE // last keyword token value 20

	// symbol tokens, range between [21, 39]
	LBRACE // starts with 21
	RBRACE
	LBRACK
	RBRACK
	RPAREN
	LPAREN
	PERIOD
	COMMA
	SEMICO
	PLUS
	MINUS
	ASTERI
	SLASH
	EQ
	TILDE
	PIPE // 36
	LT
	GT
	AMPERS // ends with 39

	// identifier token value is 40
	IDENT

	// constants value range between [41, 42]
	STR_CONST // 41
	INT_CONST // 42

	// helper tokens
	SYMBOLS
	KEYWORDS
	EOF
	ILLEGAL
)

var keywords = map[string]TokenType{
	"class":       CLASS,
	"constructor": CONSTRUCTOR,
	"function":    FUNCTION,
	"method":      METHOD,
	"field":       FIELD,
	"static":      STATIC,
	"var":         VAR,
	"int":         INT,
	"char":        CHAR,
	"boolean":     BOOLEAN,
	"void":        VOID,
	"true":        TRUE,
	"false":       FALSE,
	"null":        NULL,
	"this":        THIS,
	"let":         LET,
	"do":          DO,
	"if":          IF,
	"else":        ELSE,
	"while":       WHILE,
	"return":      RETURN,
}

func IsKeyword(tok string) bool {
	_, ok := keywords[tok]
	return ok
}

func LookupIdent(ident string) TokenType {
	if tokenType, ok := keywords[ident]; ok {
		return tokenType
	}
	return IDENT
}
