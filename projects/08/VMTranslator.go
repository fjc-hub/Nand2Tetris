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
	C_IF_GOTO    = "c_if_goto"
	C_FUNCTION   = "c_function"
	C_RETURN     = "c_return"
	C_CALL       = "c_call"
)

const BOOT_FILE = "Sys"

func main() {
	path := os.Args[1]
	fileDir, pureName, fileExt, isFile := parsePath(path)
	if isFile && fileExt != ".vm" {
		log.Fatal("invalid file name")
	}
	outFile := fileDir + "/" + pureName + ".asm"
	if !isFile {
		outFile = path + "/" + pureName + ".asm"
	}

	parsers := make([]*Parser, 0)
	if isFile {
		input, err := os.OpenFile(path, os.O_RDONLY, 777)
		defer input.Close()
		if err != nil {
			log.Fatalf("open input file error: %v", err)
		}
		parsers = append(parsers, &Parser{
			input:    bufio.NewScanner(input),
			fileName: pureName,
		})
	} else {
		parsers = []*Parser{&Parser{
			fileName: BOOT_FILE,
		}}
		// isRoot := true
		walkDirNoRecur(path, func(pt string, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}
			// if isRoot {
			// 	isRoot = false
			// 	return nil
			// }
			_, pureName, ext, isFile := parsePath(pt)
			if !isFile || ext != ".vm" {
				return nil
			}
			fmt.Printf("read file: %s\n", pt)
			input, err := os.OpenFile(pt, os.O_RDONLY, 777)
			// defer input.Close()
			if err != nil {
				log.Fatalf("open input file error: %v", err)
			}
			if pureName == BOOT_FILE {
				parsers[0].input = bufio.NewScanner(input)
			} else {
				parsers = append(parsers, &Parser{
					input:    bufio.NewScanner(input),
					fileName: pureName,
				})
			}
			return nil
		})
	}

	output, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 777)
	defer output.Close()
	if err != nil {
		log.Fatalf("open output file error: %v", err)
	}
	codeWriter := &CodeWriter{
		output:   bufio.NewWriter(output),
		fileName: pureName,
	}
	if !isFile {
		codeWriter.WriteInit()
	}
	for _, parser := range parsers {
		codeWriter.fileName = parser.fileName
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
				codeWriter.writeLabel(parser.arg1())
			case C_GOTO:
				codeWriter.writeGoto(parser.arg1())
			case C_IF_GOTO:
				codeWriter.writeIfGotoGoto(parser.arg1())
			case C_FUNCTION:
				codeWriter.writeFunction(parser.arg1(), parser.arg2())
			case C_CALL:
				codeWriter.writeCall(parser.arg1(), parser.arg2())
			case C_RETURN:
				codeWriter.writeReturn()
			default:
				log.Fatalf("error command type in switch: %v", cmdType)
			}
		}
	}
	codeWriter.close()
	fmt.Println("translator exit successfully")
}

//***************************************************utils***********************************************************//

func parsePath(path string) (string, string, string, bool) {
	fileName := filepath.Base(path)
	fileExt := filepath.Ext(path)
	pureName := fileName
	if idx := strings.LastIndexByte(fileName, '.'); idx != -1 {
		pureName = fileName[:idx]
	}
	fileDir := filepath.Dir(path)
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal("os.Stat error: ", err, path)
	}
	isFile := !fileInfo.IsDir()
	return fileDir, pureName, fileExt, isFile
}

func walkDirNoRecur(path string, todo func(string, error) error) error {
	filenames, err := filepath.Glob(path + "/*")
	if err != nil {
		log.Fatal(err)
	}
	for _, name := range filenames {
		err = todo(name, err)
	}
	return nil
}

//*******************************************************************************************************************//

