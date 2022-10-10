package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

const (
	KEYWORD      = "keyword"
	SYMBOL       = "symbol"
	IDENTIFIER   = "identifier"
	INT_CONST    = "integerConstant"
	STRING_CONST = "stringConstant"
)

var isSymbol map[string]bool

var isKeyword map[string]bool

var escape map[string]string

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
		"-": true,
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
		"void":        true,
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

	escape = map[string]string{
		"<":  "&lt;",
		">":  "&gt;",
		"\"": "&quot;",
		"&":  "&amp;",
	}
}

func main() {
	filePath := os.Args[1]
	files := readDir(filePath)
	// dir := filepath.Dir(filePath)
	baseName := filepath.Base(filePath)
	outputfile, err := os.OpenFile(baseName+".xml", os.O_WRONLY|os.O_CREATE, 777)
	if err != nil {
		log.Fatal(err)
	}
	output := bufio.NewWriter(outputfile)
	for _, file := range files {
		tokenizer := JackTokenizerConstructor(bufio.NewScanner(file))
		engine := CompilationEngineConstructor(&tokenizer, output)
		// one .jack file should contain only one class
		engine.CompileClass()
	}
}

//*******************************************************************************************************************//
//**************************************************** Utils ********************************************************//

func testTokenizer() {
	filePath := os.Args[1]
	files := readDir(filePath)
	// dir := filepath.Dir(filePath)
	baseName := filepath.Base(filePath)
	output, err := os.OpenFile(baseName+".out.xml", os.O_WRONLY|os.O_CREATE, 777)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		tokenizer := JackTokenizerConstructor(bufio.NewScanner(file))
		tokenizer.WriteToXML(bufio.NewWriter(output))
	}
}

func readDir(path string) []*os.File {
	fileExt := filepath.Ext(path)
	if fileExt == ".jack" {
		file, _ := os.Open(path)
		return []*os.File{file}
	}
	//
	filenames, _ := filepath.Glob(path + "/*")
	ans := []*os.File{}
	for i := range filenames {
		ext := filepath.Ext(filenames[i])
		if ext != ".jack" {
			continue
		}
		file1, _ := os.Open(filenames[i])
		ans = append(ans, file1)
	}
	return ans
}

//*******************************************************************************************************************//
//************************************************** JackTokenizer **************************************************//

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

