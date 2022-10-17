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

// variable kind
const (
	// class scope
	STATIC = "static"
	FIELD  = "field"
	// subroutine scope: method, function, constructor
	ARG  = "arg" // function argument
	VAR  = "var" // local variable
	NONE = "NONE"
)

// segment
const (
	LOCAL    = "local"
	ARGUMENT = "argument"
	// STATIC   = "static"
	THIS     = "this"
	THAT     = "that"
	POINTER  = "pointer"
	CONSTANT = "constant"
	TEMP     = "temp"
)

var isSymbol map[string]bool

var isKeyword map[string]bool

var isUnaryOp map[string]bool

var isOp map[string]bool

var isStatePrefix map[string]bool

var escape map[string]string

var kind2Seg map[string]string

var isBuiltInClass map[string]bool

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

	kind2Seg = map[string]string{
		"var":    "local",
		"arg":    "argument",
		"field":  "this", // can't access a field variable out the class in jack
		"static": "static",
	}

	isBuiltInClass = map[string]bool{
		"Math":     true,
		"String":   true,
		"Array":    true,
		"Output":   true,
		"Screen":   true,
		"Keyboard": true,
		"Memory":   true,
		"Sys":      true,
	}
}

func main() {
	filePath := os.Args[1]
	files, newFiles := readDir(filePath)
	// dir := filepath.Dir(filePath)
	// baseName := filepath.Base(filePath)
	for i, file := range files {
		tokenizer := JackTokenizerConstructor(bufio.NewScanner(file))

		outputfile, err := os.OpenFile(newFiles[i], os.O_WRONLY|os.O_CREATE, 777)
		if err != nil {
			log.Fatal(err)
		}
		output := bufio.NewWriter(outputfile)

		engine := CompilationEngineConstructor(&tokenizer, output)
		// one .jack file should contain only one class
		engine.NextToken()
		engine.CompileClass()
		// engine.Flush()

		output.Flush()
	}
}

//*******************************************************************************************************************//
//**************************************************** Utils ********************************************************//

func testTokenizer() {
	filePath := os.Args[1]
	files, _ := readDir(filePath)
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

func readDir(path string) ([]*os.File, []string) {
	fileExt := filepath.Ext(path)
	if fileExt == ".jack" {
		file, _ := os.Open(path)
		fileDir := filepath.Dir(path)
		fileName := strings.TrimSuffix(filepath.Base(path), ".jack")
		return []*os.File{file}, []string{fileDir + "/" + fileName + ".vm"}
	}
	// path is a directory
	newFiles := make([]string, 0)
	filenames, _ := filepath.Glob(path + "/*")
	ans := []*os.File{}
	for i := range filenames {
		ext := filepath.Ext(filenames[i])
		if ext != ".jack" {
			continue
		}
		file1, _ := os.Open(filenames[i])
		ans = append(ans, file1)
		fileName := strings.TrimSuffix(filenames[i], ".jack") + ".vm"
		newFiles = append(newFiles, fileName)
	}
	return ans, newFiles
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
	ST        *SymbolTable
	writer    *VMWriter
	tokenizer *JackTokenizer
}

func CompilationEngineConstructor(tokenizer *JackTokenizer, output *bufio.Writer) CompilationEngine {
	tmp := VMWriterConstructor(output)
	return CompilationEngine{
		ST:        SymbolTableConstructor(),
		writer:    &tmp,
		tokenizer: tokenizer,
	}
}

func (c *CompilationEngine) CompileClass() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileClass fail")
	}
	c.WriteKeyword()                 //  class
	className := c.writeIdentifier() //  className
	c.ST.startClass(className)
	c.NextToken() //  {
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
	c.NextToken() //  }
}