//******************************************************Parser*******************************************************//
//Parser: parses each VM command into its lexical elements
type Parser struct {
	input           *bufio.Scanner
	fileName        string
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
		p.var_arg1 = compnents[0]
	case "label":
		p.var_commandType = C_LABEL
		p.var_arg1 = compnents[1]
	case "goto":
		p.var_commandType = C_GOTO
		p.var_arg1 = compnents[1]
	case "if-goto":
		p.var_commandType = C_IF_GOTO
		p.var_arg1 = compnents[1]
	case "function":
		p.var_commandType = C_FUNCTION
		p.var_arg1 = compnents[1]
		p.var_arg2, err = strconv.Atoi(compnents[2])
		if err != nil {
			log.Fatalf("switch case string to integer error: %v\n", err)
		}
	case "call":
		p.var_commandType = C_CALL
		p.var_arg1 = compnents[1]
		p.var_arg2, err = strconv.Atoi(compnents[2])
		if err != nil {
			log.Fatalf("switch case string to integer error: %v\n", err)
		}
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

//*******************************************************************************************************************//

//****************************************************CodeWriter*****************************************************//

//CodeWriter: writes the assembly code that implements the parsed command
type CodeWriter struct {
	fileName         string
	output           *bufio.Writer
	truthCnt         int // record the number of labels of assembly command gt,eq,lt
	static_cnt       int // the number of static variable of this VM file (named fileName)
	func_retAddr_cnt int
	func_name        string
}

// write to the output file the assembly code that implements the given arithmetic command
func (c *CodeWriter) writeArithmetic(command string) {
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
		cnt := c.truthCnt
		c.truthCnt++
		c.output.WriteString(fmt.Sprintf(`
			@SP
			AM=M-1
			D=M
			A=A-1
			D=M-D
			@eq.true.%v
			D;JEQ
			@SP
			A=M-1
			M=0
			@eq.skip.%v
			0; JMP
			(eq.true.%v)
			@SP
			A=M-1
			M=-1
			(eq.skip.%v)
		`, cnt, cnt, cnt, cnt),
		)
	case "gt":
		cnt := c.truthCnt
		c.truthCnt++
		c.output.WriteString(fmt.Sprintf(`
			@SP
			AM=M-1
			D=M
			A=A-1
			D=M-D
			@gt.true.%v
			D;JGT
			@SP
			A=M-1
			M=0
			@gt.skip.%v
			0; JMP
			(gt.true.%v)
			@SP
			A=M-1
			M=-1
			(gt.skip.%v)
		`, cnt, cnt, cnt, cnt),
		)
	case "lt":
		cnt := c.truthCnt
		c.truthCnt++
		c.output.WriteString(fmt.Sprintf(`
			@SP
			AM=M-1
			D=M
			A=A-1
			D=M-D
			@lt.true.%v
			D;JLT
			@SP
			A=M-1
			M=0
			@lt.skip.%v
			0; JMP
			(lt.true.%v)
			@SP
			A=M-1
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
	default:
		log.Fatalf("writeArithmetic: invalid command %v", command)
	}
}

// write to the output file the assembly code that implements the given arithmetic command,
// where command is C_PUSH
func (c *CodeWriter) writePush(segment string, index int) {
	switch segment {
	case "argument":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@ARG
			A=M+D
			D=M`, index),
		)
	case "local":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@LCL
			A=M+D
			D=M`, index),
		)
	case "static":
		idx := strconv.Itoa(index)
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=M`, c.fileName+"."+idx),
		)
	case "constant":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A`, index),
		)
	case "this":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@THIS
			A=M+D
			D=M`, index),
		)
	case "that":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@THAT
			A=M+D
			D=M`, index),
		)
	case "pointer":
		str := ""
		if index == 0 {
			str = "THIS"
		} else {
			str = "THAT"
		}
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=M`, str),
		)
	case "temp":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@R5
			A=A+D
			D=M`, index),
		)
	default:
		log.Fatalf("invalid segement %v\n", segment)
		return
	}
	// store the value of D-register into *SP and increase SP word by 1
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
func (c *CodeWriter) writePop(segment string, index int) {
	switch segment {
	case "argument":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@ARG
			D=M+D
		`, index),
		)
	case "local":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@LCL
			D=M+D
		`, index),
		)
	case "static":
		idx := strconv.Itoa(index)
		c.output.WriteString(fmt.Sprintf(`
			@SP
			AM=M-1
			D=M
			@%v
			M=D
		`, c.fileName+"."+idx),
		)
		return
	case "this":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@THIS
			D=M+D
		`, index),
		)
	case "that":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@THAT
			D=M+D
		`, index),
		)
	case "pointer":
		str := ""
		if index == 0 {
			str = "THIS"
		} else {
			str = "THAT"
		}
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
		`, str),
		)
	case "temp":
		c.output.WriteString(fmt.Sprintf(`
			@%v
			D=A
			@R5
			D=A+D
		`, index),
		)
	default:
		log.Fatalf("invalid segement %v\n", segment)
		return
	}
	// address computer word by the value of D-register
	// 按D寄存器中的值寻址（即把D中的值当作指针处理），得到一个字word，然后将栈顶那个字的值复制到这个字word中
	c.output.WriteString(`
			@R13
			M=D
			@SP
			AM=M-1
			D=M
			@R13
			A=M
			M=D
	`)
}

