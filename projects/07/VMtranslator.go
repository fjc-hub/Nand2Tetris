package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	C_ARITHMETIC = "c_arithmetic"
	C_PUSH       = "c_push"
	C_POP        = "c_pop"
	C_LABEL      = "c_label"
	C_GOTO       = "c_goto"
	C_IF         = "c_if"
	C_FUNCTION   = "c_function"
	C_RETURN     = "c_return"
	C_CALL       = "c_call"
)

func main() {
	path := os.Args[1]
	fileDir, pureName, fileExt, isFile := parsePath(path)
	if !isFile || fileExt != ".vm" {
		log.Fatal("invalid file")
	}
	outFile := fileDir + "/" + pureName + ".asm"

	input, err := os.OpenFile(path, os.O_RDONLY, 777)
	defer input.Close()
	if err != nil {
		log.Fatalf("open input file error: %v", err)
	}
	output, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 777)
	defer output.Close()
	if err != nil {
		log.Fatalf("open output file error: %v", err)
	}
	parser := &Parser{
		input: bufio.NewScanner(input),
	}
	codeWriter := &CodeWriter{
		output:   bufio.NewWriter(output),
		fileName: pureName,
	}
	for parser.hasMoreCommands() {
		parser.advance()
		cmdType := parser.commandType()
		switch cmdType {
		case C_ARITHMETIC:
			codeWriter.writeArithmetic(parser.arg1())
		case C_PUSH:
			codeWriter.writePush(parser.arg1(), parser.arg2())
		case C_POP:
			codeWriter.writePop(parser.arg1(), parser.arg2())
		case C_LABEL:
		case C_GOTO:
		case C_IF:
		case C_FUNCTION:
		case C_RETURN:
		case C_CALL:
		default:
			log.Fatalf("error command type in switch: %v", cmdType)
		}
	}
	codeWriter.close()
	fmt.Println("translator exit successfully")
}

//**********************************************utils****************************************************************//

func parsePath(path string) (string, string, string, bool) {
	fileName := filepath.Base(path)
	fileExt := filepath.Ext(path)
	idx := strings.LastIndexByte(fileName, '.')
	if idx == -1 {
		log.Fatal("invalid file")
	}
	pureName := fileName[:idx]
	fileDir := filepath.Dir(path)
	isFile := fileName[0] != '.' && fileName[0] != '/'
	return fileDir, pureName, fileExt, isFile
}

//*******************************************************************************************************************//

//Parser: parses each VM command into its lexical elements
type Parser struct {
	input           *bufio.Scanner
	currentCmd      string
	var_commandType string
	var_arg1        string
	var_arg2        int
}

func (p *Parser) hasMoreCommands() bool {
	// add currentCmd as condition to realize idempotent
	for p.currentCmd == "" && p.input.Scan() {
		str := strings.Trim(p.input.Text(), " ")
		if str == "" || str[0] == '/' {
			continue
		} else {
			comments := strings.Index(str, "//")
			if comments != -1 {
				str = str[:comments]
			}
			p.currentCmd = str
			return true
		}
	}
	return false
}

// read next command from input and treat it as current command
func (p *Parser) advance() {
	if p.currentCmd == "" {
		log.Fatalln("inoke incorrectly advance()")
	}
	var err error
	cmd := p.currentCmd
	p.currentCmd = ""
	// splits the string s around each instance of one or more consecutive white space
	compnents := strings.Fields(cmd)
	// identify var_commandType
	switch compnents[0] {
	case "push", "pop":
		if compnents[0] == "push" {
			p.var_commandType = C_PUSH
		} else {
			p.var_commandType = C_POP
		}
		p.var_arg1 = compnents[1]
		p.var_arg2, err = strconv.Atoi(compnents[2])
		if err != nil {
			log.Fatalf("switch case string to integer error: %v\n", err)
		}
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		p.var_commandType = C_ARITHMETIC
		p.var_arg1 = compnents[1]
	case "return":
		p.var_commandType = C_RETURN
	default:
		log.Fatalf("invalid key word %v\n", compnents[0])
		return
	}
}

// return the type of current command (C_ARITHMETIC represents all the arithmetic/logic commands)
func (p *Parser) commandType() string {
	if p.var_commandType == "" {
		log.Fatalln("inoke incorrectly commandType()")
	}
	return p.var_commandType
}

// return the first argument of current command.
// if command type is C_ARITHMETIC(add, sub, etc.) return itself (because it don't has arguments)
// Should not be called if current command is C_RETURN.
func (p *Parser) arg1() string {
	if p.var_commandType == C_RETURN {
		log.Fatalln("inoke incorrectly arg1()")
	}
	return p.var_arg1
}

// return the second argument of current command. Should be called
// only if current command is C_POP, C_PUSH, C_CALL or C_FUNCTION
func (p *Parser) arg2() int {
	if p.var_commandType != C_POP && p.var_commandType != C_PUSH &&
		p.var_commandType != C_CALL && p.var_commandType != C_FUNCTION {
		log.Fatalln("inoke incorrectly arg2()")
	}
	return p.var_arg2
}