func (c *CompilationEngine) CompileClassVarDec() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileClassVarDec fail")
	}
	kindName := c.WriteKeyword() //  static | field
	typeName := c.WriteType()    //  type(int | char | boolean | className)
	c.ST.curKind, c.ST.curType = kindName, typeName
	c.writeIdentifier() //  varName
	for c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == "," {
		c.NextToken() //  ,
		c.ST.curKind, c.ST.curType = kindName, typeName
		c.writeIdentifier() //  varName
	}
	c.NextToken() //  ;
}

func (c *CompilationEngine) CompileSubroutine() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileSubroutine fail")
	}
	subroutineType := c.WriteKeyword() //  constructor | function | method
	//  void | type(int | char | boolean | className)
	if c.tokenizer.TokenType() == KEYWORD && c.tokenizer.KeyWord() == "void" {
		c.WriteKeyword()
	} else {
		c.WriteType()
	}
	funcName := c.writeIdentifier() //  subroutineName
	c.ST.startSubroutine(funcName)
	c.NextToken() //  (
	if subroutineType == "method" {
		c.ST.Define(THIS, ARG, c.ST.className)
	}
	c.CompileParameterList()
	c.NextToken() //  )
	//*** local variables declaration
	c.NextToken() //  {
	nLocals := 0
	for c.tokenizer.HasMoreTokens() && c.tokenizer.TokenType() == KEYWORD {
		token := c.tokenizer.KeyWord()
		if token != "var" {
			break
		}
		nLocals += c.CompileVarDec()
	}
	c.writer.writeFunction(c.ST.className+"."+funcName, nLocals)
	//*** function's commands/statements
	if subroutineType == "constructor" {
		// extra special logic of constructor: allocate new memory block
		c.writer.writePush(CONSTANT, c.ST.VarCount(FIELD))
		c.writer.writeCall("Memory.alloc", 1)
		c.writer.writePop(POINTER, 0)
	} else if subroutineType == "method" {
		// extra special logic of method: set THIS = argument 0
		c.writer.writePush(ARGUMENT, 0)
		c.writer.writePop(POINTER, 0)
	}
	c.CompileStatements()
	c.NextToken() //  }
}

func (c *CompilationEngine) CompileParameterList() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileParameterList fail")
	}
	if c.tokenizer.TokenType() == IDENTIFIER ||
		(c.tokenizer.TokenType() == KEYWORD &&
			(c.tokenizer.KeyWord() == "int" || c.tokenizer.KeyWord() == "char" || c.tokenizer.KeyWord() == "boolean")) {
		c.WriteType()       //  type
		c.ST.curKind = ARG  //  record current symbol's kind
		c.writeIdentifier() //  varName
		for c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == "," {
			c.NextToken()       //  ,
			c.WriteType()       //  type
			c.ST.curKind = ARG  //  record current symbol's kind
			c.writeIdentifier() //  varName
		}
	}
}

// return the number of variable declared by this invocation
func (c *CompilationEngine) CompileVarDec() int {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileVarDec fail")
	}
	kindName := c.WriteKeyword() //  var
	typeName := c.WriteType()    //  type
	//  record current symbol's kind (be careful)
	c.ST.curKind, c.ST.curType = kindName, typeName
	c.writeIdentifier() //  varName
	cnt := 1
	for c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == "," {
		//  record current symbol's kind (be careful)
		c.ST.curKind, c.ST.curType = kindName, typeName
		c.NextToken()       //  ,
		c.writeIdentifier() //  varName
		cnt++
	}
	c.NextToken() //  ;
	return cnt
}

func (c *CompilationEngine) CompileStatements() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileStatements fail")
	}
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
}

func (c *CompilationEngine) CompileDo() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileDo fail")
	}
	c.WriteKeyword() //  do
	//  subroutineCall
	c.CompileTerm()
	c.writer.writePop(TEMP, 0) //  won't assign top-most value of stack
	c.NextToken()              //  ;
}