func (c *CodeWriter) writeLabel(symbol string) {
	c.output.WriteString(fmt.Sprintf("(%s$%s)\n", c.func_name, symbol))
}

func (c *CodeWriter) writeGoto(symbol string) {
	c.output.WriteString(fmt.Sprintf(`
			@%s$%s
			0;JMP
	`, c.func_name, symbol))
}

// if D == 0, then jump
func (c *CodeWriter) writeIfGotoGoto(symbol string) {
	c.output.WriteString(fmt.Sprintf(`
			@SP
			AM=M-1
			D=M			// pop the top-most value of the stack onto D-register
			@%s$%s
			D;JNE
	`, c.func_name, symbol))

}

func (c *CodeWriter) writeFunction(funcName string, nLcl int) {
	c.func_name = funcName
	c.output.WriteString(fmt.Sprintf(`
		(%s)
	`, funcName))
	for i := 0; i < nLcl; i++ {
		c.output.WriteString(`
			@SP
			A=M
			M=0
			@SP
			M=M+1
		`)
	}
}

func (c *CodeWriter) writeCall(funcName string, nArg int) {
	cnt := c.func_retAddr_cnt
	c.func_retAddr_cnt++
	// write the value of D-register onto stack
	writeStackTop := `@SP 
			A=M
			M=D
			@SP
			M=M+1`
	c.output.WriteString(fmt.Sprintf(`
			@%v.return-address.%v
			D=A		// push return-address
			%s
			@LCL	// store LCL state of caller
			D=M
			%s
			@ARG 	// store ARG state of caller
			D=M
			%s
			@THIS	// store THIS state of caller
			D=M
			%s
			@THAT	// store THAT state of caller
			D=M
			%s
			@SP 	// update new ARG
			D=M
			@%v
			D=D-A
			@5
			D=D-A
			@ARG
			M=D 	
			@SP		// update LCL pointer
			D=M
			@LCL
			M=D
			@%s		// transfer control
			0; JMP
		(%v.return-address.%v)
	`, funcName, cnt, writeStackTop, writeStackTop, writeStackTop, writeStackTop, writeStackTop, nArg, funcName, funcName, cnt))
}

func (c *CodeWriter) writeReturn() {
	c.output.WriteString(`
			@LCL 	// store return address to avoid being overwitten
			D=M
			@5
			A=D-A
			D=M
			@R13
			M=D
			@SP 	// transfer top-most value of stack onto ARG[0]
			A=M-1
			D=M
			@ARG
			A=M
			M=D
			D=A		// reset SP
			@SP
			M=D+1
					// restore all preserved states of caller
			@LCL	// restore THAT
			D=M
			@1
			A=D-A
			D=M
			@THAT
			M=D
			@LCL	// resotre THIS
			D=M
			@2
			A=D-A
			D=M
			@THIS
			M=D
			@LCL 	// restore ARG
			D=M
			@3
			A=D-A
			D=M
			@ARG
			M=D
			@LCL	// restore LCL
			D=M
			@4
			A=D-A
			D=M
			@LCL
			M=D
			@R13 	// jmp to ret-addr
			A=M
			0; JMP
	`)
}

//*******************************************************************************************************************//

func (c *CodeWriter) WriteInit() {
	c.output.WriteString(`
		@256	// set *SP = 256
		D=A
		@SP
		M=D
	`)
	c.writeCall("Sys.init", 0)
}

// flush and close
func (c *CodeWriter) close() {
	c.output.Flush()
}
