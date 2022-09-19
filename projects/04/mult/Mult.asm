// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
//
// This program only needs to handle arguments that satisfy
// R0 >= 0, R1 >= 0, and R0*R1 < 32768.

// Put your code here.

//pseudo code
// M[@R2] = 0       // 初始化, avoid M[@R0] == 0 then cause error result
// for M[@R0] > 0 {
//     M[@R0]--
//     M[@R2] = M[@R2] + M[@R1]
// }

@R2
M=0             // reset M[@R2]
(loop)
    @R0
    D=M
    @end
    D; JLE      // M[@R0] <= 0
    @R0
    M = D - 1
    @R1
    D=M
    @R2
    M=M+D
    @loop
    0; JMP

(end)
    @end
    0; JMP