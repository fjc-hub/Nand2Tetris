package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	fileName := flag.String("filename", "", "")
	flag.Parse()
	if *fileName == "" {
		fmt.Println("what's wrong with you, bro. ")
	}
	// Input: you can omit 777 file permission bit, because no os.O_CREATE flag
	file, err := os.OpenFile(*fileName, os.O_RDONLY, 777)
	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	scan := bufio.NewScanner(file)
	//一、First scan through all instruction to number it and construct symbal table
	// 1.1 initialize symbal table: insert pre-defined symbal
	ST := initSymbalTable()
	// 1.2 remove empty line, remove comments and extend ST
	instructions := make([][]byte, 0)
	for scan.Scan() {
		// The underlying array may point to data that will be overwritten
		// by a subsequent call to Scan. It does no allocation.
		bts := scan.Bytes() // better manipulate the copy of its returned byte slice
		if len(bts) == 0 || bts[0] == '/' {
			continue
		}
		if bts[0] == '(' {
			// label symbal
			label := string(bts[1 : len(bts)-1])
			ST[label] = len(instructions)
		} else {
			instructions = append(instructions, extractInstraction(bts))
		}
	}
	if err := scan.Err(); err != nil {
		log.Fatalf("scanner loop error: %v", err)
	}
	//二、Second scan all instructions: translation
	binaryCodes := make([]string, 0)
	varIndex := 16 // variable symbal address start from 16
	for _, r := range instructions {
		if r[0] == '@' {
			// A-instruction
			num := -1
			if unicode.IsDigit(rune(r[1])) {
				// number
				num, err = strconv.Atoi(string(r[1:]))
				if err != nil {
					log.Fatal(err)
				}
			} else {
				// symbal
				symbal := string(r[1:])
				ok := false // fucking rubbish golang
				num, ok = ST[symbal]
				if !ok {
					// variable symbal
					num, ST[symbal] = varIndex, varIndex
					varIndex++
				}
			}
			binaryCodes = append(binaryCodes, "0"+numToBinaryString(num)+"\n")
		} else {
			// C-instruction
			comp, dest, jump := c_split(r)
			b_comp := c_parse_comp(comp)
			b_dest := c_parse_dest(dest)
			b_jump := c_parse_jump(jump)
			binaryCodes = append(binaryCodes, "111"+b_comp+b_dest+b_jump+"\n")
		}
	}

	// Output
	name := strings.Split(*fileName, ".")[0]
	outFile, err := os.OpenFile(name+".hack", os.O_CREATE|os.O_WRONLY, 0644)
	defer outFile.Close()
	if err != nil {
		log.Fatalf("open output file fail %v", err)
	}
	for i := range binaryCodes {
		outFile.WriteString(binaryCodes[i])
	}
	fmt.Println("assembler exit successfully!!!")
}

func initSymbalTable() map[string]int {
	st := map[string]int{
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
		"SCREEN": 16384,
		"KBD":    24576,
	}
	for i := 0; i <= 15; i++ {
		st["R"+strconv.Itoa(i)] = i
	}
	return st
}

// remove all comments and space character in a line
func extractInstraction(bts []byte) []byte {
	comments := bytes.Index(bts, []byte("//"))
	if comments != -1 {
		bts = bts[:comments]
	}
	bts = bytes.Trim(bts, " ")
	ans := make([]byte, len(bts))
	copy(ans, bts)
	return ans
}

// returned string' length always equals 15
func numToBinaryString(num int) string {
	ans := make([]byte, 0)
	for num > 1 {
		t := byte(num%2 + '0')
		ans = append(ans, t)
		num /= 2
	}
	if num == 1 {
		ans = append(ans, '1')
	}
	for len(ans) < 15 {
		ans = append(ans, '0')
	}
	// manul reversion, ): fucking golang
	for i, j := 0, len(ans)-1; i < j; i, j = i+1, j-1 {
		ans[i], ans[j] = ans[j], ans[i]
	}
	return string(string(ans))
}

func c_split(arr []byte) (string, string, string) {
	var comp, dest, jump []byte
	delim0 := bytes.IndexByte(arr, '=')
	if delim0 != -1 {
		dest = arr[:delim0]
	}
	delim1 := bytes.IndexByte(arr[delim0+1:], ';')
	if delim1 == -1 {
		comp = arr[delim0+1:]
	} else {
		comp = arr[delim0+1 : delim0+1+delim1]
		jump = arr[delim0+2+delim1:]
	}
	// trim empty character
	bytes.Trim(comp, " ")
	bytes.Trim(dest, " ")
	bytes.Trim(jump, " ")
	return string(comp), string(dest), string(jump)
}

// map the comp part of C-instruction to binary string
func c_parse_comp(str string) string {
	switch str {
	case "":
		return "1101010"
	case "0":
		return "0101010"
	case "1":
		return "0111111"
	case "-1":
		return "0111010"
	case "D":
		return "0001100"
	case "A":
		return "0110000"
	case "M":
		return "1110000"
	case "!D":
		return "0001101"
	case "!A":
		return "0110001"
	case "!M":
		return "1110001"
	case "-D":
		return "0001111"
	case "-A":
		return "0110011"
	case "-M":
		return "1110011"
	case "D+1":
		return "0011111"
	case "A+1":
		return "0110111"
	case "M+1":
		return "1110111"
	case "D-1":
		return "0001110"
	case "A-1":
		return "0110010"
	case "M-1":
		return "1110010"
	case "D+A":
		return "0000010"
	case "D+M":
		return "1000010"
	case "D-A":
		return "0010011"
	case "D-M":
		return "1010011"
	case "A-D":
		return "0000111"
	case "M-D":
		return "1000111"
	case "D&A":
		return "0000000"
	case "D&M":
		return "1000000"
	case "D|A":
		return "0010101"
	case "D|M":
		return "1010101"
	default:
		log.Fatalf("c_parse_comp error expression: %v", str)
		return ""
	}
}

// map the dest part of C-instruction to binary string
func c_parse_dest(str string) string {
	switch str {
	case "":
		return "000"
	case "M":
		return "001"
	case "D":
		return "010"
	case "MD":
		return "011"
	case "A":
		return "100"
	case "AM":
		return "101"
	case "AD":
		return "110"
	case "AMD":
		return "111"
	default:
		log.Fatalf("c_parse_dest error expression: %v", str)
		return ""
	}
}

// map the jump part of C-instruction to binary string
func c_parse_jump(str string) string {
	switch str {
	case "":
		return "000"
	case "JGT":
		return "001"
	case "JEQ":
		return "010"
	case "JGE":
		return "011"
	case "JLT":
		return "100"
	case "JNE":
		return "101"
	case "JLE":
		return "110"
	case "JMP":
		return "111"
	default:
		log.Fatalf("c_parse_jump error expression: %v", str)
		return ""
	}
}