func (c *CompilationEngine) CompileLet() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileLet fail")
	}
	c.WriteKeyword()            //  let
	name := c.writeIdentifier() //  varName
	kindName, index := c.ST.KindOf(name), c.ST.IndexOf(name)
	updateArray := false
	if c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == "[" {
		updateArray = true
		// special logic for array literial
		c.NextToken() //  [
		c.CompileExpression()
		c.writer.writePush(kind2Seg[kindName], index)
		c.writer.WriteArithmetic("add")
		// c.writer.writePop("pointer", 1) // step*Ⅰ
		c.NextToken() //  ]
	}
	c.NextToken() //  =
	c.CompileExpression()
	if updateArray {
		// use temp 0 segement store the result of expression in "let varName([exp])? = expression"
		// the 'that' segment(assigned at step*Ⅰ) be override, when in case of "let a[i] = b[j]"
		c.writer.writePop(TEMP, 0) // after this operation, the top-most value of stack is the address of array
		c.writer.writePop(POINTER, 1)
		c.writer.writePush(TEMP, 0)
		c.writer.writePop(THAT, 0) // that 0 is a pointer that pointed to the base address of current array
	} else {
		c.writer.writePop(kind2Seg[kindName], index)
	}
	c.NextToken() //  ;
}

func (c *CompilationEngine) CompileWhile() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileWhile fail")
	}
	l1 := fmt.Sprintf("whileL1.%v", c.ST.whileCnt)
	l2 := fmt.Sprintf("whileL2.%v", c.ST.whileCnt)
	c.ST.whileCnt++
	c.WriteKeyword() //  while
	c.NextToken()    //  (
	c.writer.writeLabel(l1)
	c.CompileExpression()           // cond
	c.writer.WriteArithmetic("not") // ~cond
	c.writer.writeIfGoto(l2)
	c.NextToken() //  )
	c.NextToken() // {
	c.CompileStatements()
	c.NextToken() // }
	c.writer.writeGoto(l1)
	c.writer.writeLabel(l2)
}

func (c *CompilationEngine) CompileReturn() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileReturn fail")
	}
	c.WriteKeyword() //  return
	if !(c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == ";") {
		c.CompileExpression()
	} else {
		// void subroutine must push a value into stack
		c.writer.writePush(CONSTANT, 0)
	}
	c.NextToken() //  ;
	c.writer.writeReturn()
}

func (c *CompilationEngine) CompileIf() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileIf fail")
	}
	ifTrue := fmt.Sprintf("if_true_%v", c.ST.ifCnt)
	ifFalse := fmt.Sprintf("if_false_%v", c.ST.ifCnt)
	ifEnd := fmt.Sprintf("if_end_%v", c.ST.ifCnt)
	c.ST.ifCnt++
	c.WriteKeyword()      //  if
	c.NextToken()         //  (
	c.CompileExpression() //cond
	c.writer.writeIfGoto(ifTrue)
	c.writer.writeGoto(ifFalse)
	c.writer.writeLabel(ifTrue)
	c.NextToken() //  )
	c.NextToken() //  {
	c.CompileStatements()
	c.NextToken() //  }
	c.writer.writeGoto(ifEnd)
	c.writer.writeLabel(ifFalse)
	if c.tokenizer.TokenType() == KEYWORD && c.tokenizer.KeyWord() == "else" {
		c.WriteKeyword() //  else
		c.NextToken()    //  {
		c.CompileStatements()
		c.NextToken() //  }
	}
	c.writer.writeLabel(ifEnd)
}

