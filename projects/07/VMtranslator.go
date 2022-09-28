package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
	inFile := os.Args[1]
	outFile := os.Args[2]
	input, err := os.OpenFile(inFile, os.O_RDONLY, 777)
	defer input.Close()
	if err != nil {
		log.Fatalf("open input file error: %v", err)
	}
	output, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 777)
	defer output.Close()
	if err != nil {
		log.Fatalf("open output file error: %v", err)
	}
	parser := Parser{
		input: bufio.NewScanner(input),
	}
	codeWriter := CodeWriter{
		output: bufio.NewWriter(output),
	}
	for ; parser.hasMoreCommands(); parser.advance() {
		cmdType := parser.commandType()
		switch cmdType {
		case C_ARITHMETIC:
			codeWriter.writeArithmetic(parser.arg1())
		case C_PUSH:
			codeWriter.writePushPop(C_PUSH, parser.arg1(), parser.arg2())
		case C_POP:
			codeWriter.writePushPop(C_POP, parser.arg1(), parser.arg2())
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

//Parser: parses each VM command into its lexical elements
type Parser struct {
	input           *bufio.Scanner
	currentCmd      string
	var_commandType string
	var_arg1        string
	var_arg2        int
}

func (p Parser) hasMoreCommands() bool {
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
func (p Parser) advance() {
	if p.currentCmd == "" {
		log.Fatalln("inoke incorrectly advance()")
	}
	var err error
	cmd := p.var_commandType
	p.var_commandType = ""
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
func (p Parser) commandType() string {
	if p.var_commandType == "" {
		log.Fatalln("inoke incorrectly commandType()")
	}
	return p.var_commandType
}

// return the first argument of current command.
// if command type is C_ARITHMETIC(add, sub, etc.) return itself (because it don't has arguments)
// Should not be called if current command is C_RETURN.
func (p Parser) arg1() string {
	if p.var_commandType == C_RETURN {
		log.Fatalln("inoke incorrectly arg1()")
	}
	return p.var_arg1
}

// return the second argument of current command. Should be called
// only if current command is C_POP, C_PUSH, C_CALL or C_FUNCTION
func (p Parser) arg2() int {
	if p.var_commandType == C_POP || p.var_commandType == C_PUSH ||
		p.var_commandType == C_CALL || p.var_commandType == C_FUNCTION {
		log.Fatalln("inoke incorrectly arg2()")
	}
	return p.var_arg2
}

//CodeWriter: writes the assembly code that implements the parsed command
type CodeWriter struct {
	output *bufio.Writer
}

// write to the output file the assembly code that implements the given arithmetic command
func (c CodeWriter) writeArithmetic(command string) {
	// c.output.WriteString("sasd")
	switch command {
	case "add":
		c.output.WriteString(`
			@SP // store M[sp-1] into D
			A=A-1
			D=M
			@SP // M[sp]=M[sp]-2
			M=M-2
			@SP // M[sp]=M[sp]+D
			M=M+D
		`)
	case "sub":
		c.output.WriteString(`
			@SP // store M[sp-1] into D
			A=A-1
			D=M
			@SP // M[sp]=M[sp]-2
			M=M-2
			@SP // M[sp]=M[sp]-D
			M=M-D
		`)
	case "neg":
		c.output.WriteString(`
			@SP // store M[sp-1] into A
			A=A-1
			M=-M
		`)
	case "eq":
		c.output.WriteString(`
			
		`)
	case "gt":
	case "lt":
	case "and":
	case "or":
	case "not":
	}

}

// write to the output file the assembly code that implements the given arithmetic command,
// where command is either C_PUSH or C_POP
func (c CodeWriter) writePushPop(cmd, segment string, index int) {

}

// flush and close
func (c CodeWriter) close() {
	c.output.Flush()
}