func JackTokenizerConstructor(input *bufio.Scanner) JackTokenizer {
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
	t.curLine = ""
	for t.input.Scan() {
		str := t.input.Text()
		// trim rear comment
		idx := strings.Index(str, "//")
		if idx != -1 {
			str = str[:idx]
		}
		// trim lead empty char and rear empty char
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}
		if isComment && strings.HasSuffix(str, "*/") {
			isComment = false
			continue
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
		if line[nxt] == '"' {
			nxt++
			for nxt < len(line) && line[nxt] != '"' {
				nxt++
			}
			nxt++
			break
		}
		if line[nxt] == ' ' || isSymbol[string(line[nxt])] {
			break
		}
	}
	if idx == nxt { // it must be a symbol
		nxt++
	}
	token = line[idx:nxt]
	// update t.index
	for nxt < len(line) && line[nxt] == ' ' {
		nxt++
	}
	t.index = nxt
	// process token
	if isKeyword[token] {
		t.tokenType = KEYWORD
		t.keyWord = token
	} else if isSymbol[token] {
		t.tokenType = SYMBOL
		if escape[token] != "" {
			token = escape[token]
		}
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
		log.Fatalf("unknown token: %v in %v, index: %v \n", token, t.curLine, t.index)
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

func (t *JackTokenizer) WriteToXML(output *bufio.Writer) error {
	output.WriteString("<tokens>\n")
	for t.HasMoreTokens() {
		t.Advance()
		switch t.TokenType() {
		case KEYWORD:
			val := t.KeyWord()
			output.WriteString("<" + KEYWORD + "> " + val + " </" + KEYWORD + ">\n")
		case SYMBOL:
			val := t.Symbol()
			output.WriteString("<" + SYMBOL + "> " + val + " </" + SYMBOL + ">\n")
		case IDENTIFIER:
			val := t.Identifier()
			output.WriteString("<" + IDENTIFIER + "> " + val + " </" + IDENTIFIER + ">\n")
		case INT_CONST:
			val := strconv.Itoa(t.IntVal())
			output.WriteString("<" + INT_CONST + "> " + val + " </" + INT_CONST + ">\n")
		case STRING_CONST:
			val := t.StringVal()
			output.WriteString("<" + STRING_CONST + "> " + val + " </" + STRING_CONST + ">\n")
		}
	}
	output.WriteString("</tokens>\n")
	return output.Flush()
}

//*******************************************************************************************************************//
//************************************************ CompilationEngine ************************************************//

type CompilationEngine struct {
	output    *bufio.Writer
	tokenizer *JackTokenizer
}

func CompilationEngineConstructor(tokenizer *JackTokenizer, output *bufio.Writer) CompilationEngine {
	return CompilationEngine{
		output:    output,
		tokenizer: tokenizer,
	}
}

func (c *CompilationEngine) CompileClass() {
	c.output.WriteString("<class>\n")
	c.CompileTerm() //  class
	c.CompileTerm() //  className
	c.CompileTerm() //  {
	c.CompileClassVarDec()
	c.CompileSubroutine()
	c.CompileTerm() //  }
	c.output.WriteString("</class>\n")
}

func (c *CompilationEngine) CompileClassVarDec() {
	c.output.WriteString("<classVarDec>\n")
	c.CompileTerm() //  static | field
	c.CompileTerm() //  type(int | char | boolean | className)
	c.CompileTerm() //  varName
	for c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == "," {
		c.CompileTerm() //  ,
		c.CompileTerm() //  varName
	}
	c.output.WriteString("</classVarDec>\n")
}

func (c *CompilationEngine) CompileSubroutine() {
	c.output.WriteString("<subroutineDec>\n")
	c.CompileTerm() //  constructor | function | method
	c.CompileTerm() //  void | type(int | char | boolean | className)
	c.CompileTerm() //  subroutineName
	c.CompileTerm() //  (
	c.CompileParameterList()
	c.CompileTerm() //  )
	c.CompileTerm() //  {
	c.CompileVarDec()
	c.CompileStatements()
	c.CompileTerm() //  }
	c.output.WriteString("</subroutineDec>\n")
}

func (c *CompilationEngine) CompileParameterList() {
	c.output.WriteString("<parameterList>\n")
	if c.tokenizer.TokenType() == IDENTIFIER ||
		(c.tokenizer.TokenType() == SYMBOL &&
			(c.tokenizer.Symbol() == "int" || c.tokenizer.Symbol() == "char" || c.tokenizer.Symbol() == "boolean")) {
		c.CompileTerm() //  type
		c.CompileTerm() //  varName
		for c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == "," {
			c.CompileTerm() //  ,
			c.CompileTerm() //  type
			c.CompileTerm() //  varName
		}
	}
	c.output.WriteString("</parameterList>\n")
}

func (c *CompilationEngine) CompileVarDec() {
	c.output.WriteString("<varDec>\n")
	c.CompileTerm() //  var
	c.CompileTerm() //  type
	c.CompileTerm() //  varName
	for c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == "," {
		c.CompileTerm() //  ,
		c.CompileTerm() //  varName
	}
	c.output.WriteString("</varDec>\n")
}

func (c *CompilationEngine) CompileStatements() {
	c.output.WriteString("<statements>\n")
	token := c.tokenizer.KeyWord()
	switch token {
	case "do":
	case "let":
	case "while":
	case "if":
	case "return":
	}
	c.output.WriteString("</statements>\n")
}

func (c *CompilationEngine) CompileDo() {

}

func (c *CompilationEngine) CompileLet() {

}

func (c *CompilationEngine) CompileWhile() {

}

func (c *CompilationEngine) CompileReturn() {

}

func (c *CompilationEngine) CompileIf() {

}

func (c *CompilationEngine) CompileExpression() {

}

func (c *CompilationEngine) CompileTerm() {

}

func (c *CompilationEngine) CompileExpressionList() {

}

func (c *CompilationEngine) NextToken() {
	if c.tokenizer.HasMoreTokens() {
		c.tokenizer.Advance()
	}
}

func (c *CompilationEngine) Flush() {
	c.output.Flush()
}