func (c *CompilationEngine) CompileTerm() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileTerm fail")
	}
	tokenType := c.tokenizer.TokenType()
	switch tokenType {
	case SYMBOL:
		token := c.tokenizer.Symbol()
		c.NextToken()
		if isUnaryOp[token] {
			// unaryOp must be followed by term
			c.CompileTerm()
			if token == "~" {
				c.writer.WriteArithmetic("not")
			} else { // token == "-" <=> ~x+1
				c.writer.WriteArithmetic("neg")
			}
		} else if token == "(" {
			c.CompileExpression()
			c.NextToken() //  )
		}
	case INT_CONST:
		token := c.tokenizer.IntVal()
		c.NextToken()
		c.writer.writePush(CONSTANT, token)
	case STRING_CONST:
		token := c.tokenizer.StringVal()
		c.NextToken()
		c.writer.writePush(CONSTANT, len(token))
		c.writer.writeCall("String.new", 1)
		for i := range token {
			c.writer.writePush(CONSTANT, int(token[i]))
			c.writer.writeCall("String.appendChar", 2)
		}
	case KEYWORD:
		token := c.tokenizer.KeyWord()
		// assert token in (true, false, null, this) // no that
		c.NextToken()
		switch token {
		case "true":
			c.writer.writePush(CONSTANT, 0)
			c.writer.WriteArithmetic("not")
		case "false", "null":
			c.writer.writePush(CONSTANT, 0)
		case "this": // in case of "return this;"
			c.writer.writePush(POINTER, 0)
		default:
			fmt.Println("unknown keyword: ", token)
		}
	case IDENTIFIER:
		token0 := c.tokenizer.Identifier()
		c.NextToken()
		// assert c.tokenizer.HasMoreToken() == true, there is at least a ';' to follow
		token1Type := c.tokenizer.TokenType()
		if token1Type == SYMBOL && c.tokenizer.Symbol() == "[" { //   varName[ expression ]
			c.NextToken() //  [
			c.CompileExpression()
			kindName, index := c.ST.KindOf(token0), c.ST.IndexOf(token0)
			c.writer.writePush(kind2Seg[kindName], index)
			c.writer.WriteArithmetic("add")
			c.writer.writePop(POINTER, 1)
			c.writer.writePush(THAT, 0)
			c.NextToken() //  ]
		} else if token1Type == SYMBOL && c.tokenizer.Symbol() == "(" { //  subroutineCall case 1: call a method in current method
			// ( expressionList )
			c.NextToken() //  (
			nArgs := 1
			c.writer.writePush(POINTER, 0) // first argument: the address of THIS segement
			nArgs += c.CompileExpressionList()
			c.writer.writeCall(c.ST.className+"."+token0, nArgs)
			c.NextToken() //  )
		} else if token1Type == SYMBOL && c.tokenizer.Symbol() == "." { //  subroutineCall case 2
			// .subroutineName
			c.NextToken() //  .
			subroutineName := c.tokenizer.Identifier()
			kindName, typeName, index := c.ST.KindOf(token0), c.ST.TypeOf(token0), c.ST.IndexOf(token0)
			nArgs := 0
			if c.ST.IndexOf(token0) >= 0 {
				// if token0 is a variable(static,field,var,arg), then subroutine is a method
				c.writer.writePush(kind2Seg[kindName], index)
				nArgs++
			}
			c.NextToken()
			// ( expressionList )
			c.NextToken() //  (
			nArgs += c.CompileExpressionList()
			c.writer.writeCall(typeName+"."+subroutineName, nArgs)
			c.NextToken() //  )
		} else if c.ST.IndexOf(token0) >= 0 {
			// variable(static,field,var,arg)
			kindName, index := c.ST.KindOf(token0), c.ST.IndexOf(token0)
			c.writer.writePush(kind2Seg[kindName], index)
		} else {
			fmt.Println("unknown identifier: ", token0)
		}
	}
}

func (c *CompilationEngine) CompileExpression() {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileExpression fail")
	}
	c.CompileTerm()
	for c.tokenizer.HasMoreTokens() && c.tokenizer.TokenType() == SYMBOL {
		token := c.tokenizer.Symbol()
		if !isOp[token] {
			break
		}
		c.NextToken() //  op
		c.CompileTerm()
		c.writer.WriteArithmetic(token)
	}

}

