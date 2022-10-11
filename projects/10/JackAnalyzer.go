package main

import (
	"bufio"
	"fmt"
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

var isUnaryOp map[string]bool

var isOp map[string]bool

var isStatePrefix map[string]bool

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

	isOp = map[string]bool{
		"+":     true,
		"-":     true,
		"*":     true,
		"/":     true,
		"&amp;": true, // "&"
		"|":     true,
		"&lt;":  true,
		"&gt;":  true,
		"=":     true,
	}

	isUnaryOp = map[string]bool{
		"-": true,
		"~": true,
	}

	isStatePrefix = map[string]bool{
		"do":     true,
		"if":     true,
		"while":  true,
		"let":    true,
		"return": true,
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
		engine.NextToken()
		engine.CompileClass()
		// engine.Flush()
	}
	output.Flush()
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
	lastIdx    int // curLine[index:lastIdx)
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

// has more (幂等)
func (t *JackTokenizer) HasMoreTokens() bool {
	if t.index < t.lastIdx || t.lastIdx < len(t.curLine) { // [index, lastIdx)
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
		t.lastIdx = 0
		break
	}
	if len(t.curLine) == 0 {
		return false
	}
	return true
}

// move cursor to point to next Token
func (t *JackTokenizer) Advance() {
	if t.lastIdx >= len(t.curLine) {
		t.index = t.lastIdx
		if t.HasMoreTokens() {
			t.Advance()
		}
		return
	}
	// acquire token
	token := ""
	idx, line := t.lastIdx, t.curLine
	for idx < len(line) && line[idx] == ' ' {
		idx++
	}
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
	// update t.index, t.lastIdx
	t.index = idx
	t.lastIdx = nxt
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
		log.Fatalln("KeyWord fail", t.tokenType)
	}
	return t.keyWord
}

// returns the character which is the current token
// should be called only when TokenType() is SYMBOL
func (t *JackTokenizer) Symbol() string {
	if t.tokenType != SYMBOL {
		log.Fatalln("Symbol fail", t.tokenType)
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
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileClass fail")
	}
	c.output.WriteString("<class>\n")
	c.WriteKeyword()    //  class
	c.writeIdentifier() //  className
	c.WriteSymbol()     //  {
	for c.tokenizer.HasMoreTokens() && c.tokenizer.TokenType() == KEYWORD {
		token := c.tokenizer.KeyWord()
		if token != "static" && token != "field" {
			break
		}
		c.CompileClassVarDec()
	}
	for c.tokenizer.HasMoreTokens() && c.tokenizer.TokenType() == KEYWORD {
		token := c.tokenizer.KeyWord()
		if token != "constructor" && token != "function" && token != "method" {
			break
		}
		c.CompileSubroutine()
	}
	c.WriteSymbol() //  }
	c.output.WriteString("</class>\n")
}

func (c *CompilationEngine) CompileClassVarDec() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileClassVarDec fail")
	}
	c.output.WriteString("<classVarDec>\n")
	c.WriteKeyword()    //  static | field
	c.WriteType()       //  type(int | char | boolean | className)
	c.writeIdentifier() //  varName
	for c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == "," {
		c.WriteSymbol()     //  ,
		c.writeIdentifier() //  varName
	}
	c.WriteSymbol() //  ;
	c.output.WriteString("</classVarDec>\n")
}

func (c *CompilationEngine) CompileSubroutine() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileSubroutine fail")
	}
	c.output.WriteString("<subroutineDec>\n")
	c.WriteKeyword() //  constructor | function | method
	//  void | type(int | char | boolean | className)
	if c.tokenizer.TokenType() == KEYWORD && c.tokenizer.KeyWord() == "void" {
		c.WriteKeyword()
	} else {
		c.WriteType()
	}
	c.writeIdentifier() //  subroutineName
	c.WriteSymbol()     //  (
	c.CompileParameterList()
	c.WriteSymbol() //  )
	c.output.WriteString("<subroutineBody>\n")
	c.WriteSymbol() //  {
	for c.tokenizer.HasMoreTokens() && c.tokenizer.TokenType() == KEYWORD {
		token := c.tokenizer.KeyWord()
		if token != "var" {
			break
		}
		c.CompileVarDec()
	}
	c.CompileStatements()
	c.WriteSymbol() //  }
	c.output.WriteString("</subroutineBody>\n")
	c.output.WriteString("</subroutineDec>\n")
}

func (c *CompilationEngine) CompileParameterList() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileParameterList fail")
	}
	c.output.WriteString("<parameterList>\n")
	if c.tokenizer.TokenType() == IDENTIFIER ||
		(c.tokenizer.TokenType() == KEYWORD &&
			(c.tokenizer.KeyWord() == "int" || c.tokenizer.KeyWord() == "char" || c.tokenizer.KeyWord() == "boolean")) {
		c.WriteType()       //  type
		c.writeIdentifier() //  varName
		for c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == "," {
			c.WriteSymbol()     //  ,
			c.WriteType()       //  type
			c.writeIdentifier() //  varName
		}
	}
	c.output.WriteString("</parameterList>\n")
}

