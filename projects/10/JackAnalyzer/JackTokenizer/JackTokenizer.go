package jacktokenizer

import (
	"bufio"
	"log"
	"strconv"
	"strings"
	"unicode"
)

const (
	KEYWORD      = "KEYWORD"
	SYMBOL       = "SYMBOL"
	IDENTIFIER   = "IDENTIFIER"
	INT_CONST    = "INT_CONST"
	STRING_CONST = "STRING_CONST"
)

var isSymbol map[string]bool

var isKeyword map[string]bool

func init() {
	isSymbol = map[string]bool{
		"{": true,
		"}": true,
		"(": true,
		")": true,
		"[": true,
		"]": true,
		".": true,
		",": true,
		";": true,
		"+": true,
		"*": true,
		"/": true,
		"&": true,
		"|": true,
		"<": true,
		">": true,
		"=": true,
		"~": true,
	}

	isKeyword = map[string]bool{
		"class":       true,
		"constructor": true,
		"function":    true,
		"method":      true,
		"field":       true,
		"static":      true,
		"var":         true,
		"int":         true,
		"char":        true,
		"boolean":     true,
		"true":        true,
		"false":       true,
		"null":        true,
		"this":        true,
		"let":         true,
		"do":          true,
		"if":          true,
		"else":        true,
		"while":       true,
		"return":      true,
	}
}

type JackTokenizer struct {
	input      *bufio.Scanner
	curLine    string // current line (without lead or rear space)
	index      int
	tokenType  string
	keyWord    string
	symbol     string
	intVal     int
	stringVal  string
	identifier string
}

func Constructor(input *bufio.Scanner) JackTokenizer {
	return JackTokenizer{
		input: input,
	}
}

// has more
func (t *JackTokenizer) HasMoreTokens() bool {
	if t.index < len(t.curLine) {
		return true
	}
	isComment := false
	for t.input.Scan() {
		str := t.input.Text()
		// trim lead empty char and rear empty char
		str = strings.Trim(str, " ")
		if str == "" {
			continue
		}
		if strings.HasSuffix(str, "*/") {
			isComment = false
		}
		if isComment {
			continue
		}
		// process 3 different comments
		if strings.HasPrefix(str, "//") {
			continue
		} else if strings.HasPrefix(str, "/*") {
			if strings.HasSuffix(str, "*/") {
				continue
			}
			isComment = true
			continue
		}
		// assert str must a valid command
		t.curLine = str
		t.index = 0
		break
	}
	if len(t.curLine) == 0 {
		return false
	}
	return true
}

// move cursor to point to next Token
func (t *JackTokenizer) Advance() {
	if t.index >= len(t.curLine) {
		log.Fatalln("Advance fail")
	}
	// acquire token
	idx, line := t.index, t.curLine
	token := ""
	nxt := idx
	for ; nxt < len(line); nxt++ {
		if line[nxt] == ' ' || isSymbol[strconv.Itoa(int(line[nxt]))] {
			break
		}
	}
	if idx == nxt { // it must be a symbol
		nxt++
	}
	t.index = nxt
	token = line[idx:nxt]
	// process token
	if isKeyword[token] {
		t.tokenType = KEYWORD
		t.keyWord = token
	} else if isSymbol[token] {
		t.tokenType = SYMBOL
		t.symbol = token
	} else if unicode.IsDigit(rune(token[0])) {
		t.tokenType = INT_CONST
		val, err := strconv.Atoi(token)
		if err != nil {
			log.Fatalf("strconv.Atoi fail, %v \n", token)
		}
		t.intVal = val
	} else if token[0] == '"' {
		t.tokenType = STRING_CONST
		t.stringVal = token[1 : len(token)-1]
	} else if unicode.IsLetter(rune(token[0])) || token[0] == '_' {
		t.tokenType = IDENTIFIER
		t.identifier = token
	} else {
		log.Fatalf("unknown token: %v \n", token)
	}
}

// Returns the type of the current token
func (t *JackTokenizer) TokenType() string {
	return t.tokenType
}

// returns the keywords which is the current token
// should be called only when TokenType() is KEYWORD
func (t *JackTokenizer) KeyWord() string {
	if t.tokenType != KEYWORD {
		log.Fatalln("KeyWord fail")
	}
	return t.keyWord
}

// returns the character which is the current token
// should be called only when TokenType() is SYMBOL
func (t *JackTokenizer) Symbol() string {
	if t.tokenType != SYMBOL {
		log.Fatalln("Symbol fail")
	}
	return t.symbol
}

// returns the identifier which is the current token
// should be called only when TokenType() is IDENTIFIER
func (t *JackTokenizer) Identifier() string {
	if t.tokenType != IDENTIFIER {
		log.Fatalln("Identifier fail")
	}
	return t.identifier
}

// returns the integer value of the current token
// should be called only when TokenType() is INT_CONST
func (t *JackTokenizer) IntVal() int {
	if t.tokenType != INT_CONST {
		log.Fatalln("IntVal fail")
	}
	return t.intVal
}

// returns the string value of the current token (without the double quotes)
// should be called only when TokenType() is STRING_CONST
func (t *JackTokenizer) StringVal() string {
	if t.tokenType != STRING_CONST {
		log.Fatalln("StringVal fail")
	}
	return t.stringVal
}