// return the number of expression
func (c *CompilationEngine) CompileExpressionList() int {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("CompileExpressionList fail")
	}
	cnt := 0
	if c.tokenizer.TokenType() != SYMBOL || c.tokenizer.Symbol() != ")" {
		c.CompileExpression()
		cnt = 1
		for c.tokenizer.HasMoreTokens() && c.tokenizer.TokenType() == SYMBOL && c.tokenizer.Symbol() == "," {
			c.NextToken() //  ,
			c.CompileExpression()
			cnt++
		}
	}
	return cnt
}

func (c *CompilationEngine) NextToken() {
	if c.tokenizer.HasMoreTokens() {
		c.tokenizer.Advance()
	} else {
		log.Fatal("NextToken fail")
	}
}

func (c *CompilationEngine) writeIdentifier() string {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("writeIdentifier fail")
	}
	token := c.tokenizer.Identifier()
	curKind := c.ST.curKind // record current symbol's kind
	c.ST.curKind = ""       // turn off
	curType := c.ST.curType // record current symbol's type
	c.ST.curType = ""       // turn off
	if curKind != "" && curType != "" {
		// add new symbol into symbol table
		// when ethier kind or type are not empty.
		c.ST.Define(token, curKind, curType)
	}
	c.NextToken()
	return token
}

func (c *CompilationEngine) WriteKeyword() string {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("WriteKeyword fail")
	}
	token := c.tokenizer.KeyWord()
	c.NextToken()
	if token == "static" || token == "field" || token == "var" {
		c.ST.curKind = token // record current symbol's kind
	}
	return token
}

func (c *CompilationEngine) WriteType() string {
	if !c.tokenizer.HasMoreTokens() {
		log.Fatal("WriteType fail")
	}
	typeName := ""
	if c.tokenizer.TokenType() == KEYWORD {
		typeName = c.WriteKeyword()
	} else if c.tokenizer.TokenType() == IDENTIFIER {
		typeName = c.writeIdentifier()
	} else {
		log.Fatal("WriteType encounter invalid token", c.tokenizer.TokenType())
	}
	c.ST.curType = typeName // // record current symbol's types
	return typeName
}

//*******************************************************************************************************************//
//*************************************************** SymbolTable ***************************************************//

type SymbolTable struct {
	ifCnt     int
	whileCnt  int
	className string // current class's name
	funcName  string // current subroutine's name
	curKind   string
	curType   string
	staticST  map[string]*STEntry
	fieldST   map[string]*STEntry
	argST     map[string]*STEntry
	varST     map[string]*STEntry
}

type STEntry struct {
	name     string
	kindName string // static field arg var
	typeName string // string int char
	index    int
}

func SymbolTableConstructor() *SymbolTable {
	return &SymbolTable{}
}

// reset the class level symbol table
func (st *SymbolTable) startClass(className string) {
	st.ifCnt = 0
	st.whileCnt = 0
	st.className = className
	st.staticST = make(map[string]*STEntry)
	st.fieldST = make(map[string]*STEntry)
	st.argST = make(map[string]*STEntry)
	st.varST = make(map[string]*STEntry)
}

// reset the subroutine level symbol table
func (st *SymbolTable) startSubroutine(funcName string) {
	st.ifCnt = 0
	st.whileCnt = 0
	st.funcName = funcName
	st.argST = make(map[string]*STEntry)
	st.varST = make(map[string]*STEntry)
}

// add a new STEntry into corresponding Symbol table
func (st *SymbolTable) Define(name, kindName, typeName string) {
	entry := &STEntry{
		name:     name,
		kindName: kindName,
		typeName: typeName,
	}
	switch kindName {
	case STATIC:
		entry.index = len(st.staticST)
		st.staticST[entry.name] = entry
	case FIELD:
		entry.index = len(st.fieldST)
		st.fieldST[entry.name] = entry
	case ARG:
		entry.index = len(st.argST)
		st.argST[entry.name] = entry
	case VAR:
		entry.index = len(st.varST)
		st.varST[entry.name] = entry
	}
}