func (c *CompilationEngine) CompileVarDec() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileVarDec fail")
	}
	c.output.WriteString("<varDec>\n")
	c.WriteKeyword()    //  var
	c.WriteType()       //  type
	c.writeIdentifier() //  varName
	for c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == "," {
		c.WriteSymbol()     //  ,
		c.writeIdentifier() //  varName
	}
	c.WriteSymbol() //  ;
	c.output.WriteString("</varDec>\n")
}

func (c *CompilationEngine) CompileStatements() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileStatements fail")
	}
	c.output.WriteString("<statements>\n")
	for c.tokenizer.HasMoreTokens() && c.tokenizer.TokenType() == KEYWORD {
		token := c.tokenizer.KeyWord()
		if !isStatePrefix[token] {
			break
		}
		switch token {
		case "do":
			c.CompileDo()
		case "let":
			c.CompileLet()
		case "while":
			c.CompileWhile()
		case "if":
			c.CompileIf()
		case "return":
			c.CompileReturn()
		}
	}
	c.output.WriteString("</statements>\n")
}

func (c *CompilationEngine) CompileDo() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileDo fail")
	}
	c.output.WriteString("<doStatement>\n")
	c.WriteKeyword() //  do
	//  subroutineCall
	token0 := c.tokenizer.Identifier()
	c.NextToken()
	c.output.WriteString(fmt.Sprintf("<identifier> %s </identifier>\n", token0))
	tokenType1 := c.tokenizer.TokenType()
	if tokenType1 == SYMBOL && c.tokenizer.Symbol() == "(" { //  subroutineCall case 1
		// ( expressionList )
		c.output.WriteString("<symbol> ( </symbol>\n")
		c.NextToken() //  (
		c.CompileExpressionList()
		c.NextToken() //  )
		c.output.WriteString("<symbol> ) </symbol>\n")
	} else if tokenType1 == SYMBOL && c.tokenizer.Symbol() == "." { //  subroutineCall case 2
		// .subroutineName
		c.output.WriteString("<symbol> . </symbol>\n")
		c.NextToken() //  .
		subroutineName := c.tokenizer.Identifier()
		c.output.WriteString(fmt.Sprintf("<identifier> %s </identifier>\n", subroutineName))
		c.NextToken()
		// ( expressionList )
		c.output.WriteString("<symbol> ( </symbol>\n")
		c.NextToken() //  (
		c.CompileExpressionList()
		c.NextToken() //  )
		c.output.WriteString("<symbol> ) </symbol>\n")
	}
	c.WriteSymbol() //  ;
	c.output.WriteString("</doStatement>\n")
}

func (c *CompilationEngine) CompileLet() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileLet fail")
	}
	c.output.WriteString("<letStatement>\n")
	c.WriteKeyword()    //  let
	c.writeIdentifier() //  varName
	if c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == "[" {
		c.WriteSymbol() //  [
		c.CompileExpression()
		c.WriteSymbol() //  ]
	}
	c.WriteSymbol() //  =
	c.CompileExpression()
	c.WriteSymbol() //  ;
	c.output.WriteString("</letStatement>\n")
}

func (c *CompilationEngine) CompileWhile() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileWhile fail")
	}
	c.output.WriteString("<whileStatement>\n")
	c.WriteKeyword() //  while
	c.WriteSymbol()  //  (
	c.CompileExpression()
	c.WriteSymbol() //  )
	c.WriteSymbol() // {
	c.CompileStatements()
	c.WriteSymbol() // }
	c.output.WriteString("</whileStatement>\n")
}

func (c *CompilationEngine) CompileReturn() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileReturn fail")
	}
	c.output.WriteString("<returnStatement>\n")
	c.WriteKeyword() //  return
	if !(c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == ";") {
		c.CompileExpression()
	}
	c.WriteSymbol() //  ;
	c.output.WriteString("</returnStatement>\n")
}

func (c *CompilationEngine) CompileIf() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileIf fail")
	}
	c.output.WriteString("<ifStatement>\n")
	c.WriteKeyword() //  if
	c.WriteSymbol()  //  (
	c.CompileExpression()
	c.WriteSymbol() //  )
	c.WriteSymbol() //  {
	c.CompileStatements()
	c.WriteSymbol() //  }
	if c.tokenizer.TokenType() == KEYWORD && c.tokenizer.KeyWord() == "else" {
		c.WriteKeyword() //  else
		c.WriteSymbol()  //  {
		c.CompileStatements()
		c.WriteSymbol() //  }
	}
	c.output.WriteString("</ifStatement>\n")
}