//CodeWriter: writes the assembly code that implements the parsed command
type CodeWriter struct {
	output     *bufio.Writer
	count      int // record the number of assembly command
	fileName   string
	static_cnt int // the number of static variable of this VM file (named fileName)
}

// write to the output file the assembly code that implements the given arithmetic command
func (c CodeWriter) writeArithmetic(command string) {
	switch command {
	case "add":
		c.output.WriteString(`
			@SP 	// A=sp
			AM=M-1 	// A=(*sp)-1; (*sp)--
			D=M 	// D=*(*sp)
			A=A-1	// A=(*sp)-1
			M=M+D
		`)
	case "sub":
		c.output.WriteString(`
			@SP		// A=sp
			AM=M-1	// A=(*sp)-1; (*sp)--
			D=M		// D=*(*sp)
			A=A-1	// A=(*sp)-1
			M=M-D
		`)
	case "neg":
		c.output.WriteString(`
			@SP		// A=sp
			A=M-1	// A=(*sp)-1
			M=-M
		`)
	case "eq":
		cnt := c.count
		c.count++
		c.output.WriteString(fmt.Sprintf(`
			@SP
			AM=M-1
			D=M
			A=A-1
			@eq.true.%v
			M-D;JEQ
			M=0
			@eq.skip.%v
			0; JMP
			(eq.true.%v)
			M=-1
			(eq.skip.%v)
		`, cnt, cnt, cnt, cnt),
		)
	case "gt":
		cnt := c.count
		c.count++
		c.output.WriteString(fmt.Sprintf(`
			@SP
			AM=M-1
			D=M
			A=A-1
			@gt.true.%v
			M-D;JGT
			M=0
			@gt.skip.%v
			0; JMP
			(gt.true.%v)
			M=-1
			(gt.skip.%v)
		`, cnt, cnt, cnt, cnt),
		)
	case "lt":
		cnt := c.count
		c.count++
		c.output.WriteString(fmt.Sprintf(`
			@SP
			AM=M-1
			D=M
			A=A-1
			@lt.true.%v
			M-D;JLT
			M=0
			@lt.skip.%v
			0; JMP
			(lt.true.%v)
			M=-1
			(lt.skip.%v)
		`, cnt, cnt, cnt, cnt),
		)
	case "and":
		c.output.WriteString(`
			@SP		// A=sp
			AM=M-1	// A=(*sp)-1; (*sp)--
			D=M		// D=*(*sp)
			A=A-1	// A=(*sp)-1
			M=M&D
		`)
	case "or":
		c.output.WriteString(`
			@SP		// A=sp
			AM=M-1	// A=(*sp)-1; (*sp)--
			D=M		// D=*(*sp)
			A=A-1	// A=(*sp)-1
			M=M|D
		`)
	case "not":
		c.output.WriteString(`
			@SP
			A=M-1
			M=!M
		`)
	}
	c.output.WriteString(`
		@SP
		M=M+1
	`)
}

// write to the output file the assembly code that implements the given arithmetic command,
// where command is C_PUSH
func (c CodeWriter) writePush(segment string, index int) {
	switch segment {
	case "argument":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@ARG
			A=M+D
			D=M
		`, index))
	case "local":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@LCL
			A=M+D
			D=M
		`, index))
	case "static":
		cnt := strconv.Itoa(c.static_cnt)
		c.static_cnt++
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@%v
			A=M+D
			D=M
		`, index, c.fileName+cnt))
	case "constant":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
		`, index))
	case "this":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@THIS
			A=M+D
			D=M
		`, index))
	case "that":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@THAT
			A=M+D
			D=M
		`, index))
	case "pointer":
		if index == 0 {
			c.output.WriteString(`
				@THIS
				D=M
			`)
		} else {
			c.output.WriteString(`
				@THAT
				D=M
			`)
		}
	case "temp":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@R5
			A=M+D
			D=M
		`, index))
	default:
		log.Fatalf("invalid segement %v\n", segment)
		return
	}
	c.output.WriteString(`
		@SP		// (*sp) = D
		A=M
		M=D
		@SP		// sp++
		M=M+1 
	`)
}

// write to the output file the assembly code that implements the given arithmetic command,
// where command is C_POP
func (c CodeWriter) writePop(segment string, index int) {
	switch segment {
	case "argument":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@ARG
			D=M+D
		`, index))
	case "local":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@LCL
			D=M+D
		`, index))
	case "static":
		cnt := strconv.Itoa(c.static_cnt)
		c.static_cnt++
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@%v
			D=M+D
		`, index, c.fileName+cnt))
	case "this":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@THIS
			D=M+D
		`, index))
	case "that":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@THAT
			D=M+D
		`, index))
	case "pointer":
		if index == 0 {
			c.output.WriteString(`
				@THIS
				D=M
			`)
		} else {
			c.output.WriteString(`
				@THAT
				D=M
			`)
		}
	case "temp":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@R5
			D=M+D
		`, index))
	default:
		log.Fatalf("invalid segement %v\n", segment)
		return
	}
	c.output.WriteString(`
		@R13
		M=D
		@SP
		AM=M-1
		D=M
		@R13
		M=D
	`)
}

// flush and close
func (c CodeWriter) close() {
	c.output.Flush()
}