func (st *SymbolTable) VarCount(kindName string) int {
	switch kindName {
	case STATIC:
		return len(st.staticST)
	case FIELD:
		return len(st.fieldST)
	case ARG:
		return len(st.argST)
	case VAR:
		return len(st.varST)
	default:
		return -1
	}
}

// return the kind of the identifier specified by name
func (st *SymbolTable) KindOf(name string) string {
	if st.varST[name] != nil {
		return st.varST[name].kindName
	}
	if st.argST[name] != nil {
		return st.argST[name].kindName
	}
	if st.fieldST[name] != nil {
		return st.fieldST[name].kindName
	}
	if st.staticST[name] != nil {
		return st.staticST[name].kindName
	}
	return NONE
}

func (st *SymbolTable) TypeOf(name string) string {
	if st.varST[name] != nil {
		return st.varST[name].typeName
	}
	if st.argST[name] != nil {
		return st.argST[name].typeName
	}
	if st.fieldST[name] != nil {
		return st.fieldST[name].typeName
	}
	if st.staticST[name] != nil {
		return st.staticST[name].typeName
	}
	if isBuiltInClass[name] {
		return name
	}
	return name
}

func (st *SymbolTable) IndexOf(name string) int {
	if st.varST[name] != nil {
		return st.varST[name].index
	}
	if st.argST[name] != nil {
		return st.argST[name].index
	}
	if st.fieldST[name] != nil {
		return st.fieldST[name].index
	}
	if st.staticST[name] != nil {
		return st.staticST[name].index
	}
	return -1
}

//*******************************************************************************************************************//
//**************************************************** VMWriter *****************************************************//

type VMWriter struct {
	output *bufio.Writer
}

func VMWriterConstructor(output *bufio.Writer) VMWriter {
	return VMWriter{
		output: output,
	}
}

func (w *VMWriter) writePush(segment string, index int) {
	w.output.WriteString(fmt.Sprintf("push %s %v\n", segment, index))
}

func (w *VMWriter) writePop(segment string, index int) {
	w.output.WriteString(fmt.Sprintf("pop %s %v\n", segment, index))
}

func (w *VMWriter) WriteArithmetic(command string) {
	if !isOp[command] {
		w.output.WriteString(command + "\n")
		return
	}
	switch command {
	case "+":
		w.output.WriteString("add\n")
	case "-":
		w.output.WriteString("sub\n")
	case "*":
		w.output.WriteString("call Math.multiply 2\n")
	case "/":
		w.output.WriteString("call Math.divide 2\n")
	case "&amp;", "&":
		w.output.WriteString("and\n")
	case "|":
		w.output.WriteString("or\n")
	case "&lt;", "<":
		w.output.WriteString("lt\n")
	case "&gt;", ">":
		w.output.WriteString("gt\n")
	case "=":
		w.output.WriteString("eq\n")
	default:
		log.Fatal("unknown arithmetic command: ", command)
	}

}

func (w *VMWriter) writeLabel(label string) {
	w.output.WriteString(fmt.Sprintf("label %s\n", label))
}

func (w *VMWriter) writeGoto(label string) {
	w.output.WriteString(fmt.Sprintf("goto %s\n", label))
}

func (w *VMWriter) writeIfGoto(label string) {
	w.output.WriteString(fmt.Sprintf("if-goto %s\n", label))
}

func (w *VMWriter) writeCall(funcName string, nArgs int) {
	w.output.WriteString(fmt.Sprintf("call %s %v\n", funcName, nArgs))
}

func (w *VMWriter) writeFunction(funcName string, nLocals int) {
	w.output.WriteString(fmt.Sprintf("function %s %v\n", funcName, nLocals))
}

func (w *VMWriter) writeReturn() {
	w.output.WriteString("return\n")
}

func (w *VMWriter) close() {
	if w.output == nil {
		return
	}
	w.output.Flush()
}