func (c *CompilationEngine) CompileTerm() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileTerm fail")
	}
	c.output.WriteString("<term>\n")
	tokenType := c.tokenizer.TokenType()
	switch tokenType {
	case SYMBOL:
		token := c.tokenizer.Symbol()
		c.NextToken()
		c.output.WriteString(fmt.Sprintf("<symbol> %s </symbol>\n", token))
		if isUnaryOp[token] {
			// unaryOp must be followed by term
			c.CompileTerm()
		} else if token == "(" {
			c.CompileExpression()
			c.NextToken() //  )
			c.output.WriteString("<symbol> ) </symbol>\n")
		}
	case INT_CONST:
		token := c.tokenizer.IntVal()
		c.NextToken()
		c.output.WriteString(fmt.Sprintf("<integerConstant> %v </integerConstant>\n", token))
	case STRING_CONST:
		token := c.tokenizer.StringVal()
		c.NextToken()
		c.output.WriteString(fmt.Sprintf("<stringConstant> %s </stringConstant>\n", token))
	case KEYWORD:
		token := c.tokenizer.KeyWord()
		// asset token in (true, false, null, this)
		c.NextToken()
		c.output.WriteString(fmt.Sprintf("<keyword> %s </keyword>\n", token))
	case IDENTIFIER:
		token0 := c.tokenizer.Identifier()
		c.NextToken()
		c.output.WriteString(fmt.Sprintf("<identifier> %s </identifier>\n", token0))
		// assert c.tokenizer.HasMoreToken() == true, there is at least a ';' to follow
		tokenType1 := c.tokenizer.TokenType()
		if tokenType1 == SYMBOL && c.tokenizer.Symbol() == "[" { //   varName[ expression ]
			c.output.WriteString("<symbol> [ </symbol>\n")
			c.NextToken() //  [
			c.CompileExpression()
			c.NextToken() //  ]
			c.output.WriteString("<symbol> ] </symbol>\n")
		} else if tokenType1 == SYMBOL && c.tokenizer.Symbol() == "(" { //  subroutineCall case 1
			// ( expressionList )
			c.output.WriteString("<symbol> ( </symbol>\n")
			c.NextToken() //  (
			c.CompileExpressionList()
			c.NextToken() //  )
			c.output.WriteString("<symbol> ) </symbol>\n")
		} else if tokenType1 == SYMBOL && c.tokenizer.Symbol() == "." { //  subroutineCall case 2
			// .subroutineName
			c.output.WriteString("<symbol> . </symbol>\n")
			c.NextToken() //  .
			subroutineName := c.tokenizer.Identifier()
			c.output.WriteString(fmt.Sprintf("<identifier> %s </identifier>\n", subroutineName))
			c.NextToken()
			// ( expressionList )
			c.output.WriteString("<symbol> ( </symbol>\n")
			c.NextToken() //  (
			c.CompileExpressionList()
			c.NextToken() //  )
			c.output.WriteString("<symbol> ) </symbol>\n")
		}
	}
	c.output.WriteString("</term>\n")
}

func (c *CompilationEngine) CompileExpression() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileExpression fail")
	}
	c.output.WriteString("<expression>\n")
	c.CompileTerm()
	for c.tokenizer.HasMoreTokens() && c.tokenizer.TokenType() == SYMBOL {
		token := c.tokenizer.Symbol() //  op
		if !isOp[token] {
			break
		}
		c.output.WriteString(fmt.Sprintf("<symbol> %s </symbol>\n", token))
		c.NextToken()
		c.CompileTerm()
	}
	c.output.WriteString("</expression>\n")
}

func (c *CompilationEngine) CompileExpressionList() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileExpressionList fail")
	}
	c.output.WriteString("<expressionList>\n")
	if c.tokenizer.TokenType() != SYMBOL || c.tokenizer.Symbol() != ")" {
		c.CompileExpression()
		for c.tokenizer.HasMoreTokens() && c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == "," {
			c.output.WriteString("<symbol> , </symbol>\n")
			c.NextToken() //  ,
			c.CompileExpression()
		}
	}
	c.output.WriteString("</expressionList>\n")
}

func (c *CompilationEngine) NextToken() {
	if c.tokenizer.HasMoreTokens() {
		c.tokenizer.Advance()
	} else {
		log.Fatal("NextToken fail")
	}
}

func (c *CompilationEngine) WriteSymbol() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("WriteSymbol fail")
	}
	token := c.tokenizer.Symbol()
	c.output.WriteString(fmt.Sprintf("<symbol> %s </symbol>\n", token))
	c.NextToken()
}

func (c *CompilationEngine) writeIdentifier() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("writeIdentifier fail")
	}
	token := c.tokenizer.Identifier()
	c.output.WriteString(fmt.Sprintf("<identifier> %s </identifier>\n", token))
	c.NextToken()
}

func (c *CompilationEngine) WriteKeyword() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("WriteKeyword fail")
	}
	token := c.tokenizer.KeyWord()
	c.output.WriteString(fmt.Sprintf("<keyword> %s </keyword>\n", token))
	c.NextToken()
}

func (c *CompilationEngine) WriteType() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("WriteType fail")
	}
	if c.tokenizer.TokenType() == KEYWORD {
		c.WriteKeyword()
	} else if c.tokenizer.TokenType() == IDENTIFIER {
		c.writeIdentifier()
	} else {
		log.Fatal("WriteType encounter invalid token", c.tokenizer.TokenType())
	}
}

func (c *CompilationEngine) Flush() {
	c.output.Flush()
}
