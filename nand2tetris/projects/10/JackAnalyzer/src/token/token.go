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
	// symbols ( } ; = ...
	LBRACE = iota
	RBRACE
	RPAREN
	LPAREN
	RBRACK
	LBRACK
	PERIOD
	COMMA
	SEMICO

	// operators, +, -, *, /, & ...
	PLUS
	MINUS
	ASTERI
	SLASH
	AMPERS
	PIPE
	GT
	LT
	EQ
	TILDE

	// constants "abc", 4532 ...
	STR_CONST
	INT_CONST

	// indentifiers and keywords, abc_d, if, class ...
	IDENT
	CLASS
	CONSTRUCTOR
	FUNCTION
	METHOD
	FIELD
	STATIC
	VAR
	INT
	CHAR
	BOOLEAN
	VOID
	TRUE
	FALSE
	NULL
	THIS
	LET
	DO
	IF
	ELSE
	WHILE
	RETURN

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

func LookupIdent(ident string) TokenType {
	if tokenType, ok := keywords[ident]; ok {
		return tokenType
	}
	return IDENT
}
